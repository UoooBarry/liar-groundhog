package session

import (
	"errors"
	"fmt"

	"uooobarry/liar-groundhog/internal/liar"
	"uooobarry/liar-groundhog/internal/message"
	"uooobarry/liar-groundhog/utils"
)

const MAX_PLAYERS = 4
const MIN_PLAYERS_TO_START = 4
const BET_CARD = liar.Ace

type Room struct {
	RoomUUID           string
	Players            []*Player
	engine             liar.Engine
	OwnerUUID          string
	playerCards        map[*Player][]liar.Card
	CurrentPlayerIndex int
	BetCard            liar.Card
	LastPlaceCards     []liar.Card
}

func CreateRoom(ownerUUID string, gameEngine *liar.Engine) (*Room, error) {
	return roomManager.CreateRoom(ownerUUID, gameEngine)
}

func FindRoom(uuid string) (*Room, bool) {
	return roomManager.FindRoom(uuid)
}

func (room *Room) FindPlayerByUsername(username *string) (*Player, bool) {
	for _, player := range room.Players {
		if player.Username == *username {
			return player, true
		}
	}

	return nil, false
}

func (room *Room) FindPlayerByUUID(uuid *string) (*Player, bool) {
	for _, player := range room.Players {
		if player.SessionUUID == *uuid {
			return player, true
		}
	}

	return nil, false
}

func (room *Room) Broadcast(msg message.MessageInterface) {
	for _, player := range room.Players {
		if player.Conn != nil {
			message.SendResponse(player.Conn, msg)
		}

	}
}

func (room *Room) SendPrivateMessages(messageFn func(*Room, *Player) message.MessageInterface) {
	for _, player := range room.Players {
		if player.Conn != nil {
			message.SendResponse(player.Conn, messageFn(room, player))
		}
	}
}

func (room *Room) PublishRoomInfo() {
	room.Broadcast(GetInfoMessage(room))
}

func (room *Room) PublishPlayerHoldingCards() {
	room.SendPrivateMessages(GetUserCardMessages)
}

func GetUserCardMessages(room *Room, p *Player) message.MessageInterface {
	return message.PlayerHoldingCardsMessage{
		Type:         "player_holding_cards",
		HoldingCards: room.playerCards[p],
		SessionUUID:  p.SessionUUID,
		Username:     p.Username,
	}
}

func GetInfoMessage(room *Room) message.RoomInfoMessage {
	playerListInfo := utils.MapSlice(room.Players, func(p *Player) message.PublicPlayerMessage {
		return message.PublicPlayerMessage{
			Username: p.Username,
		}
	})
	return message.RoomInfoMessage{
		Type:        "room_info",
		PlayerCount: room.PlayerCount(),
		PlayerList:  playerListInfo,
		GameState:   room.engine.GetState(),
		Surivals:    room.getCurrentSurivals(),
	}
}

func validPlayerJoin(room *Room, playerUUID string) (*Player, error) {
	player, exist := FindSession(playerUUID)

	if !exist {
		return nil, fmt.Errorf("Player session not exist '%s'", playerUUID)
	}

	if len(room.Players) >= MAX_PLAYERS {
		return player, errors.New("The current game room is full.")
	}

	if _, inRoom := room.FindPlayerByUsername(&player.Username); inRoom {
		return player, fmt.Errorf("A player name '%s' is already in this room", player.Username)
	}

	return player, nil
}

func (room *Room) AddPlayer(playerUUID string) error {
	player, error := validPlayerJoin(room, playerUUID)

	if error != nil {
		return error
	}
	room.Players = append(room.Players, player)
	player.RoomUUID = room.RoomUUID
	room.PublishRoomInfo()
	return nil
}

func (room *Room) PlayerCount() int {
	return len(room.Players)
}

func (room *Room) TryStartGame(playerUUID *string) error {
	player, exist := room.FindPlayerByUUID(playerUUID)
	if !exist || room.OwnerUUID != player.SessionUUID {
		return errors.New("Invalid player")
	}
	if err := room.engine.StartGame(); err != nil {
		return err
	}
	if room.PlayerCount() < MIN_PLAYERS_TO_START {
		return errors.New("Require at least 4 players to start the game.")
	}

	room.PublishRoomInfo()
	err := room.dealCards()
	if err != nil {
		return err
	}

	return nil
}

func (room *Room) dealCards() error {
	for _, player := range room.Players {
		cards, err := room.engine.DealCards(5)
		if err != nil {
			return err
		}
		room.playerCards[player] = cards
	}
	room.PublishPlayerHoldingCards()
	return nil
}

func (room *Room) validGameState() error {
	if room.engine.GetState() != liar.StateInGame {
		return errors.New("Game is not started")
	}

	return nil
}

func (room *Room) PlayerPlaceCard(playerUUID string, cards []liar.Card) error {
	p, err := room.validatePlayerTurn(playerUUID)
	if err != nil {
		return err
	}

	currentPlayerCards := room.playerCards[p]
	remainCards, err := room.engine.PlaceCard(currentPlayerCards, cards)
	if err != nil {
		return err
	}

	msg := message.RoomBoardPlayerPlaceCardMessage{
		Type:     "player_place_cards",
		Username: p.Username,
		Number:   len(cards),
	}

	// After cards placed
	room.LastPlaceCards = cards
	room.playerCards[p] = remainCards

	room.nextRound()

	room.Broadcast(msg)
	room.PublishRoomInfo()
	return nil
}

func (room *Room) PlayerDeclare(playerUUID string, doubt bool) error {
	p, err := room.validatePlayerTurn(playerUUID)
	if err != nil {
		return err
	}

	// If the current player choice to doubt
	lastPlayer := room.Players[room.GetLastPlayerIndex()]
	result, err := room.engine.Declare(doubt)
	if err != nil {
		return err
	}

	msg := message.RoomBoardCastDeclareMessage{
		Refname: p.Username,
		Suspect: lastPlayer.Username,
		Result:  result,
	}
	room.Broadcast(msg)
	if result == liar.Lied {
		err := room.killPlayer(room.Players[room.GetLastPlayerIndex()])
		if err != nil {
			return err
		}
	} else {
		err := room.killPlayer(p)
		if err != nil {
			return err
		}

	}

	room.PublishRoomInfo()
	return nil
}

func (room *Room) killPlayer(p *Player) error {
	p.Alive = false

	if room.getCurrentSurivals() == 1 {
		err := room.endGame()
		// Unexpected game state
		if err != nil {
			return err
		}
	}

	return nil
}

func (room *Room) endGame() error {
	if err := room.engine.EndGame(); err != nil {
		return err
	}
	room.PublishRoomInfo()

	return nil
}

func (room *Room) nextRound() {
	nextPlayerI := func(currentIndex int) int {
		for i := 1; i <= MAX_PLAYERS; i++ {
			nextIndex := (currentIndex + i) % MAX_PLAYERS
			if room.Players[nextIndex].Alive {
				return nextIndex
			}
		}

		// if no surivals
		return -1
	}

	// 更新当前玩家索引
	room.CurrentPlayerIndex = nextPlayerI(room.CurrentPlayerIndex)
}

func (room *Room) GetLastPlayerIndex() int {
	if room.CurrentPlayerIndex == 0 {
		return MAX_PLAYERS - 1
	}

	return room.CurrentPlayerIndex - 1
}

func (room *Room) GetPlayerIndex(p *Player) int {
	for i, rp := range room.Players {
		if rp == p {
			return i
		}
	}

	return -1
}

func (room *Room) validatePlayerTurn(playerUUID string) (*Player, error) {
	p, exist := FindSession(playerUUID)
	if !exist {
		return nil, errors.New("Not existed player")
	}

	if i := room.GetPlayerIndex(p); i == -1 || i != room.CurrentPlayerIndex {
		return nil, errors.New("Not your turn")
	}

	return p, nil
}

func (room *Room) getCurrentSurivals() int {
	return utils.SliceCount(room.Players, func(p *Player) bool {
		return p.Alive == true
	})
}

func (room *Room) GetCurrentAction() liar.GameAction {
	return room.engine.CurrentAction
}

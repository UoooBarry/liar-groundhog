package session

import (
	"fmt"
	"log"
	"sync"
	"uooobarry/liar-groundhog/internal/liar"

	"github.com/google/uuid"
)

type RoomManager struct {
	sync.Mutex
	rooms map[string]*Room
}

var roomManager = &RoomManager{
	rooms: make(map[string]*Room),
}

func (rm *RoomManager) CreateRoom(ownerUUID string, gameEngine *liar.Engine) (*Room, error) {
	rm.Lock()
	defer rm.Unlock()

	owner, exist := FindSession(ownerUUID)
	if !exist {
		return nil, fmt.Errorf("Player session not exist '%s'", ownerUUID)
	}

	uuid := uuid.NewString()
	room := &Room{
		RoomUUID:           uuid,
		engine:             *gameEngine,
		CurrentPlayerIndex: 0,
		OwnerUUID:          owner.SessionUUID,
		BetCard:            BET_CARD,
		Players:            []*Player{owner},
		playerCards:        make(map[*Player][]liar.Card),
	}

	rm.rooms[uuid] = room
	log.Printf("Created room UUID '%s'", uuid)
	return room, nil
}

func (rm *RoomManager) FindRoom(uuid string) (*Room, bool) {
	rm.Lock()
	defer rm.Unlock()
	return rm.rooms[uuid], rm.rooms[uuid] != nil
}


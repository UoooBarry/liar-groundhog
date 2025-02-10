/**
 * Generate login message
 * @param {string} username - Description of the parameter
 * @returns {Object} - Description of the return value
 */
function loginMessage(username) {
    return {
        type: "login",
        username: username
    }
}

function createRoomMessage(user) {
    return {
        type: "room_create",
        sessionuuid: user.sessionuuid,
        username: user.username
    }
}

function joinRoomMessage(user, roomuuid) {
    return {
        type: "room_join",
        sessionuuid: user.sessionuuid,
        username: user.username,
        roomuuid: roomuuid
    }
}

function startRoomMessage(user, roomuuid) {
    return {
        type: "room_start",
        sessionuuid: user.sessionuuid,
        username: user.username,
        roomuuid: roomuuid
    }
}

function start() {
    const people = [
        {
            username: 'barry'
        },
        {
            username: 'tom'
        },
        {
            username: 'joey'
        },
        {
            username: 'alice'
        }
    ]
    const socket = new WebSocket('ws://localhost:8080/ws');
    let roomuuid = "";
    let startGameFlag = false;

    // Connection opened
    socket.addEventListener("open", (event) => {
        console.log("onoepn", event)

        people.forEach((user) => {
            socket.send(JSON.stringify(loginMessage(user.username)))
        })
    });

    // Listen for messages
    socket.addEventListener("message", (event) => {
        console.log("Message from server ", event.data);
        const response = JSON.parse(event.data)

        // a user login success
        if (response.type === "login") {
            let user = people.find((p) => p.username == response.username);
            if (user) {
                user.sessionuuid = response.sessionuuid;
                console.log('login success:', user);
            } else {
                console.error("unfond local user")
            }
        }
        // room_create success
        if (response.type === 'room_create') {
            people[0].roomuuid = response.roomuuid
            roomuuid = response.roomuuid
            people.slice(1).forEach((p) => {
                socket.send(JSON.stringify(joinRoomMessage(p, roomuuid)));
                p.uuid = roomuuid;
            })
        }

        if (response.type === 'room_info' && response.game_state === 'preparing' && response.player_count === 4 && !startGameFlag) {
            console.log('room is full, ready to start the game')
            startGameFlag = true
            socket.send(JSON.stringify(startRoomMessage(people[0], roomuuid)))
        }

        // if all logged in
        if (people.every((p) => p.sessionuuid && p.sessionuuid !== "" && !p.roomuuid)) {
            // barry create the room
            console.log(createRoomMessage(people[0]))
            socket.send(JSON.stringify(createRoomMessage(people[0])))
        }

        if (response.type === 'error') {
            console.error(response.content);
            socket.close();
            return;
        }
    });
}

start()

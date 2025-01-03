package session

import "sync"

var rooms = struct {
	sync.Mutex
	data map[string]Session // uuid -> Session
}{
	data: make(map[string]Session),
}

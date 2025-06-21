package pkg

import "time"

type LobbyID int

type Lobby struct {
	id       LobbyID
	host     PeerID
	members  map[PeerID]*Peer
	sealedAt time.Time
}

func NewLobby(host *Peer) *Lobby {
	host.id = 1;
	return &Lobby{
		id:      LobbyID(NewID()),
		host:    1,
		members: make(map[PeerID]*Peer),
	}
}

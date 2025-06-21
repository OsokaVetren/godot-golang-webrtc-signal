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
	L := &Lobby{
		id:      LobbyID(NewID()),
		host:    host.id,
		members: make(map[PeerID]*Peer),
	}
	L.members[host.id] = host;
	return L;
}

package pkg

import "time"

type Lobby struct {
	id       LobbyId
	hostGlobal     PeerId
	members  map[Local Id]*Peer
	global2local map[PeerId]LocalId
	local2global map[LocalId]PeerId
	sealedAt time.Time
}

func NewLobby(host *Peer) *Lobby {
	L := &Lobby{
		id:      LobbyId(NewID()),
		hostGlobal:    host.id,
		members:      make(map[ID]*Peer),
    		global2local: make(map[PeerId]Local Id),
    		local2global: make(map[Local Id]PeerId),
	}
	L.members[1] = host;
	L.global2local[host.id] = 1
	L.local2global[1] = host.id
	return L;
}

func (L *Lobby) AddMember(p *Peer) Local Id {
  var next Local Id = 2
  for {
    if _, ok := L.members[next]; !ok {
      break
    }
    next++
  }
  L.members[next] = p
  L.global2local[p.id] = next
  L.local2global[next] = p.id
  return next
}

func (L *Lobby) Local Id(p *Peer) Local Id {
	return L.global2local[p.id] 
}

func (L *Lobby) PeerByLocal(id Local Id) *Peer {
	return L.members[id] 
}

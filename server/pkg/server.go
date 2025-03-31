package pkg

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

const (
	READ_TIMEOUT  = 60 * time.Second
	WRITE_TIMEOUT = 10 * time.Second
	// Interval for sending ping messages to WebSocket connections. Needs to be
	// less than the read timeout for the connection.
	PING_INTERVAL = (READ_TIMEOUT * 9) / 10
)

var (
	NEW_LINE = []byte{'\n'}
)

type Server struct {
	hub *Hub
}

// Initialize a Hub and start listening for messages
func (server *Server) Run() {
	server.hub = NewHub()
	go server.hub.Run()
}

// Initialize a Peer for a WebSocket connection.
func (server *Server) InitPeer(ws *websocket.Conn) {
	peer := NewPeer(ws)

	go peerToHub(peer, server.hub)
	go peerToWs(peer)
}

// Next read deadline for the WebSocket connection is now + READ_TIMEOUT.
func nextDeadline(timeout time.Duration) time.Time {
	return time.Now().Add(timeout)
}

// Connect a peer to the hub and pump WebSocket messages from peer to hub.
// Disconnect from hub and close WebSocket when done.
func peerToHub(peer *Peer, hub *Hub) {
	fmt.Println("[Server] peerToHub")

	defer func() {
		fmt.Println("[Server] peerToHub exiting")
		hub.disconnect <- peer
		peer.ws.Close()
	}()

	hub.connect <- peer

	// Set read deadline for the WebSocket connection and reset it every time a
	// message is received.
	peer.ws.SetReadDeadline(nextDeadline(READ_TIMEOUT))
	peer.ws.SetPongHandler(func(string) error {
		peer.ws.SetReadDeadline(nextDeadline(READ_TIMEOUT))
		return nil
	})

	fmt.Println("[Server] Starting message loop")

	for {
		_, msg, err := peer.ws.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		fmt.Println("[Server] msg:", string(msg))
		hub.peer_msg <- NewPeerMsg(peer.id, msg)
	}
}

// Sends messages from the Peer.send channel to the WebSocket connection.
// The Hub sends messages to the Peer.send channel.
func peerToWs(peer *Peer) {
	pinger := time.NewTicker(PING_INTERVAL)
	defer func() {
		pinger.Stop()
		peer.ws.Close()
	}()

	for {
		select {
		case msg, ok := <-peer.send:
			if !ok {
				fmt.Println("[Server] Peer send channel closed")
				peer.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// It's possible we have queued messages waiting, so flush all of
			// them.
			for {
				peer.ws.SetWriteDeadline(nextDeadline(WRITE_TIMEOUT))

				err := peer.ws.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					fmt.Println("[Server] Error writing message:", err)
					return
				}

				if len(peer.send) == 0 {
					break
				}

				msg = <-peer.send
			}

		// Ping the WebSocket connection to keep it alive. Update the write timeout.
		case <-pinger.C:
			peer.ws.SetWriteDeadline(nextDeadline(WRITE_TIMEOUT))
			if err := peer.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Println("[Server] Error writing ping message:", err)
				return
			}
		}
	}
}

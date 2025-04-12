# Golang WebRTC Signaling Server

WebRTC signal server with SSL support written in Go. Supports hosting and joining lobbies. Includes Godot client.

## Project Architecture

- /client (`SignalWsClient` class) - Godot client code. Provides a `SignalLobby` node via a Godot plugin.
  - `SignalLobby` node (signal_lobby_node.gd and signal_lobby_node.tscn) - Handles signals from `SignalWsClient` class and creates a WebRTC mesh using Godot's multiplayer apis
  - `SignalWsClient` class (scripts/client.gd) - manages the client websocket connection and a `SignalWsPeer` instance
  - `SignalWsMsg` class (scripts/msg.gd) - defines message types and handles serialization / deserialization of messages.
  - `SignalWsPeer` class (scripts/peer.gd) -
- /server - Golang signaling server

### Messages

The server and client communicate over a websocket connection to establish a peer to peer WebRTC connection. Messages are identified by a numeric enum type:

```gdscript
enum Type {
	INVALID = 0,
	CONNECTED,
	HOST,
	JOIN,
	PEER_CONNECT,
	PEER_DISCONNECT,
	OFFER,
	ANSWER,
	CANDIDATE,
	SEAL
}
```

Messages are passing in a simple serialized format in the form of `type|pid|data`:

- type - numeric string for the msg type
- pid - peer id identifiying a peer in the WebRTC mesh
- data - optional data used by some message type

## Installation

### Client

The client code can be installed as a Godot plugin

1. Copy the `/client` folder to the `/addons` folder in a Godot project
1. Enable the `Godot Golang WebRTC Signal` plugin in the Plugins tab of Project settings

   ![Plugin Settings](docs/plugin-settings.png)

1. The plugin provides a `SignalLobby` node that can be added to a scene

   ![SignalLobby](docs/signal-lobby.png)

1. After adding the node to a scene, you can right-click on the node and select `% Access as Unique Name`. The node should now show a `%` character.

   ![SignalLobby % Access as Unique Name](docs/signal-lobby-node.png)

1. You can now connect to the `player_added` and `lobby_sealed` signals:

   ```gdscript
   @onready var signal_lobby = %SignalLobby # SignalLobby node set to "% Access as Unique Name"

   func _ready() -> void:
      signal_lobby.player_added.connect(_on_player_added)
      signal_lobby.lobby_sealed.connect(_on_lobby_sealed)
   ```

## Systemd Config

The server can be configured via systemd on a linux server such as a Raspberry Pi.

1. Open a new config file for the service:

   ```sh
   sudo nano /etc/systemd/system/gowebrtcsignal.service
   ```

1. Copy / paste the following template filling in <USER> and <APPDIR> placeholders

   ```ini
   [Unit]
   Description = Go WebRTC Signaling Server

   [Service]
   Type             = simple
   Restart          = always
   RestartSec       = 5s
   StandardOutput   = append:<APPDIR>/stdout.log
   StandardError    = append:<APPDIR>/stderr.log
   ExecStart        = <APPDIR>/signalserver
   WorkingDirectory = <APPDIR>
   User             = <USER>

   [Install]
   WantedBy = multi-user.target
   ```

1. Enable and start the service

```sh
sudo systemctl enable gowebrtcsignal.service
sudo systemctl start gowebrtcsignal.service
sudo systemctl status gowebrtcsignal.service
```

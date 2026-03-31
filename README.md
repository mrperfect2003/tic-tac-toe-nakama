# Tic Tac Toe Nakama Backend

A multiplayer Tic Tac Toe backend built using Nakama (server-authoritative architecture) and Golang.  
This project demonstrates real-time match handling, RPC-based match creation, and game state synchronization.

---

## Author

Keshav Raj  
LinkedIn: https://www.linkedin.com/in/keshavraj18  
Portfolio: https://mr-keshav.netlify.app/

---

## Overview

This project implements a real-time multiplayer Tic Tac Toe game using Nakama.  
The server controls all game logic, ensuring fairness, consistency, and synchronization between players.

---

## Features

- Server-authoritative game logic  
- Real-time multiplayer match handling  
- Turn-based gameplay (X and O)  
- Winner and draw detection  
- Game state synchronization using dispatcher  
- Custom RPC for match creation  
- Fully containerized using Docker  

---

## Tech Stack

- Backend: Golang  
- Game Server: Nakama  
- Database: CockroachDB  
- Containerization: Docker, Docker Compose  
- Communication: WebSocket / gRPC (via Nakama)  

---

## Project Structure

```

tic-tac-toe-nakama/
│
├── api/
│   └── rpc.go              # RPC functions (create match)
│
├── match/
│   ├── match.go            # Match lifecycle (Init, Join, Loop, etc.)
│   ├── logic.go            # Game logic (winner detection)
│   └── state.go            # Game state definition
│
├── main.go                 # Entry point (register match and RPC)
├── Dockerfile.build        # Plugin build configuration
├── docker-compose.yml      # Nakama and database setup
├── go.mod / go.sum         # Dependencies
└── README.md

```

---

## Requirements

Make sure the following are installed:

- Docker  
- Docker Compose  
- Go (optional, for local development)

---

## How to Run

### 1. Clone Repository

```

git clone [https://github.com/mrperfect2003/tic-tac-toe-nakama](https://github.com/your-username/tic-tac-toe-nakama.git)
cd tic-tac-toe-nakama

```

### 2. Build Nakama Plugin

```

docker build -f Dockerfile.build -t nakama-plugin-builder .

```

### 3. Start Services

```

docker compose up --build

```

---

## Access Services

- Nakama API: http://localhost:7350  
- Nakama Console: http://localhost:7351  

---

## How It Works

```

Open Nakama Console

Go to:

http://localhost:7351

Default login:

username: admin
password: password

### Match Creation

Client calls RPC:

```

create_match

```

Server creates a match and returns:

```

{
"match_id": "some-match-id"
}

```

---

### Gameplay

Players send moves in the following format:

```

{
"row": 0,
"col": 2
}

```

---

### Server Responsibilities

- Validate player turn  
- Validate move position  
- Update board state  
- Check for winner or draw  
- Broadcast updated state to all players  

---

## Match Lifecycle

- MatchInit: Initialize game state  
- MatchJoinAttempt: Validate player join  
- MatchJoin: Add players to match  
- MatchLoop: Process moves and update game  
- MatchLeave: Handle player leaving  
- MatchTerminate: Cleanup resources  

---

## Game State

The game state includes:

- 3x3 board  
- List of players  
- Current turn  
- Winner  
- Move count  
- Game status  

---

## Real-time Communication

The server uses:

```

dispatcher.BroadcastMessage(...)

```

to send updated game state to all connected players.

---

## Notes

- CockroachDB is running in insecure mode for development  
- This project is intended for learning and demonstration purposes  

---

## Future Improvements

- Add matchmaking  
- Add authentication  
- Add leaderboard  
- Build frontend client (React or Unity)  
- Add reconnect support  

---

## Support

If you found this project useful:

- Star the repository  
- Fork the project  
- Contribute improvements  

---

## License

This project is open for learning and educational use.
```

---

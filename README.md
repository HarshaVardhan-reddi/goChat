# goChat

Real-time one-to-one chat server built with Go, Gorilla WebSocket, MySQL, and Redis. Designed around a concurrent connection hub, event-driven message routing, and Redis pub/sub for cross-instance status tracking.

## Architecture

```
┌────────┐  WS Upgrade   ┌──────────────────────────────────────────────┐
│ Client ├───────────────►│  Server (:3000)                             │
└────────┘  JWT Auth      │                                             │
                          │  ┌─────┐    ┌──────────────┐                │
                          │  │ Hub │◄──►│ WsConnection │ (per user)     │
                          │  └──┬──┘    │  ├ broker chan│                │
                          │     │       │  ├ ctx/cancel │                │
                          │     │       └──────┬───────┘                │
                          │     │         ┌────┴────┐                   │
                          │     │    ┌────┤Goroutines├────┐             │
                          │     │    │    └─────────┘    │             │
                          │     │    ▼                   ▼             │
                          │     │  MessageReader    MessageWriter      │
                          │     │    │                   ▲             │
                          │     │    ▼                   │             │
                          │     │  EventProcessor       │             │
                          │     │    ├─ MESSAGE ────► target.broker    │
                          │     │    ├─ SUBSCRIBE ──► Redis Sub       │
                          │     │    └─ UNSUBSCRIBE                   │
                          │     │                                      │
                          │  ┌──┴───────────┐                          │
                          │  │    Redis      │                          │
                          │  │  pub/sub per  │                          │
                          │  │  status:{uid} │                          │
                          │  └──────────────┘                          │
                          └──────────────────────────────────────────────┘
```

### Request Flow

1. Client authenticates via `POST /api/v1/users/login` or `/signup` and receives a JWT
2. Client opens a WebSocket at `WS /api/v1/ws/startchat` with the JWT
3. Server upgrades the connection, registers it in the Hub (`userID -> WsConnection`)
4. Two goroutines spin up: `MessageReader` (inbound) and `MessageWriter` (outbound)
5. Inbound messages are deserialized into a `WsEvent` envelope and dispatched by the `EventProcessor`
6. On connect/disconnect, the server publishes `ACTIVE`/`INACTIVE` to the Redis channel `status:{userId}`

## Design Decisions

| Decision | Rationale | Trade-off |
|----------|-----------|-----------|
| **`sync.Map` for Hub** | Lock-free reads for a read-heavy workload (message routing lookups far exceed connect/disconnect writes) | Not optimal for write-heavy scenarios; would need sharded maps or a different structure at very high churn |
| **Broker channel per connection** | Decouples message production from WebSocket writes. The reader goroutine never blocks on a slow client — it drops the event into the channel and moves on | Bounded buffer (1024) means backpressure is possible under extreme load; requires monitoring |
| **Two goroutines per connection** | Separates read and write paths for non-blocking I/O. Clean shutdown via `context.WithCancel` — cancelling the context tears down both goroutines and the Redis subscription | Goroutine cost per connection (~4KB stack each); acceptable for the expected connection count |
| **Redis pub/sub for status** | Enables status events to propagate across multiple server instances, not just in-process. Each user's status gets its own channel (`status:{userId}`) for targeted subscriptions | Adds Redis as an infrastructure dependency; pub/sub is fire-and-forget (no persistence of missed status events) |
| **Polymorphic JSON with explicit `ToJSON`/`FromJSON`** | Go's `encoding/json` recurses infinitely when a type implements `json.Marshaler` and also gets marshaled as a field. Explicit methods on the `Message` interface avoid this while keeping the `WsEvent` envelope generic | Slightly more boilerplate than using struct tags alone |
| **Interface-driven repositories** | `UserRepository` interface lets controllers remain testable and decoupled from the MySQL/GORM implementation | Minor indirection cost; pays off when adding test doubles or swapping storage |

## Event System

All messages flow through a `WsEvent` envelope. The `details` field holds a polymorphic payload dispatched by `EventType`:

```json
{
  "from": { "id": 1 },
  "to": { "id": 2 },
  "src": 2,
  "details": {
    "type": 1,
    "message": { "..." }
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

| EventType | Code | Payload | Direction |
|-----------|------|---------|-----------|
| `MESSAGE` | 1 | `ChatMessage` (text + attachments) | Client → Server → Client |
| `SUBSCRIBE` | 2 | `SubscribeMessage` (target user ID) | Client → Server |
| `UNSUBSCRIBE` | 3 | `UnsubscribeMessage` (target user ID) | Client → Server |
| `STATUS_UPDATE` | 4 | `StatusMessage` (user ID, status, last seen) | Server → Client |

## Project Structure

```
├── main.go                          # Entry point — loads config, initializes DB, starts HTTP server
├── configs/
│   ├── router.go                    # Gorilla Mux routes, CORS middleware
│   └── databases/                   # MySQL config loader, GORM initialization
├── internals/
│   ├── auth/                        # JWT handling, credential & password validation
│   ├── controllers/v1/
│   │   ├── users/                   # Signup, login endpoints
│   │   └── chats/                   # WebSocket upgrade endpoint
│   ├── models/
│   │   ├── ws_event.go              # WsEvent envelope, polymorphic marshaling
│   │   └── events/                  # ChatMessage, StatusMessage, SubscribeMessage, UnsubscribeMessage
│   ├── services/
│   │   ├── event_processor.go       # Routes events by type to handlers
│   │   ├── redis.go                 # Redis singleton, pub/sub helpers
│   │   ├── chats/                   # MessageReader/MessageWriter goroutines
│   │   └── users/                   # Login, signup, token services
│   ├── repositories/                # UserRepository interface + MySQL implementation
│   └── helpers/                     # HTTP response utilities
├── client/                          # Test client (Node.js WebSocket client)
└── db/migrator/                     # Database migration tool
```

## Getting Started

### Prerequisites

- Go 1.21+
- MySQL on `localhost:3306`
- Redis on `localhost:6379`

### Setup

```bash
# Clone and navigate
git clone <repo-url> && cd chat1-1

# Configure environment
cp .env.example .env   # Edit with your DB credentials

# Run database migrations
./dbmigrator

# Start the server
go run main.go         # Listens on :3000

# (Optional) Run the test client
cd client && npm install && node test_client.js
```

### API

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/users/signup` | Create account |
| POST | `/api/v1/users/login` | Authenticate, returns JWT |
| WS | `/api/v1/ws/startchat` | Upgrade to WebSocket (JWT required) |

## Roadmap

- [ ] **Message persistence** — Store messages in MongoDB/MySQL with delivery status tracking (sent/delivered/read)
- [ ] **Message history API** — Fetch paginated conversation history via REST or WebSocket command
- [ ] **Offline message queue** — Buffer messages for offline users and deliver on reconnect
- [ ] **Delivery receipts** — Read/delivered acknowledgments propagated back to sender
- [ ] **Horizontal scaling** — Multiple server instances behind a load balancer, with Redis handling cross-instance message routing
- [ ] **Rate limiting** — Per-connection message throttling to prevent abuse

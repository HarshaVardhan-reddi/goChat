# Chat Application Todo List

## Feat 1: Send Message Functionality (One-to-One Chat)
- [ ] **WebSocket Manager/Hub**: Create a central manager to handle client connections, registrations, and deregistrations.
- [ ] **Message Protocol**: Define a JSON structure for messages (e.g., sender_id, receiver_id, content, timestamp).
- [ ] **Server-side Routing**:
    - [ ] Map authenticated users to their WebSocket connections.
    - [ ] Implement logic to find the destination user's connection and deliver the message.
    - [ ] Handle cases where the destination user is offline (queue or just store in DB).
- [ ] **Concurrency**: Ensure thread-safe access to the connection map.

## Feat 2: Active Status Tracking
- [ ] **Status Management**: Update user's "online" status in the database or an in-memory store upon connection/disconnection.
- [ ] **Real-time Updates**:
    - [ ] Broadcast a user's online status change to their contacts or interested parties.
    - [ ] Implement an endpoint or WebSocket event to fetch initial online status of users.

## Feat 3: Persistent Messages & History
- [ ] **Database Schema**: 
    - [ ] Design a `Message` collection in MongoDB (preferred for chat logs) or a table in MySQL.
    - [ ] Fields: `id`, `sender_id`, `receiver_id`, `content`, `status` (sent/delivered/read), `timestamp`.
- [ ] **Persistence Logic**:
    - [ ] Save every incoming message to the database before (or after) delivery.
    - [ ] Implement message status updates (e.g., mark as 'read' when the receiver opens the chat).
- [ ] **History Retrieval**:
    - [ ] Implement an API endpoint or WebSocket command to fetch the latest N messages for a specific conversation.
    - [ ] Ensure messages are loaded when a user joins a chat session.

## Miscellaneous
- [ ] **Authentication**: Ensure WebSocket connections are authenticated (e.g., via JWT).
- [ ] **Error Handling**: Robust handling of connection drops, invalid message formats, and database failures.

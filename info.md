Project MVP Timeline - 

Day 1  -

- Set up the development environment.

- Create initial project structure for both frontend and backend.


Day 2: Backend - Core Distributed Logic

- Implement WebSocket server setup in Go, Set up the basic structure for handling client connections and messages.

- Implement the Document State Service : State Synchronizer, Conflict Resolution module using message queues

- Implement the Persistence Layer Set up MongoDB Implement basic data storage and retrieval functions.


Day 3: Backend - Advanced Features

- Implement the Replication Service Ensure document changes are replicated across multiple nodes with consistency.

- Implement the Recovery module Implement crash recovery mechanisms Ensure data integrity and consistency post-recovery.

- IDK IF required but - Set up the Load Balancer Distribute client connections across multiple WebSocket servers. Handle failover and load balancing logic.

Day 4: Frontend - Basic Implementation

- Set up the React/TS frontend Basic UI components for document editing.

- Implement WebSocket client in the frontend Handle real-time updates. Display changes from other users.

- Integrate URL TO new session logic - check [notepad.pw](https://notwpad.pw/).

Day 5: Frontend - Advanced Features

- Implement collaborative editing features - Real-time cursor position updates User presence indicators.

- Implement caching/migration strategies to minimize response time. Frontend caching mechanisms. Optimize data fetching and synchronization.

- Optional(Linked List browser -> stores webpage in back button) - Implement additional UI features Document history and version control.

Day 6: Integration and Testing

- Integrate frontend and backend.
- Conduct end-to-end testing.
- Real-time collaboration scenarios.
- Crash recovery scenarios.
- Load balancing and failover scenarios.

Day 7: Finalization and Deployment

- Fix any bugs identified during testing. Optimize performance and scalability.
- Prepare documentation.
- Deploy the application.
- Ensure all services are running and stable.
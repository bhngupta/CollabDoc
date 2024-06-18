# CollabDoc


CollabDoc is a real-time collaborative document editor built with Go. It allows multiple users to edit a document simultaneously while maintaining consistency and handling conflicts seamlesly.

Every wanted a paper napkin on the internet to quickly jot down suff? Look no further, Collab Doc is a just the right tool. It can  be shared to multiple people to collaborate and also be saved for later use! (upto 90 days)

## Features 

- Real-time Collaboration: Multiple users can edit the document simultaneously, with changes reflected in real-time.
- Simple Persistence: Document state is saved to a JSON file, allowing for easy state recovery.
- WebSocket Communication: Utilizes WebSockets for low-latency, bi-directional communication between.
- Open Source.

## The WHY?

The project was sparked by the idea of having a paper napkin on internet for jotting down ideas and sharing them effortlessly. Drawing inspiration from Google Docs, our goal was to create a seamless experience without the necessity of a login layer, while preserving essential features such as real-time collaboration and updates.

## Technology Stack

+ Backend (Go)
  + Concurrency: Goroutines and channels make it easy to handle multiple connections and concurrent operations efficiently
  + Performance: Compiled nature and efficient memory management ensure low-latency and high-throughput
  + Standard Library: Ideal for scalable WebSocket server development.

+ Frontend (React Typescript)
  + Leveraged declarative component-based architecture, facilitating UI updates in real-time with WebSocket API

+ Persistence (JSON file): 
  + Simple data storage due to its lightweight nature and ease of integration for file operations

## Architecture

![CollabDoc Architecture Diagram](https://github.com/bhngupta/CollabDoc/blob/main/misc/arch-diagram.png?raw=true)

## Core Concepts 

1. Collaborative Editing

- Concurrency Management: Go's concurrency model with goroutines and channels is leveraged to manage simultaneous editing sessions efficiently. Each user's edits are processed concurrently, allowing for smooth collaboration without blocking operations.

- Operational Transformation (OT): CollabDoc utilizes Operational Transformation (OT) to enable real-time collaboration among multiple users. OT ensures that concurrent edits to a document are handled seamlessly, maintaining consistency and resolving conflicts. This transformation ensures that all users see a consistent view of the document, even when edits overlap in time.

2. Broadcasting Changes

- WebSocket Communication: CollabDoc uses WebSocket technology to facilitate low-latency, bi-directional communication between the server and clients. When a user makes changes to a document, the server broadcasts these changes to all connected clients in real-time. This broadcasting mechanism ensures that updates are immediately visible to all collaborators, enhancing the collaborative editing experience.

By leveraging Go's concurrency features and OT for conflict resolution, coupled with WebSocket broadcasting for real-time updates, CollabDoc provides a robust platform for collaborative document editing without relying on vector clocks.


## Getting Started

### Prerequisites

Ensure you have the following installed:

1. **Git**
2. **Go** (1.22+)
3. **Node.js** (20+)
4. **NPM** (10+)

### Install from source

#### Backend Setup

0. Fork the repository

1. Clone the repository:

   ```bash
   git clone https://github.com/your_username/CollabDoc
   cd CollabDoc
   ```

2. Install dependencies:

  ```bash
  go mod tidy
  ```

3. Run the server:

  ```bash
  go run cmd/server/main.go
  ```

#### Frontend Setup

1. Navigate to the frontend directory:

  ```bash
  cd frontend
  ```

2. Install dependencies:

  ```bash
  npm install
  ```

3. Start development client:

  ```bash
  npm run dev
  ```

Start collaborating in real-time!

#### Configurations

- Client Port: By default, client runs on port 3000.
- Server Port: By default, backend runs on port 8080.
- WebSocket Endpoint: WebSocket server listens on `/ws`.

### Usage

1. Open your browser and navigate to http://localhost:3000.
2. Create or join a document by entering a document ID.
3. Start collaborating in real-time!

## Future Plans

- [x] ~~Implement broadcasting to clinets on same document ID~~
- [x] ~~Handle multiple clients concurrently using multi-threading or an event loop~~
- [ ] create new arch diagram
- [ ] Write integration tests 
- [ ] Implement a more robust persistence layer - MongoDB 
- [ ] Extend the editor to support rich text and code formatting
- [ ] Delete data after set time

## License

This project is licensed under the [MIT License](https://opensource.org/license/MIT) - see the [LICENSE](https://github.com/bhngupta/CollabDoc/blob/main/LICENSE) file for details.


## Contributing

Contributions are welcome! Please fork the repository and create a pull request with your changes.

### Contributors 

<table>
	<tbody>
		<tr>
      <td align="center">
          <a href="https://github.com/bhngupta" >
              <img src="https://avatars.githubusercontent.com/u/44861163?v=4" width="100;" alt="Bhanu"/>
              <br />
              <sub>Bhanu Gupta</sub>
          </a>
      </td>
      <td align="center">
          <a href="https://github.com/PranavN1234">
              <img src="https://avatars.githubusercontent.com/u/44135759?v=4" width="100;" alt="Pranav"/>
              <br />
              <sub >Pranav Iyer</sub>
          </a>
      </td>
		</tr>
	<tbody>
</table>

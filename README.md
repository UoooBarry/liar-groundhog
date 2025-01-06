# Liar Groundhog

Liar Groundhog is a work-in-progress project that aims to create a duplicate version of the popular "Liar's Bar" experience using modern web technologies. The project combines a Go backend with a Vue frontend to deliver a seamless and engaging platform for users.

## Overview
Liar Groundhog replicates the functionality and design of Liar's Bar, focusing on:
- A fun, interactive environment for users to participate in games and activities.
- Real-time communication powered by WebSockets.
- A scalable and maintainable architecture.

## Features
### Backend (Go):
- WebSocket support for real-time communication.

### Frontend (Vue):
- Dynamic and responsive user interface.
- Real-time updates reflecting the game's state.
- Smooth user interactions for joining games, playing, and viewing results.

## Work in Progress
The project is still under development, and the following areas are actively being worked on:
- Game logic implementation.
- Enhanced error handling and validation for player actions.
- Comprehensive testing for backend and frontend components.
- UI/UX improvements for a more immersive experience.
- Documentation and deployment setup.

## Getting Started
### Prerequisites
- Go 1.20 or later.
- Node.js with pnpm installed.
- Docker for containerized development.

### Development Setup
1. Clone the repository:
   ```bash
   git clone https://github.com/UoooBarry/liar-groundhog.git
   cd liar-groundhog
   ```
2. Start the development environment using Docker Compose:
   ```bash
   make dev
   ```
3. Access the Vue frontend at `http://localhost:5173` and the Go backend at `http://localhost:8080`.

### Running Tests
Run the backend tests:
```bash
make go-test
```

# FitClassMaster
FitClassMaster provides gyms, trainers, and members with a structured digital platform for organizing fitness classes and tracking workout progress.

FitClassMaster is a comprehensive digital platform designed for gyms, trainers, and fitness enthusiasts. It provides a structured environment for organizing fitness classes, managing workout plans, tracking session progress in real-time using WebSockets, and facilitating communication between trainers and members.

## 🚀 Features

- **User Management**: Multi-tier access control for Members, Trainers, and Admins.
- **Class Scheduling**: Trainers can create, edit, and manage fitness classes. Members can browse and enroll in classes.
- **Workout Planning**: Create detailed workout plans with specific exercises, sets, reps, and notes.
- **Real-time Workout Tracking**: Perform workouts with a live tracker. Trainers can monitor member progress in real-time via WebSockets.
- **Messaging System**: Built-in chat system for direct communication between trainers and members.
- **Dashboard**: Personalized views for members and trainers to track schedules and active sessions.
- **Modern UI**: Built with Go Templates, and HTMX for a smooth, single-page application feel without the complexity of a heavy frontend framework.

## 🛠 Tech Stack

- **Backend**: [Go](https://go.dev/) (Golang)
- **Router**: [Chi](https://github.com/go-chi/chi)
- **Database**: [GORM](https://gorm.io/) (Object-Relational Mapper)
- **Real-time**: [WebSockets](https://github.com/gorilla/websocket)
- **Session Management**: [Gorilla Sessions](https://github.com/gorilla/sessions)
- **Frontend**: Go Templates, [HTMX](https://htmx.org/)

## 🏃 Getting Started

### Prerequisites

- Go 1.22 or higher
- A SQL Server instance (or modify the driver in `internal/config/db.go` for other databases)

### Environment Variables

The application requires the following environment variables:

- `DB_DSN`: Database connection string (e.g., `sqlserver://username:password@localhost:1433?database=FitClassMaster`)
- `SESSION_KEY`: A secret key for session encryption.
- `DEV_TEMPLATES`: (Optional) Set to `1` to enable hot-reloading of HTML templates.

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/FitClassMaster.git
   cd FitClassMaster
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the application:
   ```bash
   go run cmd/fitclassmaster/main.go
   ```

The server will start at `http://localhost:8080`.

## 📂 Project Structure

- `cmd/`: Main entry point of the application.
- `internal/config/`: Database and session configuration.
- `internal/handlers/`: HTTP request handlers.
- `internal/middlewares/`: Custom HTTP middlewares (Auth, Roles, etc.).
- `internal/models/`: GORM database models.
- `internal/repositories/`: Database abstraction layer.
- `internal/services/`: Business logic layer.
- `internal/templates/`: HTML templates and rendering logic.
- `internal/websockets/`: WebSocket hub and connection logic.
- `internal/static/`: Static assets (CSS, JS).

## 📄 License

This project is licensed under the GPL-3.0 License.

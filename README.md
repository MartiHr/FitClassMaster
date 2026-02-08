# [FitClassMaster](https://fitclass-app.onrender.com/)

**FitClassMaster** is a comprehensive web-based platform designed to bridge the gap between gym management, trainer planning, and member performance. It moves fitness tracking away from scattered spreadsheets into a single, real-time ecosystem, featuring class scheduling, detailed workout planning, and live session monitoring.

🌐 **Live Demo:** [https://fitclass-app.onrender.com/](https://fitclass-app.onrender.com/)

---

## 📋 Functionality

The system operates on a role-based access model, providing specific tools for three distinct user groups:

### 1. Members
* **Personal Dashboard:** View active subscriptions, upcoming classes, and workout history.
* **Class Enrollment:** Browse the schedule and book spots in group fitness classes.
* **Workout Execution:** A mobile-friendly interface to start workout sessions and log data (reps, sets, weight) in real-time.
* **Communication:** Built-in chat system to message trainers directly.

### 2. Trainers
* **Workout Plan Management:** Create and edit detailed workout templates with specific exercises, ordering, and instructions.
* **Live Tracking:** Monitor active client sessions in real-time. Using WebSockets, the trainer sees every set completed by the client instantly without refreshing the page.
* **Class Management:** Schedule, modify, or cancel group classes and view attendee lists.

### 3. Admins
* **User Management:** Oversee all registered users and manage role assignments (Member/Trainer/Admin).
* **System Oversight:** Full access to all platform data for maintenance and support.

---

## 🏗 Architecture

The project follows a **Layered Architecture** to ensure Separation of Concerns, maintainability, and testability.

1.  **Presentation Layer (Handlers & Templates):**
    * Handles HTTP requests and renders HTML using **Go Templates** (Server-Side Rendering).
    * Utilizes **HTMX** for dynamic, partial page updates, providing a Single-Page Application (SPA) feel without heavy JavaScript frameworks.
2.  **Business Logic Layer (Services):**
    * Contains the core rules of the application (e.g., class capacity validation, session state management).
    * Isolates business logic from the database and HTTP transport layers.
3.  **Data Access Layer (Repositories):**
    * Abstracts direct database interactions using **GORM**.
    * Handles all CRUD operations and queries.

---

## 💻 Tech Stack

### Backend
* **Language:** [Go (Golang) 1.22+](https://go.dev/)
* **Router:** [Chi](https://github.com/go-chi/chi) - A lightweight, idiomatic router for Go.
* **ORM:** [GORM](https://gorm.io/) - For object-relational mapping.
* **Real-time:** [Gorilla WebSocket](https://github.com/gorilla/websocket) - For bidirectional communication.
* **Sessions:** Gorilla Sessions - For secure user session management.

### Frontend
* **Template Engine:** Go `html/template`.
* **HTMX:** For AJAX requests and dynamic DOM manipulation.
* **CSS:** Custom styling with FontAwesome for icons.

### Infrastructure & DevOps
* **Database:** PostgreSQL.
* **Containerization:** Docker & Docker Compose.
* **Cloud Platform:** Render.com (Auto-deploy via Git).

---

## 🗂 Data Model

The application uses a relational database (PostgreSQL) with the following key entities:

* **Users:** Stores profile data, password hashes, and roles (Admin, Trainer, Member).
* **Classes:** Represents group sessions with schedules, capacity, and assigned trainers.
* **Enrollments:** A Many-to-Many relationship table between Users and Classes.
* **WorkoutPlans:** Templates for structured workouts created by trainers.
* **WorkoutExercises:** Defines the specific exercises, order, and targets within a plan.
* **Sessions:** Represents an active or completed workout instance.
* **SessionLogs:** Granular data for every set performed (weight/reps) linked to a session.
* **Messages & Conversations:** Stores chat history between users.

---

## ⚙️ Configuration & Installation

### Prerequisites
* **Docker Desktop** (Recommended) OR **Go 1.22+** and a local **PostgreSQL** instance.

### Option 1: Running with Docker (Recommended)
This method automatically builds the application and sets up the database in isolated containers.

1.  Clone the repository:
    ```bash
    git clone [https://github.com/your-username/FitClassMaster.git](https://github.com/your-username/FitClassMaster.git)
    cd FitClassMaster
    ```
2.  Build and start the containers:
    ```bash
    docker compose up --build
    ```
3.  Access the application at: `http://localhost:8080`

*To stop the application:* `docker compose down`

### Option 2: Local Manual Setup
1.  Create a `.env` file in the root directory:
    ```env
    PORT=8080
    DB_DSN="host=localhost user=postgres password=yourpass dbname=fitclass port=5432 sslmode=disable"
    SESSION_KEY="your-secret-key"
    ```
2.  Install dependencies:
    ```bash
    go mod tidy
    ```
3.  Run the application:
    ```bash
    go run main.go
    ```
---

## 📚 References

* [Go Programming Language Documentation](https://go.dev/doc/)
* [GORM Guides](https://gorm.io/docs/)
* [HTMX Documentation](https://htmx.org/docs/)
* [Chi Router](https://github.com/go-chi/chi)
* [Gorilla WebSocket](https://github.com/gorilla/websocket)

<p align="center">
<img src="https://img.shields.io/badge/kood%2FSisu-Car_viewer-brightgreen?logo=gitea&logoColor=white&labelColor=8A2BE2">
<img src="https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white">
<img src="https://img.shields.io/badge/style-Red%20Cars-ff0055">
</p>

# Car Viewer

A robust, server-side rendered web application designed to browse car catalog, compare vehicle specifications, and provide personalized recommendations based on user behavior.

This project was built as a comprehensive exercise in modern backend development, focusing on Clean Architecture, separation of concerns, and advanced data processing in Go.

It features a privacy-first recommendation engine that tracks user session history via HTTP-only cookies to generate personalized suggestions without client-side analytics scripts.

## Key learnings and Results:
- Architecting a scalable Go application using Clean Architecture principles to maintain a strict, logical separation between the transport layer, business logic, and data access.
- Developing modular, testable codebases by leveraging Go Interfaces and Dependency Injection for high component maintainability.
- Integrating external REST APIs via an abstracted data layer, allowing for seamless transitions between data sources (API/DB) without impacting core logic.
- Improving application performance through in-memory caching (with future Redis swap in mind) and UX personalized experiences via HTTP cookies.

## Key Features

### 1. Clean Architecture
The application maintains strict boundaries between layers to ensure testability and interchangeability of components:

**Delivery Layer** (`internal/controller/`): Handles HTTP transport, cookie parsing, and response formatting. It contains no business logic.

**Domain Layer** (`internal/usecase/`): Contains pure business rules (data retrieval, catalog filtering, recommendation logic). It has no knowledge of the database or HTTP.

**Data Layer** (`internal/repository/`): Manages data retrieval from the external API. It utilizes a reusable HTTP client (`pkg/httpclient`) to fetch resources and handles the mapping of raw DTOs into internal domain entities.

### 2. Concurrency & State Management

**Thread-Safe Caching:** The custom cache package (`pkg/cache`) uses `sync.RWMutex` to protect shared state during concurrent access.

**Janitor Pattern:** A background goroutine actively monitors and evicts expired cache items. This prevents memory leaks while keeping the main execution thread unblocked.

**Graceful Shutdown:** The server captures OS signals (`SIGTERM`, `SIGINT`) and utilizes `context.WithTimeout` to finish processing active requests before shutting down connections.

**Behavioral Recommendation Engine:** A custom algorithm that analyzes session history to determine brand affinity and category preferences. It employs a 4-slot strategy:
1.  *Resume Journey:* The user's most viewed vehicle.
2.  *Resume Journey 2:* The user's second most viewed vehicle.
3.  *Brand Loyalty:* A different model from the user's most visited manufacturer.
4.  *Competitor Comparison:* A model from the second most visited manufacturer.
5.  *Discovery:* A popular vehicle from the user's preferred category.

**Dynamic Comparison Grid:** A responsive, CSS Grid-based comparison tool that adapts layout columns based on the number of selected vehicles, implemented without heavy JavaScript frameworks.

**Server-Side Rendering:** High-performance HTML delivery using Go's `html/template` engine.

**Resilient Data Layer:** Handles external API failures gracefully with fallback strategies and strict data sanitization.

## Prerequisites

* **Go:** Version 1.22 or higher.
* **Node.js & NPM:** Required to run the external Cars API.

## Setup and Installation

The system consists of two parts: the external Data API (Node.js) and the Viewer Application (Go). Both must be running for the application to function correctly.

### 1. Start the Cars API (Backend Service)

> [!IMPORTANT]
> The data source is located in the separate `carapi` directory.
> Please refer to the `README.md` within that folder for detailed configuration options.

**Quick Start:**

1. Navigate to the API directory:

```bash
cd carapi
```

2. Install dependencies:
```bash
make build
# OR
npm install
```

3. Start the API server:
```bash
make run
```

> [!WARNING]
> By default, the API runs on `http://localhost:3000`.
> So, the Viewer config is set to request data from this host.
> If needed, you can change it in the config/local/local.json file.

### 2. Start the Viewer Application (Frontend)

Once the API is running, open a new terminal window to start the Go application.

1. Navigate to the project root.
2. Build the application
```bash
go build -o viewer cmd/viewer/main.go
```
2. Run the application:
```bash
go run cmd/viewer/main.go
```
3. Access the application in your browser:
```bash
http://localhost:8080
```

## Usage Guide

1.  **Browse Catalog:** Navigate to the homepage to see a list of available vehicles fetched from the Cars API.
2.  **View Details:** Click on any vehicle to view detailed specifications (Horsepower, Transmission, etc.).
3.  **Test Recommendations:**
    * Visit specific car models (e.g., view 3 different Audi models).
    * Return to the homepage or a different car page.
    * Observe the "Recommended for You" section, which will now prioritize Audi models and some other manufacturers based on your session history.

## Project Structure

The project follows the Standard Go Project Layout:

```text
viewer/
├── carapi/                     # External Node.js API (Data Source)
├── cmd/
│   └── viewer/
│       └── main.go             # Application entry point (Wires dependencies & starts App)
├── config/
│   └── local/
│       └── local.json          # Configuration for local environment
├── internal/
│   ├── app/                    # Composition Root (Initializes Core, Repository, Web)
│   ├── config/                 # Configuration structs and parsing logic
│   ├── controller/
│   │   └── httpserver/         # HTTP Transport Layer
│   │       ├── cookies/        # Secure Cookie logic (Session management)
│   │       ├── handlers/       # HTTP Handlers (Presentation logic)
│   │       ├── middleware/     # Request processing (Log, Recover, Context)
│   │       ├── router.go       # Route registration and routes
│   │       ├── server.go       # HTTP Server lifecycle (Graceful shutdown)
│   │       └── templates.go    # Custom Template Engine (Clone & Parse)
│   ├── domain/                 # Core Business Entities (Car, specs, manufacturers, filters)
│   ├── lib/
│   │   ├── adapter/            # Type-safe Adapters (e.g., Cache -> Domain)
│   │   └── e/                  # Error wrapping utilities
│   ├── repository/
│   │   └── webapi/             # Data Access Layer (Fetches from Node API)
│   └── usecase/
│       └── carstore/           # Business Logic (Catalog, filters, Recommendations)
├── pkg/                        # Reusable Library Code (No domain dependencies)
│   ├── cache/                  # Thread-safe Cache with Janitor
│   ├── httpclient/             # Resilient HTTP Client wrapper
│   └── logger/                 # Structured Logger setup
├── static/                     # Frontend Assets
│   ├── assets/                 # Images & Icons
│   ├── css/                    # Stylesheets
│   └── templates/              # HTML Templates (Layouts, Pages, Partials)
├── go.mod                      # Go Module definitions
├── TODO.md                     # Todo list with ideas and tasks
└── README.md                   # Project Documentation
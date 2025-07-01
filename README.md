# Game Integration API - Go Microservice Demo

## How to

### Quick Start

1. Copy `.env.example` to `.env` and update the configuration variables as needed.

1. Start the services using `docker-compose`:

```sh
docker-compose up --build
```

1. Seed the database (for testing):

Use the endpoint `POST /seed` to fill the auth table with this data:

```
{ID: 34633089486, Username: "player_34633089486", Password: "demo123!"}
{ID: 34679664254, Username: "player_34679664254", Password: "demo123!"}
{ID: 34616761765, Username: "player_34616761765", Password: "demo123!"}
{ID: 34673635133, Username: "player_34673635133", Password: "demo123!"}
```

1. Access Swagger documentation:

```sh
http://localhost:3000/swagger/index.html
```

### Notes

* Better to have [Docker](https://www.docker.com/) and [docker-compose](https://docs.docker.com/compose/) installed.
* Check the `Makefile` for additional useful commands

## About

* **`POST /auth`**: Authenticate players and provide a JWT.
* **`GET /player-info`**: Retrieve user details, including balance and currency.
* **`POST /withdraw`**: Process withdrawals (bet placements).
* **`POST /deposit`**: Handle deposits (bet settlements).
* **`POST /cancel`**: Roll back a previous transaction.

### Architecture

The system is designed using Clean Architecture principles to ensure low coupling and high cohesion:

* **Core Business Logic**: Encapsulated in `/service`.
* **Models**: Located in `/models`, defining the application's core entities.
* **Interfaces and Adapters**:

  * `/repository`: Handles database interactions.
  * `/transport`: Manages HTTP and other communication interfaces.

### Technologies Used

* **Web Framework**: [Echo](https://echo.labstack.com/)
* **ORM**: [Bun](https://bun.uptrace.dev/)
* **Database**: PostgreSQL

## Challenges and Design Decisions

### 1- Handling Unreliable Wallet Service

To interface with the mock wallet service:

* Added retry mechanisms to handle transient failures.
* Ensured idempotency to avoid duplicated transactions.

### 2- Choosing an ORM

* **Preferred Tool**: [sqlc](https://sqlc.dev/) for its raw SQL flexibility and schema-driven approach. + [goose](https://pressly.github.io/goose/) for managing migrations.
* **Current Choice**: Bun was chosen for this project due to its lightweight design, despite its documentation gaps.

### 3- Running Wallet Client on Linux (Fedora) amd64

1. **Install Dependencies** (once):

```sh
sudo dnf install qemu-user qemu-user-binfmt
```

2. **Set up emulators**:

```sh
docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
```

3. **Run Wallet Client**:

```sh
docker run --platform=linux/arm64 --rm -it -p 8000:8000 kentechsp/wallet-client
```

## Future Improvements

* **Improve Resiliency**: Explore an event-sourcing architecture (saga patterns) to better decouple services, keep an event log and mitigate the impact of service outages.
* **Enhanced Testing**: Add unit, integration, and end-to-end tests.
* **Monitoring**: Integrate tools like Prometheus and Grafana for observability.

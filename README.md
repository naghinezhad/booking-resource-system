# Booking Resource System

A backend service written in **Go** that allows users to reserve limited resources (such as rooms, cars, or servers) for a specific time range.  
The main challenge of this system is **preventing double booking under high concurrency** while handling a large number of simultaneous requests.

This project focuses on learning and demonstrating **concurrency handling in Go**, **distributed locking**, **caching**, and **high‑load request management**.

---

# Project Goals

The main goal of this project is to practice and demonstrate:

- Concurrency in Go (goroutines, channels, mutex)
- Race condition detection and prevention
- Handling high concurrent requests
- Preventing double booking
- Distributed locking using Redis
- Caching availability checks
- Stateless API design
- Load handling and request limiting
- Logging and monitoring basics
- Containerized deployment using Docker

---

# System Overview

Users can reserve a resource for a specific time range.  
When a reservation request arrives, the system checks whether the resource is available during the requested period.

If the resource is available:

- the reservation is created

If the resource is already booked:

- the request fails with a conflict response

The system must guarantee that **only one reservation succeeds for overlapping time ranges**, even if thousands of requests arrive simultaneously.

---

# API Endpoints

### Check Availability

GET /availability

Query Parameters:

- resource_id
- start_time
- end_time

Returns whether the resource is available in the requested time range.

---

### Create Reservation

POST /reserve

Request body example:

{
"resource_id": "room-1",
"start_time": "2026-04-10T10:00:00Z",
"end_time": "2026-04-10T11:00:00Z"
}

Responses:

- 201 Created → reservation successful
- 409 Conflict → resource already reserved
- 400 Bad Request → invalid input

---

### Get Reservations

GET /reservations

Optional query:

resource_id

Returns existing reservations for a resource.

---

# Database Design

The system uses **MongoDB**.

Collections:

resources

Fields:

- id
- name

reservations

Fields:

- id
- resource_id
- start_time
- end_time

Before inserting a new reservation, the system checks that there is **no overlapping reservation** for the same resource.

---

# Concurrency Handling

Concurrent requests can create race conditions where multiple users attempt to reserve the same resource simultaneously.

To prevent this, the system uses several techniques:

- goroutines for concurrent processing
- mutex protection where necessary
- distributed locking with Redis
- request limiting to prevent overload

Only one request should succeed when multiple requests attempt to reserve the same time slot.

---

# Redis Usage

Redis is used for two purposes:

### Caching

Availability checks can be cached to reduce database load.

Cache entries are invalidated after a successful reservation.

### Distributed Locking

A Redis lock (using SETNX) is used to ensure that only one process can attempt to create a reservation for a resource at a time.

This prevents double booking across multiple instances of the application.

---

# Stateless vs Stateful

The API itself is **stateless**, meaning each request contains all necessary information.

However:

- Redis cache is stateful
- Redis locking is stateful
- MongoDB stores persistent reservation data

This architecture allows the API servers to scale horizontally.

---

# Logging

The system logs important events for monitoring and debugging.

Three types of logs are used:

Success  
Logged when a reservation is created successfully.

Failure  
Logged when an internal error occurs.

Conflict  
Logged when a reservation fails because the resource is already reserved.

Logs are structured using a logging library such as zap.

---

# Load Handling

The system limits the number of concurrent requests using a configurable limit.

Environment variable:

MAX_CONCURRENT_REQUESTS

If the number of incoming requests exceeds this limit, additional requests are dropped.

This prevents the system from becoming overloaded.

---

# Configuration

Configuration is loaded from environment variables.

Example .env file:

SERVER_PORT=8080

MONGO_URI=mongodb://mongo:27017  
MONGO_DB=booking

REDIS_ADDR=redis:6379  
REDIS_PASS=

MAX_CONCURRENT_REQUESTS=1000

---

# Running the Project

The easiest way to run the project is with Docker Compose.

Start the system:

docker compose up --build

The following services will run:

- Go API server
- MongoDB
- Redis

The API will be available at:

http://localhost:8080

---

# Docker Services

The system runs three containers:

app  
The Go application that handles API requests.

mongo  
MongoDB database storing reservations.

redis  
Redis instance used for caching and distributed locking.

---

# Load Testing

A load testing script can simulate thousands of concurrent requests to verify system behavior.

Example scenario:

10,000 simultaneous requests try to reserve the same resource and time range.

Expected result:

- Only one request succeeds
- All other requests fail with conflict
- No double booking occurs

---

# Definition of Done

The project is considered complete when:

- Reservations do not overlap
- Double booking is prevented under concurrent requests
- The system successfully handles at least 10,000 concurrent reservation attempts
- Redis cache is used for availability checks
- Distributed locking is implemented
- Logging clearly records success, failure, and conflict events
- The project can be run easily with Docker

---

# Future Improvements

Possible improvements include:

- Adding Prometheus metrics
- Implementing request tracing with correlation IDs
- Adding authentication
- Improving cache strategies
- Implementing horizontal scaling with multiple API instances
- More advanced load testing

---

# Technologies Used

Go  
Gin (HTTP framework)  
MongoDB  
Redis  
Docker  
Docker Compose

---

This project is designed as a **learning exercise for building reliable concurrent backend systems in Go** while addressing real-world problems such as race conditions, distributed locking, and high-load request handling.

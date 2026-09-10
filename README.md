# Hotel Booking System

A production-oriented RESTful Hotel Booking System backend built with **Go**, **Gin**, **GORM**, and **PostgreSQL**.

The system provides authentication, role-based authorization, hotel and room management, hotel search, room availability, bookings, reviews, favorites, and a mock payment workflow with transaction handling.

---

## Features

### Authentication & Profile

* User registration
* User login
* BCrypt password hashing
* JWT authentication
* Authentication middleware
* Role-based authorization
* Get current authenticated user
* Update profile
* Change password

### Hotel Management

* Create hotel
* Get all hotels
* Get hotel by ID
* Update hotel
* Delete hotel
* Hotel details with rooms
* Hotel search and filtering

### Room Management

* Create room
* Get rooms by hotel
* Get room by ID
* Update room
* Delete room
* Room type management
* Room capacity management
* Room pricing
* Room status management

### Booking System

* Create booking
* Get all bookings
* Get user's bookings
* Get booking by ID
* Update booking
* Cancel booking
* Booking status management
* Server-side price calculation
* Guest capacity validation
* Booking date validation
* Double-booking prevention

### Room Availability

* Check room availability by date range
* Detect conflicts with existing bookings
* Exclude the current booking when updating a booking

### Hotel Search

Hotels can be filtered by:

* Hotel name
* City
* Country
* Number of stars
* Minimum price
* Maximum price
* Room type
* Minimum room capacity

### Reviews & Ratings

* Create hotel review
* Get hotel reviews
* Calculate average hotel rating
* Prevent duplicate reviews
* Validate ratings from 1 to 5

### Favorites

* Add hotel to favorites
* Get user's favorite hotels
* Remove hotel from favorites
* Prevent duplicate favorites

### Payments

* Create payment for a booking
* Get user's payments
* Get payment by ID
* Get all payments
* Update payment status
* Cash payment support
* Card payment support
* Server-side payment amount calculation
* Prevent duplicate payments
* UUID transaction ID generation

### Payment State Management

The payment system uses controlled state transitions:

```text
Booking: Pending
       |
       v
Payment: Pending
       |
       +----------------> Failed
       |                    |
       |                    v
       |             Booking remains Pending
       |
       v
Payment: Paid
       |
       v
Booking: Confirmed
       |
       v
Payment: Refunded
       |
       v
Booking: Cancelled
```

Payment and booking updates are handled inside database transactions to maintain data consistency.

---

## Tech Stack

| Technology        | Purpose                      |
| ----------------- | ---------------------------- |
| Go                | Backend programming language |
| Gin               | HTTP web framework           |
| GORM              | ORM                          |
| PostgreSQL        | Relational database          |
| JWT               | Authentication               |
| BCrypt            | Password hashing             |
| Swagger / OpenAPI | API documentation            |
| UUID              | Transaction ID generation    |
| Postman           | API testing                  |
| Git & GitHub      | Version control              |

---

## Architecture

The project follows a layered architecture with clear separation of responsibilities.

```text
                         HTTP Request
                              |
                              v
                       +--------------+
                       |    Handler   |
                       +--------------+
                              |
                              v
                       +--------------+
                       |    Service   |
                       +--------------+
                              |
                              v
                       +--------------+
                       |  Repository  |
                       +--------------+
                              |
                              v
                       +--------------+
                       |  PostgreSQL  |
                       +--------------+
```

### Handler Layer

Responsible for:

* HTTP requests
* Request binding
* Request parameter parsing
* HTTP status codes
* API responses

### Service Layer

Responsible for:

* Business logic
* Validation
* Authorization rules
* Booking rules
* Availability checks
* Price calculations
* Payment state transitions
* Transaction handling

### Repository Layer

Responsible for:

* Database queries
* CRUD operations
* GORM interactions
* Database transaction repositories

### Model Layer

Contains:

* Database models
* Request models
* Data structures

---

## Project Structure

```text
hotel-booking-system/
│
├── cmd/
│   ├── api/
│   │   └── main.go
│   │
│   └── seed/
│       └── main.go
│
├── configs/
│   └── config.go
│
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── internal/
│   │
│   ├── auth/
│   │   ├── handler/
│   │   ├── model/
│   │   ├── repository/
│   │   └── service/
│   │
│   ├── user/
│   │   ├── handler/
│   │   ├── model/
│   │   ├── repository/
│   │   └── service/
│   │
│   ├── profile/
│   │   ├── handler/
│   │   ├── model/
│   │   ├── repository/
│   │   └── service/
│   │
│   ├── hotel/
│   │   ├── handler/
│   │   ├── model/
│   │   ├── repository/
│   │   └── service/
│   │
│   ├── room/
│   │   ├── handler/
│   │   ├── model/
│   │   ├── repository/
│   │   └── service/
│   │
│   ├── booking/
│   │   ├── handler/
│   │   ├── model/
│   │   ├── repository/
│   │   └── service/
│   │
│   ├── review/
│   │   ├── handler/
│   │   ├── model/
│   │   ├── repository/
│   │   └── service/
│   │
│   ├── favorite/
│   │   ├── handler/
│   │   ├── model/
│   │   ├── repository/
│   │   └── service/
│   │
│   └── payment/
│       ├── handler/
│       ├── model/
│       ├── repository/
│       └── service/
│
├── migrations/
│
├── pkg/
│   ├── database/
│   ├── hash/
│   ├── jwt/
│   ├── middleware/
│   ├── response/
│   └── validator/
│
├── routes/
│   └── routes.go
│
├── .env
├── .gitignore
├── go.mod
└── go.sum
```

---

## Authentication

The API uses **JWT Bearer Authentication**.

After a successful login, the API returns a JWT token.

Protected endpoints require:

```text
Authorization: Bearer <JWT_TOKEN>
```

The authentication middleware extracts the authenticated user's:

```text
userID
role
```

from the JWT claims.

---

## Roles & Authorization

The system supports three main roles:

| Role     | Access                                                      |
| -------- | ----------------------------------------------------------- |
| Customer | Bookings, reviews, favorites, payments, profile             |
| Manager  | Hotel and room management, bookings and payments management |
| Admin    | Full administrative access                                  |

Role-based authorization is implemented using middleware.

---

## Booking Flow

```text
Customer
    |
    v
Create Booking
    |
    v
Pending
    |
    v
Create Payment
    |
    +---------> Failed
    |              |
    |              v
    |           Pending
    |
    v
Paid
    |
    v
Confirmed
```

The booking system validates:

* Room existence
* Room availability
* Room capacity
* Check-in date
* Check-out date
* Booking duration
* Server-side total price

---

## Payment Flow

Payment amounts are calculated from the booking on the server.

The client cannot control the final payment amount.

### Payment Methods

```text
cash
card
```

### Payment Statuses

```text
pending
paid
failed
refunded
```

### Transaction ID

Successful payments receive a unique UUID transaction ID.

### Database Transactions

Payment and booking status updates are executed inside a database transaction.

This ensures that related updates succeed or fail together.

---

## API Documentation

The project includes Swagger / OpenAPI documentation.

When the API is running, Swagger UI is available at:

```text
http://localhost:8081/swagger/index.html
```

Swagger documents the main API modules:

* Authentication
* Profile
* Hotels
* Hotel Search
* Hotel Details
* Rooms
* Room Availability
* Bookings
* Reviews
* Favorites
* Payments

---

## API Endpoints

### Authentication

```text
POST   /auth/register
POST   /auth/login
GET    /auth/me
PUT    /auth/profile
PUT    /auth/change-password
```

### Hotels

```text
GET    /hotels
GET    /hotels/search
GET    /hotels/:id
GET    /hotels/:id/details
POST   /hotels
PUT    /hotels/:id
DELETE /hotels/:id
```

### Rooms

```text
POST   /hotels/:hotel_id/rooms
GET    /hotels/:hotel_id/rooms
GET    /rooms/:id
PUT    /rooms/:id
DELETE /rooms/:id
GET    /rooms/:id/availability
```

### Bookings

```text
POST   /bookings
GET    /bookings
GET    /bookings/my
GET    /bookings/:id
PUT    /bookings/:id
DELETE /bookings/:id
PATCH  /bookings/:id/status
```

### Reviews

```text
POST   /hotels/:id/reviews
GET    /hotels/:id/reviews
GET    /hotels/:id/rating
```

### Favorites

```text
POST   /hotels/:id/favorite
GET    /favorites
DELETE /hotels/:id/favorite
```

### Payments

```text
POST   /bookings/:id/payment
GET    /payments/my
GET    /payments/:id
GET    /payments
PATCH  /payments/:id/status
```

---

## Environment Variables

Create a `.env` file in the project root:

```env
APP_NAME=Hotel Booking System
PORT=8081

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=hotel_booking_db
DB_SSLMODE=disable

JWT_SECRET=your_secret_key
JWT_EXPIRE=24h
```

> Never commit real database credentials or JWT secrets to GitHub.

---

## Getting Started

### Prerequisites

Make sure you have:

* Go installed
* PostgreSQL installed and running
* Git installed

### 1. Clone the repository

```bash
git clone https://github.com/Mosteben/hotel-booking-system.git
cd hotel-booking-system
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Create the database

Create a PostgreSQL database named:

```text
hotel_booking_db
```

### 4. Configure environment variables

Create `.env` and configure your PostgreSQL and JWT settings.

### 5. Run the API

```bash
go run ./cmd/api
```

The API will run on:

```text
http://localhost:8081
```

### 6. Open Swagger

```text
http://localhost:8081/swagger/index.html
```

---

## Testing

API endpoints have been tested using **Postman**.

The project also contains unit tests for the payment service.

Run the test suite with:

```bash
go test ./...
```

---

## API Response Format

The API uses a consistent response structure.

### Success Response

```json
{
  "success": true,
  "message": "operation completed successfully",
  "data": {}
}
```

### Error Response

```json
{
  "success": false,
  "message": "error message"
}
```

---

## Database

The project uses **PostgreSQL** with **GORM**.

```text
Application
     |
     v
GORM
     |
     v
PostgreSQL
```

Database models are managed through the application's database migration/AutoMigrate setup.

---

## Security

The backend implements:

* JWT authentication
* BCrypt password hashing
* Role-based authorization
* Protected routes
* User ownership validation
* Server-side price calculation
* Input validation
* Booking conflict prevention
* Payment state validation
* Database transactions

---

## Future Improvements

Planned improvements include:

* React frontend integration
* Real payment gateway integration
* Stripe or Paymob integration
* Email notifications
* Booking confirmation emails
* Password reset
* Pagination
* Advanced hotel filtering
* Admin dashboard
* More automated tests
* Docker support
* CI/CD pipeline
* Production deployment
* Monitoring and logging

---

## Author

**Mostafa Alaaeldin**

Backend-focused Software Engineer specializing in **Go, Gin, REST APIs, and backend system architecture**.

---

## License

This project was developed as a backend engineering and portfolio project.

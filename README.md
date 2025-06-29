# Chirpy 🐦

A modern, RESTful API for a Twitter-like social media platform built with Go. Chirpy allows users to create, read, and manage short messages called "chirps" with full authentication and authorization.

## Features

- **User Management**: Register, login, and manage user accounts
- **Chirp Operations**: Create, read, update, and delete chirps
- **Authentication**: JWT-based authentication with refresh tokens
- **Authorization**: User-specific chirp management
- **Webhooks**: Integration with external services (Polka)
- **Metrics**: Built-in monitoring and analytics
- **Database**: PostgreSQL with SQLC for type-safe database operations

## Tech Stack

- **Backend**: Go 1.23.2
- **Database**: PostgreSQL
- **ORM**: SQLC (SQL Compiler)
- **Authentication**: JWT tokens
- **HTTP Server**: Standard library `net/http`
- **Environment**: GoDotEnv for configuration

## Project Structure

```
chirpy/
├── assets/                 # Static assets
├── internal/
│   ├── auth/              # Authentication logic
│   └── database/          # Database models and queries
├── sql/
│   ├── queries/           # SQLC query definitions
│   └── schema/            # Database migrations
├── handler_*.go           # HTTP request handlers
├── main.go               # Application entry point
├── go.mod                # Go module definition
└── sqlc.yaml             # SQLC configuration
```

## API Endpoints

### Health & Metrics
- `GET /api/healthz` - Health check
- `GET /admin/metrics` - Application metrics
- `POST /admin/reset` - Reset metrics

### User Management
- `POST /api/users` - Create new user
- `PUT /api/users` - Update user profile
- `POST /api/login` - User authentication
- `POST /api/refresh` - Refresh JWT token
- `POST /api/revoke` - Revoke refresh token

### Chirp Operations
- `POST /api/chirps` - Create new chirp
- `GET /api/chirps` - Get all chirps (with optional filtering)
- `GET /api/chirps/{chirp_id}` - Get specific chirp
- `DELETE /api/chirps/{chirp_id}` - Delete chirp

### Webhooks
- `POST /api/polka/webhooks` - Handle external webhooks

### Static Files
- `/app/*` - Serve static files

## Database Schema

### Users Table
- `id` (UUID, Primary Key)
- `created_at` (Timestamp)
- `updated_at` (Timestamp)
- `email` (VARCHAR, Unique)
- `hashed_password` (VARCHAR)
- `is_chirpy_red` (Boolean)

### Chirps Table
- `id` (UUID, Primary Key)
- `created_at` (Timestamp)
- `updated_at` (Timestamp)
- `body` (TEXT)
- `user_id` (UUID, Foreign Key to users)

### Refresh Tokens Table
- `token` (VARCHAR, Primary Key)
- `created_at` (Timestamp)
- `updated_at` (Timestamp)
- `user_id` (UUID, Foreign Key to users)
- `expires_at` (Timestamp)
- `revoked_at` (Timestamp, Nullable)

## Setup & Installation

### Prerequisites
- Go 1.23.2 or higher
- PostgreSQL database
- SQLC (for code generation)

### Environment Variables
Create a `.env` file with the following variables:

```env
DB_URL=postgres://username:password@localhost:5432/chirpy_db
PLATFORM=your_platform_name
JWT_SECRET=your_jwt_secret_key
POLKA_KEY=your_polka_webhook_key
```

### Database Setup
1. Create a PostgreSQL database
2. Run database migrations:
   ```bash
   # Install goose (migration tool)
   go install github.com/pressly/goose/v3/cmd/goose@latest
   
   # Run migrations
   goose -dir sql/schema postgres "your_db_url" up
   ```

### Code Generation
Generate database code with SQLC:
```bash
sqlc generate
```

### Running the Application
```bash
go run .
```

The server will start on port 8080.

## API Usage Examples

### Create a User
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'
```

### Login
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'
```

### Create a Chirp
```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"body": "Hello, Chirpy!"}'
```

### Get All Chirps
```bash
# Get all chirps (ascending order)
curl http://localhost:8080/api/chirps

# Get all chirps (descending order)
curl http://localhost:8080/api/chirps?sort=desc

# Get chirps by specific user
curl http://localhost:8080/api/chirps?author_id=USER_UUID
```

## Development

### Adding New Queries
1. Add SQL queries to `sql/queries/*.sql`
2. Run `sqlc generate` to generate Go code
3. Use the generated methods in your handlers

### Database Migrations
1. Create new migration files in `sql/schema/`
2. Use goose format with `-- +goose Up` and `-- +goose Down`
3. Run migrations with `goose up`

### Testing
```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.
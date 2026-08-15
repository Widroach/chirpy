# Chirpy

A guided project that demonstrates building a production-style HTTP server in Go with a REST API, JSON responses, authentication (JWT), authorization (refresh tokens, API keys, ownership checks), and a PostgreSQL database managed with sqlc and goose.

## Features

- RESTful JSON API for users and chirps (posts)
- JWT-based access tokens and rotating refresh tokens
- Middleware-based authorization (bearer tokens, refresh tokens, Polka API key)
- Password hashing with Argon2id
- Profanity filtering on chirp bodies
- Request metrics tracking and an admin reset endpoint
- SQL queries generated from `.sql` files with sqlc

## Requirements

- Go 1.26+
- PostgreSQL
- sqlc (to regenerate database code after schema/query changes)
- goose (to run database migrations)

## Install

```sh
git clone https://github.com/Widroach/chirpy.git
cd chirpy
go mod download
```

## Setup

Create a `.env` file in the project root:

```
DB_URL=postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable
PLATFORM=dev
SECRET=<a long random string used to sign JWTs>
POLKA_KEY=<a long random string used for webhooks>
```

| Variable    | Description                                                        |
| ----------- | ------------------------------------------------------------------ |
| `DB_URL`    | PostgreSQL connection string                                        |
| `PLATFORM`  | `dev` enables the `/admin/reset` endpoint; any other value disables it |
| `SECRET`    | Secret used to sign and verify JWTs                                 |
| `POLKA_KEY` | API key expected in the `Authorization: ApiKey` header on webhook calls |

Generate the database code and run migrations:

```sh
sqlc generate
goose postgres "$DB_URL" up
```

Start the server:

```sh
go run .
```

The server listens on `:8080`.

## API Endpoints

### Health

| Method | Path            | Description          |
| ------ | --------------- | -------------------- |
| GET    | `/api/healthz`  | Returns `{"status":"OK"}` |

### Users

| Method | Path          | Auth  | Description                                   |
| ------ | ------------- | ----- | --------------------------------------------- |
| POST   | `/api/users`  | None  | Create a user. Body: `{"email","password"}`. Returns the user (201). |
| POST   | `/api/login`  | None  | Login. Body: `{"email","password"}`. Returns `id`, `created_at`, `updated_at`, `email`, `is_chirpy_red`, `token` (JWT), and `refresh_token` (200). |
| PUT    | `/api/users`  | JWT   | Update the authenticated user. Body: `{"email","password"}`. Returns the updated user. |

### Chirps

| Method | Path                     | Auth | Description |
| ------ | ------------------------ | ---- | ----------- |
| GET    | `/api/chirps`            | None | List chirps, newest first. Optional query params: `author_id` (filter by author), `sort=asc` (oldest first). |
| GET    | `/api/chirps/{chirpId}`  | None | Get a single chirp. |
| POST   | `/api/chirps`            | JWT  | Create a chirp. Body: `{"body"}` (max 140 chars). Profanity is censored. Returns the chirp (201). |
| DELETE | `/api/chirps/{chirpId}`  | JWT  | Delete a chirp you own (403 if it belongs to someone else). Returns 204. |

### Refresh Tokens

| Method | Path           | Auth          | Description |
| ------ | -------------- | ------------- | ----------- |
| POST   | `/api/refresh` | Refresh token | Issues a new JWT (and rotates the refresh token). |
| POST   | `/api/revoke`  | Refresh token | Revokes a refresh token. Returns 204. |

### Polka Webhook (authorization)

| Method | Path                   | Auth      | Description |
| ------ | ---------------------- | --------- | ----------- |
| POST   | `/api/polka/webhooks`  | Polka key | Body: `{"event","data":{"user_id"}}`. On `user.upgraded`, sets the user's `is_chirpy_red` to true. Returns 204. |

### Admin

| Method | Path              | Description |
| ------ | ----------------- | ----------- |
| GET    | `/admin/metrics`  | HTML page showing total request count. |
| POST   | `/admin/reset`    | Deletes all users and chirps and resets metrics. Only enabled when `PLATFORM=dev`. |
| GET    | `/app/`           | Serves static files from the project root. |

## Authentication

All authenticated endpoints expect an `Authorization` header:

- JWT endpoints: `Authorization: Bearer <access_token>`
- Refresh token endpoints: `Authorization: Bearer <refresh_token>`
- Webhook endpoint: `Authorization: ApiKey <polka_key>`

Authorization is enforced via middleware: JWTs are validated and the user id is injected into the request context, refresh tokens are checked against the database, and the Polka key must match `POLKA_KEY`. Chirp deletion also checks that the chirp belongs to the authenticated user (403 otherwise).

## Project Structure

```
cmd/            (or main.go)
http/
  controller/   request handlers
  middleware/   auth + metrics middleware
internal/
  auth/         password hashing, JWT, tokens
  database/     sqlc-generated Go code
sql/
  schema/       goose migrations
  queries/      sqlc queries
routes.go       route registration
```

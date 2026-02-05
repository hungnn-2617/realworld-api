# RealWorld API

A RESTful API implementation of [RealWorld](https://github.com/gothinkster/realworld) spec using Go, Gin, GORM and PostgreSQL.

## Quick Start

```bash
# Clone and setup
git clone <repository-url>
cd realworld-api
cp .env.example .env

# Run with Docker
docker-compose up -d

# Or run locally
go run cmd/app/main.go
```

## API Endpoints

| Method | Endpoint                           | Auth | Description      |
| ------ | ---------------------------------- | :--: | ---------------- |
| POST   | `/api/users`                       |  -   | Register         |
| POST   | `/api/users/login`                 |  -   | Login            |
| GET    | `/api/user`                        |  ✓   | Get current user |
| PUT    | `/api/user`                        |  ✓   | Update user      |
| GET    | `/api/profiles/:username`          |  ○   | Get profile      |
| POST   | `/api/profiles/:username/follow`   |  ✓   | Follow user      |
| DELETE | `/api/profiles/:username/follow`   |  ✓   | Unfollow user    |
| GET    | `/api/articles`                    |  ○   | List articles    |
| GET    | `/api/articles/feed`               |  ✓   | User feed        |
| GET    | `/api/articles/:slug`              |  -   | Get article      |
| POST   | `/api/articles`                    |  ✓   | Create article   |
| PUT    | `/api/articles/:slug`              |  ✓   | Update article   |
| DELETE | `/api/articles/:slug`              |  ✓   | Delete article   |
| POST   | `/api/articles/:slug/favorite`     |  ✓   | Favorite         |
| DELETE | `/api/articles/:slug/favorite`     |  ✓   | Unfavorite       |
| GET    | `/api/articles/:slug/comments`     |  ○   | Get comments     |
| POST   | `/api/articles/:slug/comments`     |  ✓   | Add comment      |
| DELETE | `/api/articles/:slug/comments/:id` |  ✓   | Delete comment   |
| GET    | `/api/tags`                        |  -   | Get tags         |

**Auth:** ✓ Required | ○ Optional | - None

## Testing

```bash
go test ./...
```

## License

MIT

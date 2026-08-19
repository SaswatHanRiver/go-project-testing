# Go Product API

A REST API built with Go as part of Divii's one-week Go learning initiative.
Built by a Java Spring Boot developer learning Go from scratch.

## What does it do?

A CRUD REST API for managing products, backed by PostgreSQL.

**Tech stack:**
| Layer | Go | Spring Boot equivalent |
|---|---|---|
| Web framework | Gin | Spring MVC |
| ORM | GORM | JPA / Hibernate |
| Database | PostgreSQL | PostgreSQL |
| Docs | Swaggo | SpringDoc / Swagger UI |
| Logging | log/slog | SLF4J + Logback |

**Endpoints:**
```
GET    /api/products        → list all products
GET    /api/products/:id    → get one product
POST   /api/products        → create product
PUT    /api/products/:id    → update product
DELETE /api/products/:id    → delete product (soft delete)
GET    /swagger/index.html  → Swagger UI
```

**Go-specific features implemented:**
- Goroutine + channel background job worker — fires on every create/update/delete
- `context.WithTimeout` on every request (5s deadline)
- `context.WithCancel` for graceful shutdown on SIGINT/SIGTERM
- `log/slog` structured JSON logging
- Unit tests with mock repository (no DB required)
- Multi-stage Dockerfile

## How to run

**Prerequisites:** Go 1.21+, PostgreSQL running locally

```bash
# 1. Create the database
psql -U postgres -c "CREATE DATABASE go_testing;"

# 2. Install dependencies
go mod tidy

# 3. Generate Swagger docs
go install github.com/swaggo/swag/cmd/swag@latest
swag init

# 4. Configure .env
cp .env.example .env   # edit DB credentials if needed

# 5. Run
go run main.go

# 6. Run tests
go test ./...
```

**Using Docker:**
```bash
docker build -t go-product-api .
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=postgres \
  -e DB_NAME=go_testing \
  go-product-api
```

Swagger UI: http://localhost:8080/swagger/index.html

## Where did Claude get it right, and where did I have to fix its output?

**What Claude got right:**
- The layered project structure (controller → service → repository) mapped naturally from Spring Boot. The Gin routing, GORM AutoMigrate, and Swagger annotations all worked without changes.
- The goroutine + channel worker pattern was correct and clean on the first attempt.
- `log/slog` integration was straightforward.

**What I had to understand and verify myself:**
- The repository interface for mocking — Claude generated it but I had to understand *why* Go needs an explicit interface for mocking (unlike Mockito which can mock concrete classes). In Go, mocking only works through interfaces.
- `defer cancel()` — Claude added this in every handler but I had to look up what it actually does. It releases the context's resources when the function returns. Skipping it causes a goroutine leak.
- The difference between a value receiver `(r ProductRepository)` and a pointer receiver `(r *ProductRepository)`. Claude used `*ProductRepository` for the interface implementation — I had to verify this was correct because Go won't satisfy an interface if the receiver type doesn't match.

## Which Go concept confused me most?

**Pointers and value vs pointer receivers.**

Coming from Java, everything is a reference by default. In Go, structs are copied by value unless you use a pointer. This is where I got burned:

```go
// WRONG - this creates a copy, changes are lost
func (r ProductRepository) Create(product models.Product) error { ... }

// CORRECT - pointer receiver, modifies the actual struct
func (r *ProductRepository) Create(product *models.Product) error { ... }
```

The rule I now understand:
- Use pointer receivers (`*T`) when the method needs to modify the struct, or when the struct is large
- Use value receivers (`T`) for small read-only structs
- **All methods on a type should be consistent** — if one needs a pointer receiver, make them all pointer receivers

The second thing that confused me was `context.Context`. In Spring Boot, timeouts and cancellation are handled by annotations or config. In Go, you pass a `context` explicitly through every function call. It felt verbose at first, but I now understand why — it makes cancellation explicit and traceable, rather than hidden in a framework.

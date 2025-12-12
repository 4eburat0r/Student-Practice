# user-service (MVP)

## Run locally (dev)
cd backend/services/user-service
docker-compose -f build/docker-compose.yml down -v
docker-compose -f build/docker-compose.yml up --build

- Postgres will initialize and run SQL in migrations/ on first volume creation.
- Service available on http://localhost:8081

## Endpoints
- POST /users { email, password, role } -> create user
- GET /users/:id -> get user
Protected (requires Authorization header):
- GET /users/students/me
- PUT /users/students/me
- DELETE /users/students/me
- GET /users/employers/me
- PUT /users/employers/me
- DELETE /users/employers/me

For local dev you can use header: Authorization: Bearer devtoken (middleware dev shortcut).

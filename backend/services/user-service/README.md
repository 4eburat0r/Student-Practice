# user-service (MVP)

Run:
cd backend/services/user-service
docker-compose -f build/docker-compose.yml up --build

Endpoints:
- POST /users { email, password, role } -> create user
- GET /users/:id -> get user (no password)
- Protected (Authorization: Bearer <token>):
  - GET /users/students/me
  - PUT /users/students/me
  - DELETE /users/students/me
  - GET /users/employers/me
  - PUT /users/employers/me
  - DELETE /users/employers/me

Notes:
- Auth introspection endpoint: AUTH_INTROSPECT_URL expects POST {"token":"..."} and returns {user_id, role, active}
- DB auto-init uses SQL in migrations folder (first run only)

Port: 8081
# Overview

Plan de implementación para Reputation Ranking. La vista `leaderboard` se define en SQL y se asegura en arranque vía `EnsureLeaderboardView`. Endpoints: GET Top 3; GET leaderboard (offset, limit, student_id opcional) con data + current_user. Stack: Go, Gin, PostgreSQL.

# Tasks

- [x] 1. Modelo de datos y DTOs
  - [x] 1.1 DTOs en `internal/models/ranking_dto.go`
  - [x] 1.2 `LeaderboardRow` en `internal/repositories/ranking_repository.go`
  - [x] 1.3 Interfaces `RankingRepository` + servicio en `internal/services/ranking/ranking_service.go`

- [x] 2. Checkpoint – Tipos compilando

- [x] 3. RankingRepository
  - [x] 3.1 Tests de servicio con `fakeRankingRepo` en `ranking_service_test.go`
  - [x] 3.2 Implementación GORM `ranking_repository_gorm.go` (vista `leaderboard`)

- [x] 4. Checkpoint – RankingRepository implementado

- [x] 5. RankingService
  - [x] 5.1 Tests en `ranking_service_test.go`
  - [x] 5.2 `Top3` / `Leaderboard` en `ranking_service.go`

- [x] 6. Checkpoint – Servicio implementado

- [x] 7. Manejo de errores y validación
  - [x] 7.1 Normalización offset/limit y UUID `student_id` en `ranking_handler.go`
  - [x] 7.2 Tests en `ranking_handler_test.go`

- [x] 8. RankingHandler – Top 3
  - [x] 8.1 Tests handler
  - [x] 8.2 `Top3` handler

- [x] 9. RankingHandler – Leaderboard
  - [x] 9.1 Tests handler
  - [x] 9.2 `Leaderboard` handler

- [x] 10. Integración
  - [x] 10.1 Rutas en `main.go`: `GET /api/ranking/top3`, `GET /api/ranking/leaderboard` (grupo **`/api` público**, sin `staffAPI`)
- [x] 10.1d Lookup público: query `email` o `student_id_code` (o `student_id` UUID); exclusión mutua; resolución vía `StudentRepository` + `current_user`
  - [x] 10.2 `go test ./...`

- [x] 11. Checkpoint final

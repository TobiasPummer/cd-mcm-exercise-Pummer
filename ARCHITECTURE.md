# Architecture Documentation – Product Catalog API

## Request Flow

An incoming HTTP request durchläuft folgende Schichten:

HTTP Request
│
▼
┌─────────────────┐
│  Gorilla Mux    │  ← Router: matched URL + HTTP Method
│     Router      │     z.B. GET /products/1 → GetProduct Handler
└────────┬────────┘
│
▼
┌─────────────────┐
│    Handler      │  ← Liest Request, validiert Input,
│ (handler.go /   │     schreibt HTTP Response (JSON)
│ postgres_han..) │
└────────┬────────┘
│
▼
┌─────────────────┐
│  Store Interface│  ← Abstraktion: GetAll, GetByID,
│                 │     Create, Update, Delete
└────────┬────────┘
│
┌────┴────┐
▼         ▼
┌────────┐ ┌──────────┐
│Memory  │ │Postgres  │
│Store   │ │Store     │
│(RAM)   │ │(DB)      │
└────────┘ └──────────┘

### Ablauf am Beispiel: POST /products

1. Client sendet HTTP POST an `/products` mit JSON Body
2. Gorilla Mux Router erkennt die Route und ruft `CreateProduct` Handler auf
3. Handler dekodiert den JSON Body in ein `Product` Struct
4. Handler ruft `p.Validate()` auf – bei leerem Name → 400 Bad Request
5. Handler ruft `store.Create(p)` auf
6. Store speichert das Produkt (Memory oder PostgreSQL) und gibt es mit ID zurück
7. Handler antwortet mit 201 Created + JSON des erstellten Produkts

---

## Project Structure
cmd/api/main.go              → Entry point: startet Server, wählt Store
internal/model/product.go    → Product Struct + Validate()
internal/store/memory.go     → MemoryStore Implementierung
internal/store/postgres.go   → PostgresStore Implementierung
internal/handler/handler.go  → HTTP Handler (MemoryStore)
internal/handler/postgres_handler.go → HTTP Handler (PostgresStore)
internal/handler/handler_test.go     → Unit Tests
Dockerfile                   → Multi-stage Build
docker-compose.yml           → Lokale Entwicklung mit PostgreSQL
.github/workflows/ci.yml     → GitHub Actions CI Pipeline

## MemoryStore vs PostgresStore

| Kriterium        | MemoryStore                        | PostgresStore                        |
|------------------|------------------------------------|--------------------------------------|
| Persistenz       | Nein – Daten weg nach Neustart     | Ja – Daten bleiben dauerhaft         |
| Geschwindigkeit  | Sehr schnell (kein Netzwerk)       | Etwas langsamer (DB-Verbindung)      |
| Skalierung       | Nur eine Instanz möglich           | Mehrere API-Instanzen möglich        |
| Setup            | Kein Setup nötig                   | PostgreSQL muss laufen               |
| Verwendung       | Unit Tests, lokale Entwicklung     | Produktion, Docker Compose           |

### Wann MemoryStore?
- **Unit Tests**: Kein Datenbankserver nötig, Tests laufen schnell und isoliert
- **Lokale Entwicklung**: Schnell starten ohne Docker

### Wann PostgresStore?
- **Produktion**: Daten müssen persistent gespeichert werden
- **Mehrere API-Instanzen**: Alle teilen sich dieselbe Datenbank
- **Docker Compose**: Wie in diesem Projekt konfiguriert

### Trade-offs
Der MemoryStore verliert alle Daten beim Neustart – ideal für Tests, aber unbrauchbar
für Produktion. Der PostgresStore braucht eine laufende Datenbankverbindung, bietet
dafür aber echte Persistenz und erlaubt horizontale Skalierung der API.
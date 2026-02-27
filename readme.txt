# Rick & Morty API

Author: Wiktor Szydłowski

Go backend service that uses the [Rick and Morty API](https://rickandmortyapi.com/) and exposes two new endpoints for combined search and character pair analysis.

## Running the Service

**With Docker (recommended):**
```
docker compose up
```

**Without Docker:**
```
go run main.go
```

The server starts on `http://localhost:8080` and pre-warms its cache on startup by fetching all characters and episodes from the API.

## Endpoints

### GET /search

Performs a combined search across characters, locations, and episodes by name. All three resource types are queried concurrently.

**Query parameters:**
- `term` (required) — the search term to match against names
- `limit` (optional) — maximum number of results to return - if omitted, all matches are returned

**Example:**
```
curl "http://localhost:8080/search?term=rick&limit=2"
```

**Response:**
```json
[
  {"name":"Rick Sanchez","type":"character","url":"https://rickandmortyapi.com/api/character/1"},
  {"name":"Adjudicator Rick","type":"character","url":"https://rickandmortyapi.com/api/character/8"}
]
```

Results are ordered as: characters first, then locations, then episodes. The limit is applied after combining all results.

---

### GET /top-pairs

Returns character pairs that have appeared together in the most episodes, sorted by shared episode count descending. Character and episode data is served from an in-memory cache (TTL: 1 hour).

**Query parameters:**
- `min` (optional) — minimum number of shared episodes for a pair to be included
- `max` (optional) — maximum number of shared episodes for a pair to be included
- `limit` (optional) — maximum number of pairs to return - defaults to 20

**Example:**
```
curl "http://localhost:8080/top-pairs?limit=2"
```

**Response:**
```json
[
  {"character1":{"name":"Rick Sanchez","url":"https://rickandmortyapi.com/api/character/1"},"character2":{"name":"Morty Smith","url":"https://rickandmortyapi.com/api/character/2"},"episodes":51},
  {"character1":{"name":"Rick Sanchez","url":"https://rickandmortyapi.com/api/character/1"},"character2":{"name":"Summer Smith","url":"https://rickandmortyapi.com/api/character/3"},"episodes":42}
]
```

## Caching

Characters and episodes are cached in memory for 1 hour. The cache is populated on startup, so the first response from `/top-pairs` is fast. The `/search` endpoint does not use the cache — it queries the upstream API directly on each request.

## Running Tests

```
go test ./...
```

Tests cover parameter parsing, pair counting logic, and live integration tests against the upstream API.

## Project Structure

```
main.go                  — server entry point, cache warm-up, route registration
handlers/search.go       — /search endpoint handler
handlers/toppairs.go     — /top-pairs endpoint handler
client/rickmorty.go      — upstream API fetch functions
client/cache.go          — in-memory cache for characters and episodes
models/types.go          — shared data types
```

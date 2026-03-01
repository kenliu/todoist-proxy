# todoist-proxy

A Go HTTP proxy for the [Todoist API v1](https://developer.todoist.com/api/v1/) that filters sync responses to a specific set of projects (and their descendants). Useful for exposing a subset of your Todoist data to a client without giving it access to your entire account.

## How it works

- **`POST /api/v1/sync`** — intercepted. The request is forwarded to Todoist, and the response is filtered to only include the allowed projects and all resources (tasks, sections, labels, etc.) that belong to them.
- **All other endpoints** — passed through transparently to `api.todoist.com`.

Clients authenticate directly with their Todoist API token as usual. The proxy forwards credentials without storing or inspecting them.

## Configuration

| Env var | Required | Description |
|---|---|---|
| `TODOIST_PROXY_ALLOW` | Yes | Comma-separated Todoist project IDs to expose (e.g. `123456,789012`). Child projects are included automatically. |
| `PORT` | No | Port to listen on (default: `8080`). |

## Running

```bash
go build -o todoist-proxy .
TODOIST_PROXY_ALLOW=123456,789012 PORT=8080 ./todoist-proxy
```

Or without building:

```bash
TODOIST_PROXY_ALLOW=123456,789012 go run .
```

Then point your Todoist client at `http://localhost:8080` instead of `https://api.todoist.com`.

## Development

```bash
# Run tests
go test ./...

# Build
go build ./...
```

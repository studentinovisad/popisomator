# Popisomator

## Development

Reqs: podman/docker compose, pnpm, golang.

1. Create the local configuration:

```sh
cp .env.example .env
```

2. In one terminal, start PostgreSQL, migrations, and the backend:

```sh
podman compose up --build backend
```

The backend is available at `http://localhost:<BACKEND_PORT>` (`8080` by default). On later starts, you can omit `--build` unless an image changed.

Alternatively, to iterate on the backend without rebuilding the container on every change, run it locally with live reload (uses the `.env` from step 1):

```sh
podman compose up postgres
cd backend
make migrate-apply # once, or after pulling new migrations
make dev
```

This uses [air](https://github.com/air-verse/air) (tracked as a `go tool` dependency, invoked as `go tool air`) to rebuild and restart the server on file changes.

3. In another terminal, install frontend deps and run:

```sh
cd frontend
pnpm install
pnpm dev
```

You can access the frontend over `http://localhost:5173`.

---

To reset local database data after changing its user, password, or database name:

```sh
podman compose down -v
```

This permanently deletes the local database volume.

## Production

### Local production demo

Build and run the same frontend, backend, migration, and Caddy containers used in production, without a domain or HTTPS:

```sh
cp .env.example .env # once
podman compose -p demo --env-file .env up --build
```

Open `http://localhost:5173`. This is a production-like integration test, so it does not hot-reload. It uses a separate database volume from development.

### Deploy on a domain

On the deployment server, point `DOMAIN` to the server and allow inbound ports 80 and 443.

1. Create the production configuration:

```sh
cp .env.production.example .env.production
```

2. Edit `.env.production`: set `DOMAIN` and replace the JWT and PostgreSQL secret placeholders with strong values. PostgreSQL passwords must be URL-safe.

3. Build and start the stack:

```sh
podman compose --env-file .env.production up -d --build
```

4. Verify it:

```sh
podman compose --env-file .env.production ps
curl -fsS https://your-domain.example/api/health
```

Open `https://your-domain.example` after the health check returns `OK`. The first start applies pending database migrations before starting the backend.

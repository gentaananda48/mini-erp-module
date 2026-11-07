# miniERPModule

A full-stack Mini ERP application for managing products, inventory
movements, and customer orders.

## Table of Contents

- [Overview](#overview)
- [Technology Stack](#technology-stack)
- [Application Architecture](#application-architecture)
- [Project Structure](#project-structure)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Run the Application with Docker](#run-the-application-with-docker)
- [Access the Application](#access-the-application)
- [Useful Docker Commands](#useful-docker-commands)
- [Run the Frontend Locally](#run-the-frontend-locally)
- [Run the Backend Locally](#run-the-backend-locally)
- [API Overview](#api-overview)
- [Database and Seed Data](#database-and-seed-data)
- [Troubleshooting](#troubleshooting)

## Overview

The application separates responsibilities into:

- **Frontend:** user interface, form handling, client-side validation,
  and API integration.
- **Backend:** business rules, request validation, transaction
  management, and standardized API responses.
- **Database:** persistent storage for products, stock movements,
  orders, and related records.

The backend is responsible for authoritative business calculations. For
example, product prices used in an order are retrieved from the
database, and order totals are calculated by the backend rather than
trusted from client input.

## Technology Stack

### Frontend

- React
- TypeScript
- Vite
- HTML and CSS
- Nginx for serving the production build and proxying API requests
- Docker multi-stage build

### Backend

- Go
- REST API
- PostgreSQL
- Database migrations
- Docker

### Infrastructure

- Docker
- Docker Compose
- PostgreSQL 17 Alpine
- `golang-migrate` migration image
- PostgreSQL client image for seed execution
- Nginx

## Project Structure

The following is the intended high-level structure. Some internal Go
package and React component names may differ as the project evolves.

```text
miniERPModule/
├── backend/
│   ├── cmd/
│   ├── internal/
│   ├── migrations/
│   ├── seeds/
│   │   ├── products.sql
│   │   ├── stock_movements.sql
│   │   └── orders.sql
│   ├── DockerFile
│   └── go.mod
│
├── frontend/
│   ├── public/
│   ├── src/
│   │   ├── features/
│   │   │   ├── products/
│   │   │   └── orders/
│   │   ├── components/
│   │   ├── lib/
│   │   ├── styles/
│   │   ├── types/
│   │   ├── App.tsx
│   │   └── main.tsx
│   ├── Dockerfile
│   ├── nginx.conf
│   ├── .dockerignore
│   ├── package.json
│   └── package-lock.json
│
└── docker-compose.yml
```

## Features

### Products and Inventory

- View products in a paginated table.
- Search products by SKU or product name.
- Filter products by status and low-stock condition.
- View product details.
- Create products.
- Edit product SKU, name, and price.
- Delete products, subject to backend business rules.
- Record stock movements using **IN** and **OUT** adjustments.
- View stock movement history for a product.
- Preview the resulting stock before submitting an adjustment.

Stock is managed through stock movements rather than by directly editing
the product stock field in the product form.

### Orders

The project includes an Orders area intended for order listing and order
creation, with order detail and cancellation workflows supported by the
backend API.

Order creation follows backend-controlled business rules:

- The client submits product SKU and quantity, not unit prices.
- The backend retrieves product prices from the database.
- The backend calculates the order total.
- Stock changes for an order are handled transactionally.
- Invalid items, non-positive quantities, duplicate SKUs, and empty
  item lists should be rejected by the backend.

Order functionality should be considered complete only for the workflows
implemented in the current source code.

## Prerequisites

To run the complete application with Docker, install:

1.  **Git** --- to clone the repository.
2.  **Docker Desktop** (macOS/Windows) or Docker Engine (Linux).
3.  **Docker Compose v2** --- available as `docker compose`.

Verify the installations:

```bash
git --version
docker --version
docker compose version
```

For local development outside Docker, also install:

- Go (use a version compatible with `backend/go.mod`).
- Node.js and npm (Node.js 22 or a version compatible with the
  frontend project).

## Run the Application with Docker

### 1. Clone the repository

```bash
git clone <YOUR_REPOSITORY_URL>
cd miniERPModule
```

Replace `<YOUR_REPOSITORY_URL>` with the actual Git repository URL.

### 2. Check the project files

Before starting, make sure the repository contains:

```text
docker-compose.yml
backend/DockerFile
backend/migrations/
backend/seeds/
frontend/Dockerfile
frontend/nginx.conf
frontend/package.json
frontend/package-lock.json
```

The Compose configuration expects the backend Dockerfile to be named
`DockerFile` with the same capitalization.

### 3. Build and start the application

Run this command from the project root:

```bash
docker compose up --build -d
```

Docker Compose will:

1.  Start PostgreSQL.
2.  Wait for PostgreSQL's health check to pass.
3.  Run database migrations.
4.  Run the SQL seed files.
5.  Start the Go backend after the seeder completes successfully.
6.  Build and start the React frontend served by Nginx.

The first build may take several minutes depending on network speed and
whether Docker already has the required base images.

### 4. Check container status

```bash
docker compose ps
```

Check the logs if a service does not start:

```bash
docker compose logs -f postgres
docker compose logs -f migrate
docker compose logs -f seeder
docker compose logs -f backend
docker compose logs -f frontend
```

Press `Ctrl+C` to stop following logs. This does not stop the
containers.

## Access the Application

Once the services are running:

Component URL

---

Frontend http://localhost:3000
Backend http://localhost:8080
PostgreSQL localhost:5432

Open the frontend:

```text
http://localhost:3000
```

The frontend sends API requests to relative paths such as:

```text
/api/products?page=1&limit=5
```

Nginx proxies those requests to the backend container.

The backend port is published on the host for direct API testing with
tools such as curl or Postman.

## Useful Docker Commands

### Start existing containers

```bash
docker compose up -d
```

### Rebuild after changing application code or Docker configuration

```bash
docker compose up --build -d
```

### Follow all service logs

```bash
docker compose logs -f
```

### Follow logs for one service

```bash
docker compose logs -f frontend
docker compose logs -f backend
```

### Stop and remove containers and network

```bash
docker compose down
```

This command does **not** delete the named PostgreSQL volume. Database
data normally remains available when the stack is started again.

### Stop the stack and delete the database volume

```bash
docker compose down -v
```

**Warning:** `-v` deletes the PostgreSQL data volume. All data stored in
that volume will be lost. Use this only when you intentionally want to
reset the local database.

### Rebuild only the frontend

```bash
docker compose build frontend
docker compose up -d frontend
```

### Rebuild only the backend

```bash
docker compose build backend
docker compose up -d backend
```

## Run the Frontend Locally

Running the frontend locally is useful during UI development. PostgreSQL
and the backend can continue running in Docker.

Start the backend and database services first:

```bash
docker compose up -d postgres migrate seeder backend
```

Then open another terminal:

```bash
cd frontend
npm ci
npm run dev
```

Vite normally serves the development application at:

```text
http://localhost:5173
```

Make sure the Vite development proxy is configured to forward `/api`
requests to the backend at `http://localhost:8080`. The exact proxy
configuration is defined in `frontend/vite.config.ts`.

To create a production build locally:

```bash
npm run build
```

To preview the production build using Vite:

```bash
npm run preview
```

## Run the Backend Locally

Running the backend locally is optional. The Docker Compose workflow is
the recommended way to start the complete application.

Start PostgreSQL and initialize the database:

```bash
docker compose up -d postgres migrate seeder
```

Then open another terminal:

```bash
cd backend
go mod download
```

Start the backend using the repository's configured Go entry point. For
example, if the application entry point is `cmd/api/main.go`:

```bash
go run ./cmd/api
```

The exact command depends on the actual entry point in `backend/cmd/`.
Check the project source before using the example above.

When running Go directly on the host, configure the database host as
`localhost` (not `postgres`) and use the database port published by
Docker Compose:

```text
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=mini_erp
DB_SSLMODE=disable
APP_PORT=8080
```

These are the development values from the current Compose configuration.
Do not use the example credentials as production credentials.

## API Overview

The API is served under the `/api` prefix.

### Products

Method Endpoint Description

---

GET `/api/products?page=1&limit=10` Get paginated products
POST `/api/products` Create a product
PUT `/api/products/{id}` Update a product
DELETE `/api/products/{id}` Delete a product

Example create-product request:

```json
{
  "sku": "PRD-099",
  "name": "Mechanical Keyboard",
  "price": 1000
}
```

### Stock Movements

---

Method Endpoint Description

---

POST `/api/products/{id}/stock` Create an IN or OUT
stock movement

GET `/api/products/{id}/movements?page=1&limit=10` Get paginated stock
movement history

---

Example stock movement request:

```json
{
  "type": "in",
  "qty": 5,
  "note": "Restock from supplier"
}
```

Supported movement types:

- `in` --- increases stock.
- `out` --- decreases stock.

### Orders

Method Endpoint Description

---

POST `/api/orders` Create an order
GET `/api/orders/{id}` Get order detail
POST `/api/orders/{id}/cancel` Cancel an order

The order-list endpoint and its query parameters should be documented
here once confirmed against the implemented backend route and handler.

Order creation payloads should contain product identifiers/SKUs and
quantities according to the backend DTO. The backend is responsible for
retrieving unit prices and calculating totals.

## Database and Seed Data

The PostgreSQL database is initialized through three Compose services:

- `postgres`: runs PostgreSQL and stores data in the `postgres_data`
  named volume.
- `migrate`: applies SQL migrations from `backend/migrations`.
- `seeder`: executes SQL seed files from `backend/seeds` after
  migrations complete.

The current seeder command executes:

```text
backend/seeds/products.sql
backend/seeds/stock_movements.sql
backend/seeds/orders.sql
```

The seed scripts should be safe to execute whenever the Compose
initialization flow runs. Use idempotent SQL or another
duplicate-protection strategy where appropriate.

### Important: persistent data

Running:

```bash
docker compose down
docker compose up --build -d
```

preserves the named PostgreSQL volume. It does not automatically create
a clean database.

If seed scripts insert fixed records without duplicate protection,
running the seeder again against an existing database may create
duplicate data or fail due to unique constraints. Review the seed SQL
before repeatedly re-running initialization.

For a completely fresh local database, use:

```bash
docker compose down -v
docker compose up --build -d
```

This permanently deletes the local database volume before recreating it.

## Troubleshooting

### Frontend cannot reach the backend

Check:

1.  The backend container is running:

    ```bash
    docker compose ps
    ```

2.  The frontend Nginx configuration proxies `/api/` to
    `http://backend:8080`.

3.  The frontend API client uses a relative base URL such as `/api`.

4.  The backend listens on `0.0.0.0:8080` inside its container, rather
    than only `127.0.0.1`.

### Frontend container fails to build

Check that:

- `frontend/package-lock.json` exists.
- `npm ci` succeeds locally.
- `npm run build` succeeds locally.
- `frontend/Dockerfile` and `frontend/nginx.conf` exist.

Inspect the build output:

```bash
docker compose build frontend
```

### Backend container fails to start

Inspect:

```bash
docker compose logs backend
docker compose logs seeder
docker compose logs migrate
```

Verify database configuration and confirm that migrations and seeding
completed successfully.

### PostgreSQL data is not reset

`docker compose down` preserves named volumes. To delete the local
database volume and initialize a fresh database:

```bash
docker compose down -v
docker compose up --build -d
```

Only use this when deleting local database data is intended.

### Port is already in use

If ports `3000`, `8080`, or `5432` are occupied, stop the process using
the port or change the host-side port mapping in `docker-compose.yml`.

For example, changing:

```yaml
ports:
  - "3000:80"
```

to:

```yaml
ports:
  - "3001:80"
```

makes the frontend available at `http://localhost:3001`.

---

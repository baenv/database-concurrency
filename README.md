# Golang template for Data concurrency

## Components

- **1 Queue server** implemented by using `asyncq` that listen and process event.
- **1 Redis server** to back **Queue**.
- **2 Indexers** stream block from Ethereum and push event to Queue.
- **1 API server**
- **1 PostgreSQL Database**

## Run locally

### Dependencies

- Currently, I am using [EntGo](https://entgo.io/docs) to generate and manage model. Please check the doc.
- Modify `docker-compose` for enable/disable proper services.

### Environment

- Create file `.env` depend on `.env.example`.

### Commands

- `make up`: start all components using `docker-compose`.
- `make up-latest`: like `up` but re-build docker images.
- `make down`: stop all components.
- `make ent-gen`: generate new ent from definition at `ent/schema`
- `make migrate-new name=<name>`: check different and generate file migration with the given name.

## Usage

- Currently, I configured to get the following characteristics
  - **2 Indexer** subscribe to ETHEREUM to get newest block and push event to **Queue**.
  - **Queue** receive events, deduplicates and log them, have no interact to DB now here.
  - **API server** serves an API to get transaction by hash.

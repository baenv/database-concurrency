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
- `make ent-init name=<name>` init new schema.

## How can I migrate new change?
- Run `make migrate-new-custom` to generate empty migration file.
- Edit this file.
- Run `make mirage-new-hash` to generate new proper.
- Run `make migrate-latest` to apply latest migration.
- Update/Add ent schema.
- Run `make ent-gen` to generate latest version of ent.

- For custom migration, after edit content of new file, please use `atlas migrate hash --dir "file://<path-to-migrations>"` to re-generate the checksum. Then you need to migrate by using flag `--baseline 20230219100638`.
    Example:
      ```
        atlas migrate apply \
        --dir "file://ent/migrate/migrations" \
        -u "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable" --baseline 20230219100638
      ```

## Usage

- Currently, I configured to get the following characteristics
  - **2 Indexer** subscribe to ETHEREUM to get newest block and push event to **Queue**.
  - **Queue** receive events, deduplicates and log them, have no interact to DB now here.
  - **API server** serves an API to get transaction by hash.


## How can I simulate scenario
- Create folder to store all files of your scenario in the folder `/scenario`.
- Create folder `/scenario/your/testdata` included:
  - `seed.up.sql`: SQL query to seed data to test.
  - `seed.down.sql`: SQL query to remove seed data.
- Create file `scenario/your/test.go` to implement scenario by using `pewpew`.
- Run `test.go` to test

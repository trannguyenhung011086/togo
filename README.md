### Overview

This is a simple backend for a good old todo service, right now this service can handle login/list/create simple tasks.

### How to run on local

- run `source .env` to add environment variables
- run `sudo docker run -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 postgres` to start postgres
- run `go run main.go`
- run `go test internal/storages/postgres/db_test.go` to unit test with db methods

### DB Schema

```sql
-- users definition

CREATE TABLE IF NOT EXISTS users (
	id TEXT NOT NULL,
	password TEXT NOT NULL,
	max_todo INTEGER DEFAULT 5 NOT NULL,
	CONSTRAINT users_PK PRIMARY KEY (id)
);

INSERT INTO users (id, password, max_todo) VALUES('firstUser', 'example', 5);

-- tasks definition

CREATE TABLE IF NOT EXISTS tasks (
	id TEXT NOT NULL,
	content TEXT NOT NULL,
	user_id TEXT NOT NULL,
    created_date TEXT NOT NULL,
	CONSTRAINT tasks_PK PRIMARY KEY (id),
	CONSTRAINT tasks_FK FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### Sequence diagram

![auth and create tasks request](https://github.com/trannguyenhung011086/togo/blob/master/docs/sequence.svg)

# TeacherPortal

## Local Development Setup

### Tools

- Docker
- GoLang
- Homebrew
- Migrate

```sh
brew install golang-migrate
```

### Setup infrastructure

#### Approach 1

Go to the directory of the project

```sh
make dockerbuild

make network

make mysql

# Might need to wait for 10s-20s for db services to start to execute the below
make createdb

make migrateup

make dockerserver
```

If above does not work, try executing the actual commands

```sh
docker build -t teacherportal:latest .

docker network create teacherportal-network

docker run --network teacherportal-network -p 3306:3306 --name mysql8 -e MYSQL_ROOT_PASSWORD=secret -d mysql:8.0.31

docker exec -it mysql8 mysql -u root -p'secret' -e "create database teacherportal;"

migrate -path db/migration -database "mysql://root:secret@tcp(localhost:3306)/teacherportal" -verbose up

docker run --name teacherportal --network teacherportal-network -p 8080:8080 -e DB_SOURCE="root:secret@(mysql8:3306)/teacherportal?parseTime=true" teacherportal:latest
```

#### Approach 2

Go to the directory of the project

```sh
docker run -p 3306:3306 --name mysql8 -e MYSQL_ROOT_PASSWORD=secret -d mysql:8.0.31

# Alternatively can run the actual command or bash into the container to create the db manually
make createdb

# or manually execute the init up-script
make migrate up

# or `go run main.go`
make server
```

### Possible errors encounters

#### Unable to create db

```
ERROR 1045 (28000): Access denied for user 'root'@'localhost' (using password: YES)
make: *** [createdb] Error 1

ERROR 2002 (HY000): Can't connect to local MySQL server through socket '/var/run/mysqld/mysqld.sock' (2)
make: *** [createdb] Error 1
```

Might need to tear down the container and re-run it.
Or wait awhile and re run the command.
Else will have to bash it to create the db

```
docker exec -it mysql8 bin/bash
mysql -uroot -psecret
create database teacherportal;
```

## Some DB Design decisions

### Using BigInt for ID

Although student can be uniquely identified with email addresses. Used BigInt, as it would probably be easier to use integer ID for production support and logging of email might have to be masked.

### Audit fields

Generally audit fields just to track the creation time, suspended time, etc. for production support and audit trails.

## Assumptions

Some assumptions when developing the stories.

### Story 4

Notification inputs with no spaces between 'mentions' will not be recognised. As from example below:
User input might have typo
e.g. trying to mention '@john' and '@gmail.com@gmail.com'.
While @john is not valid, but it is beside a mentioned with '@gmail.com' could not really tell if it is suppose to be '@john@gmail.com' and '@gmail.com'

```
@john@gmail.com@gmail.com
```

## Challenges/Limitations

### SQLC does not fully support MySQL

- Chose SQLC as it is fast and relatively simple to use
- Able to construct SQL query and automatic generation of CRUD code, reduces error.
- And reduces some time in manual mapping
- And since SQLC will process the SQL query to generate the code, it will catch and surface query errors. Do not need to wait until runtime for the errors to be surfaced.
- However, SQLC do not have full support for MySQL. E.g. Queries with 'IN' clauses are not supported.
- Hence, had to 'customise' some queries. Breaking the usual development pattern.
- Could explore using GORM or just using database/sql package

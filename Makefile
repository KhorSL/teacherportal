mysql:
	docker run --network teacherportal-network -p 3306:3306 --name mysql8 -e MYSQL_ROOT_PASSWORD=secret -d mysql:8.0.31

createdb:
	docker exec -it mysql8 mysql -u root -p'secret' -e "create database teacherportal;"

dropdb:
	docker exec -it mysql8 mysql -u root -p'secret' -e "drop database teacherportal;"

migrateup:
	migrate -path db/migration -database "mysql://root:secret@tcp(localhost:3306)/teacherportal" -verbose up

migratedown:
	migrate -path db/migration -database "mysql://root:secret@tcp(localhost:3306)/teacherportal" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/khorsl/teacherportal/db/sqlc Store

server:
	go run main.go

network:
	docker network create teacherportal-network

dockerbuild:
	docker build -t teacherportal:latest .

dockerserver:
	docker run --name teacherportal --network teacherportal-network -p 8080:8080 -e DB_SOURCE="root:secret@(mysql8:3306)/teacherportal?parseTime=true" teacherportal:latest

.PHONY: mysql createdb dropdb migrateup migratedown test sqlc mock server network dockerbuild dockerserver
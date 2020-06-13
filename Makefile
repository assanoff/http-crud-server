    export DB_HOST=localhost
	export DB_USER=postgres
	export DB_PASS=password
	export DB_NAME=crud-db
	export PG_MAXCONN=10
    export PORT=3220
	export ENDPOINT=/api/v1
	export DB_SCHEMA=test

.PHONY: build
build:
	go build -o bin/http-crud-server -v ./cmd/apiserver

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.PHONY: run
run:
	
	./bin/http-crud-server/apiserver

.PHONY: rundb
rundb:
	docker run -d --rm --name postgresql -e POSTGRES_DB=crud-db -e POSTGRES_PASSWORD=password -p 5432:5432 postgres:11

.DEFAULT_GOAL := build


#docker ps | grep postgresql
#docker exec -it f797324eba028f psql -h localhost -U postgres crud-db -W


#создание curl -vvv -XPOST -H 'Content-Type: application/json' -d '{"name":"Alice","email":"Lisa@google.com"}'  localhost:3220/api/v1/users
#curl -vvv -XPOST -H 'Content-Type: application/json' -d '{"name":"Bob","email":"Boby@google.com"}'  localhost:3220/api/v1/users

# чтение: curl -vvv localhost:3220/api/v1/users/1
# чтение: curl -vvv http://localhost:3220/api/v1/users/?field=email&val=Lisa@google.com

# обновление: curl -vvv -XPUT -H 'Content-Type: application/json' -d '{"name":"Bob","email":"456@google.com"}' localhost:3220/api/v1/users/3

# удаление: curl -vvv -XDELETE  localhost:3220/api/v1/records/1
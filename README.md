# go-microservice

A simple CRUD HTTP API using Go and MongoDB.

## Run locally

```shell
$ go run server.go
```

## API

It's a Book API, where a book has the following fields:
- `id`
- `title`
- `author`

The endpoints are:
- `GET http://localhost:8080/`
- `GET http://localhost:8080/books` 
- `POST http://localhost:8080/books`
- `GET http://localhost:8080/books/{id}`
- `PATCH http://localhost:8080/books/{id}`
- `DELETE http://localhost:8080/books/{id}`

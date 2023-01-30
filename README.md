# Example user service

This example service is built for quick service prototyping. It consists of 3 layers:  
- GRPC handlers/HTTP handlers
- Service layer which acts as code responsible for problem domain
- Storage which is responsible only for storing given data

## Running

### With Docker-compose

First run postgres. It takes a while to initialize the script. Even with `depends_on` the server  
becomes available before the DB and then shuts down because the DB isn't ready.

```
docker-compose up postgres
```

When it's up and running run the rest of the service with:

```
docker-compose up -d
```

It will expose 4 services to the outside.
- postgres on it's default port 5432
- redis on it's default port 6379
- service GRPC listening on port 9000
- service HTTP proxy listening on port 9001

### Native

For native running postgres and redis have to be provided locally on their default ports.  
Database needs a user "user" with password "password" and database named "database".  
Database also needs to have a schema initialized. It can be found in `./db-initial-scripts`.

#### Run databases with docker-compose

The easiest way of running the databases is with docker-compose.

```
docker-compose up postgres redis -d
```

#### Running the app

When databases are running, simply run:

```
go run ./cmd/server/main.go
```

## Testing
Run tests with:
```
go test ./...
```

Current test coverage can be better, but given the big scope of the project I have omitted some tests.

## Building proto

Make sure that protoc is available on your OS.

```
brew install protobuf
```

proto-gen-go and gateway as well:
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
```

Make sure that googleapis git submodule is initialized.

```
git submodule init
git submodule update

make gen_proto
```

## API docs

API documentation is generated from proto files and is available in `./proto/users.swagger.json` and `./proto/health.swagger.json`.

### googleapis

Googleapis module is required for building the protobuf files because of GRPC gateway.

### Testing database

Testing database is out of scope as it would take me too much time.


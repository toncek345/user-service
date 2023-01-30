## Build
FROM golang:alpine AS build

WORKDIR /app

COPY . .
RUN go mod download

WORKDIR /app/cmd/server
RUN go build

## Run
FROM golang:alpine

WORKDIR /app
COPY --from=build /app/cmd/server/server .

ENTRYPOINT ./server

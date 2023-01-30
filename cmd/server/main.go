package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/toncek345/userservice/server"
	"github.com/toncek345/userservice/service"
	"github.com/toncek345/userservice/storage"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	dbOpts := "host=localhost user=user password=password dbname=database sslmode=disable"
	if os.Getenv("ENV") == "compose" {
		dbOpts = "host=postgres user=user password=password dbname=database sslmode=disable"
	}

	db, err := sqlx.Open("postgres", dbOpts)
	if err != nil {
		log.Fatal("db is not available")
	}

	userStorage := &storage.UserStorageSQL{
		DB: db,
	}
	userService := &service.UserServiceImpl{
		UserStorage: userStorage,
	}

	s, err := server.NewServer(9000, 9001, userService)
	if err != nil {
		log.Fatalf("new server: %s", err)
	}

	go func() {
		log.Println("Starting grpc server")
		log.Printf("server run exited: %s\n", s.Start())
	}()
	// TODO: there should be some synchronization between because GRPC needs to be running in order to be reachable from reverse proxy.
	time.Sleep(500 * time.Millisecond)
	go func() {
		log.Println("Starting http server")
		log.Printf("server run exited: %s\n", s.StartHTTP())
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT)

	<-signalChan
	log.Println("shutting down...")
	s.Stop()
}

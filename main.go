package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/bmsandoval/vehicle-api/internal/transport/api_grpc"
	"github.com/bmsandoval/vehicle-api/internal/transport/api_http"
	"github.com/bmsandoval/vehicle-api/pkg/appcontext"
	"github.com/bmsandoval/vehicle-api/pkg/services"
	gocache "github.com/pmylund/go-cache"
	"github.com/soheilhy/cmux"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"time"
)

var Viper *viper.Viper

func init() {
	Viper = viper.New()
	Viper.AutomaticEnv()
}

func main() {
	serverPort, ok := Viper.Get("SERVER_PORT").(string)
	if serverPort == "" || ! ok {
		log.Fatal("env variable SERVER_PORT missing or invalid")
	}

	appCtx := appcontext.Context{
		Viper: Viper,
		GoContext: context.Background(),
	}

	// Create the main listener.
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", serverPort))
	if err != nil {
		log.Fatal(err)
	}

	// Create a cmux.
	m := cmux.New(l)

	// Match connections in order:
	// First grpc, then HTTP
	grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpL := m.Match(cmux.HTTP1Fast())


	pgHost, ok := Viper.Get("POSTGRES_HOST").(string)
	if pgHost == "" || ! ok {
		log.Fatal("env variable POSTGRES_HOST missing or invalid")
	}
	pgPort, ok := Viper.Get("POSTGRES_PORT").(string)
	if pgPort == "" || ! ok {
		log.Fatal("env variable POSTGRES_PORT missing or invalid")
	}
	pgUser, ok := Viper.Get("POSTGRES_USER").(string)
	if pgUser == "" || ! ok {
		log.Fatal("env variable POSTGRES_USER missing or invalid")
	}
	pgPass, ok := Viper.Get("POSTGRES_PASSWORD").(string)
	if pgPass == "" || ! ok {
		log.Fatal("env variable POSTGRES_PASS missing or invalid")
	}
	pgDB, ok := Viper.Get("POSTGRES_DB").(string)
	if pgDB == "" || ! ok {
		log.Fatal("env variable POSTGRES_DB missing or invalid")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		pgHost, pgPort, pgUser, pgPass, pgDB)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		panic(err)
	}

	grpcS, httpS, err := AcquireServers(appCtx, db)

	// Use the muxed listeners for your servers.
	go grpcS.Serve(grpcL)
	go httpS.Serve(httpL)

	// Start serving!
	m.Serve()
}

func AcquireServers(appCtx appcontext.Context, db *sql.DB) (*grpc.Server, *http.Server, error) {
	cache := gocache.New(5*time.Minute, 30*time.Second)

	serviceBundle, err := services.NewBundle(appCtx, db, cache)
	if err != nil {
		return nil, nil, err
	}

	grpcS := api_grpc.AcquireVehicleServer(appCtx, *serviceBundle)
	httpS := api_http.AcquireVehicleServer(appCtx, *serviceBundle)

	return grpcS, httpS, nil
}

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

const sqlConnectionStringVar = "PG_CONNECTION_STRING"
const listenAddress = ":8080"

func main() {
	// First, get the connection string and attempt a connection or fail
	connStr := os.Getenv(sqlConnectionStringVar)
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open the database %s=%s: %v", sqlConnectionStringVar, connStr, err)
	}
	defer conn.Close()

	// Second, wait until we get a successful connection to the db
	if err := conn.PingContext(context.TODO()); err != nil {
		log.Fatalf("failed to connect to database %s=%s: %v", sqlConnectionStringVar, connStr, err)
	}
	log.Printf("successfully connected to database %s=%s", sqlConnectionStringVar, connStr)

	// Third, link the handler function
	http.HandleFunc("GET /{$}", func(writer http.ResponseWriter, request *http.Request) {
		var versionOutput string
		if err := conn.QueryRowContext(request.Context(), "SELECT version()").Scan(&versionOutput); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(writer, "Failed to get sql version: %v", err)
			log.Printf("%s %s %s status=500: %v", request.Host, request.Method, request.RequestURI, err)
		} else {
			_, _ = fmt.Fprintf(writer, "SQL VERSION: %s", versionOutput)
			log.Printf("%s %s %s status=200", request.Host, request.Method, request.RequestURI)
		}
	})

	// Finally, run the http server
	log.Print("starting server")
	if err := http.ListenAndServe(listenAddress, http.DefaultServeMux); err != nil {
		log.Fatal(err.Error())
	}
}

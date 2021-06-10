package main

import (
	"astral/app"
	"astral/auth"
	"astral/storage"
	"flag"
	"log"
)

func main() {
	var ip, port, cert, key, users, data string
	var fillDB, clearDB bool
	flag.StringVar(&ip, "ip", "127.0.0.1", "ip address to listen on")
	flag.StringVar(&port, "port", "8083", "port to listen on")
	flag.StringVar(&cert, "cert", "", "TLS certificate file")
	flag.StringVar(&key, "key", "", "TLS key file")
	flag.StringVar(&users, "users", "users.conf", "file with user credentials")
	flag.StringVar(&data, "data", "sample.db", "sqlite3 database file")
	flag.BoolVar(&fillDB, "fill", false, "fill database with random sample data")
	flag.BoolVar(&clearDB, "clear", false, "remove all items from database")
	flag.Parse()
	auth, err := auth.New(users)
	if err != nil {
		log.Println("error creating authenticator:", err)
		return
	}
	st, err := storage.OpenSQLiteStorage(data)
	if err != nil {
		log.Println("error opening storage:", err)
		return
	}
	if clearDB {
		if err := storage.DeleteAllFromStorage(st); err != nil {
			log.Println("error deleting DB content.")
		}
		log.Println("Cleanup complete. Storage is now empty.")
		return
	}
	if fillDB {
		if err := storage.GenerateRandomPayload(st, 10); err != nil {
			log.Println("error generating random DB content.")
		}
		log.Println("Storage populated with 10 elements with random values.")
		return
	}
	srv, err := app.New(ip, port, st, auth)
	if err != nil {
		log.Println("error initializing app:", err)
		return
	}
	if cert != "" && key != "" {
		log.Fatal(srv.ListenAndServeTLS(cert, key))
	} else {
		log.Fatal(srv.ListenAndServe())
	}
}

package app

import (
	"astral/auth"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"astral/storage"
)

// Application represents the test web app
type Application struct {
	storage.Storage
	auth.Authenticator
}

// New accepts any compliant storage and authenticator and returns
// http server ready to run. Since application is using http Authorization: Basic
// it is not safe to launch server with ListenAndServe. ListenAnsServeTLS should be used.
func New(ip, port string, st storage.Storage, auth auth.Authenticator) (*http.Server, error) {
	if st == nil || auth == nil {
		return nil, errors.New("st or auth is nil")
	}
	srv := &http.Server{Addr: ip + ":" + port}
	app := Application{st, auth}
	http.Handle("/get", http.HandlerFunc(app.handleGet))
	http.Handle("/put", http.HandlerFunc(app.handlePut))
	http.Handle("/", http.HandlerFunc(redirectToGet))
	//Handling OS signals
	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-interrupts
		log.Println("exiting on signal:", sig)
		if err := st.Close(); err != nil {
			log.Printf("error closing storage: %v\n", err)
		}
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("server shutdown error: %v\n", err)
		}
	}()
	return srv, nil
}

func (app *Application) handleGet(w http.ResponseWriter, r *http.Request) {
	if !app.isValidUser(w, r) {
		return
	}
	if r.Method != "GET" {
		allowMethod(w, r, "GET")
	}
	elements, err := app.Get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	htmlGet.Execute(w, PageGet{elements})
}

func (app *Application) handlePut(w http.ResponseWriter, r *http.Request) {
	if !app.isValidUser(w, r) {
		return
	}
	if r.Method == "GET" {
		username, _, _ := r.BasicAuth()
		htmlPut.Execute(w, PagePut{username})
		return
	} else if r.Method == "POST" {
		payload := r.FormValue("payload")
		if payload == "" {
			http.Error(w, "Can not add empty element", http.StatusBadRequest)
			return
		}
		id, err := app.Put(storage.Element{Payload: payload})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		e, _ := app.Get(id)
		w.WriteHeader(http.StatusCreated)
		htmlOK.Execute(w, PageOK{e[0]})
		return
	}
	allowMethod(w, r, "GET")
}

func (app *Application) isValidUser(w http.ResponseWriter, r *http.Request) bool {
	username, password, _ := r.BasicAuth()
	if app.CredentialsAreValid(username, password) {
		return true
	}
	w.Header().Set(`WWW-Authenticate`, `Basic realm="Log in to add content", charset="UTF-8"`)
	w.WriteHeader(http.StatusUnauthorized)
	return false

}

func redirectToGet(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/get", http.StatusSeeOther)
}

func allowMethod(w http.ResponseWriter, r *http.Request, allowed string) {
	w.Header().Add("Allowed", allowed)
	http.Error(w, "405 Method not allowed", http.StatusMethodNotAllowed)
}

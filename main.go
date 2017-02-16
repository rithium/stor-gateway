package main

import (
	"./handler"
	"flag"
	"log"
	"os"
	"github.com/rithium/version"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
	"fmt"
	"time"
	"./model"
)

func init() {
	versionFlag := flag.Bool("v", false, "prints version")
	//configFlag := flag.Bool("c", false, "dumps configuration")

	flag.Parse()

	if *versionFlag {
		log.Println("Stor Data", version.GetVersion())
		os.Exit(0)
	}

	/*config.LoadConfig()

	if *configFlag {
		log.Printf("HTTP:\t%+v\n", config.HttpServer)
		log.Printf("Cassandra:\t%+v\n", config.Cassandra)

		os.Exit(0)
	}*/
}

func main() {
	hub := model.NewHub()

	go hub.Run()

	router := mux.NewRouter()

	router.HandleFunc("/data/{nodeId:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		handler.HandlePostData(hub, w, r)
	}).Methods("POST")

	router.HandleFunc("/data/{nodeId:[0-9]+}", handler.HandleGetData).Methods("GET")

	router.HandleFunc("/ws/data/{apiKey}", func(w http.ResponseWriter, r *http.Request) {
		model.ServeWs(hub, w, r)
	})

	router.HandleFunc("/health", handleHealth)

	n := negroni.New()

	// Convert panics to 500 responses
	n.Use(negroni.NewRecovery())

	// Pretty print REST requests
	//n.Use(negroni.NewLogger())

	n.UseHandler(router)

	addr := fmt.Sprintf("%s:%d", "0.0.0.0", 80)

	serv := &http.Server{
		Addr:           addr,
		Handler:        n,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("Binding HTTP on", addr)

	log.Fatal("http serv:", serv.ListenAndServe())
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

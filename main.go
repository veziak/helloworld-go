package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/hello/{username}", getBirthdayMessage).Methods("GET")
	router.HandleFunc("/hello/{username}", createOrUpdateUser).Methods("PUT")
	router.HandleFunc("/healthcheck", healthcheck).Methods("GET")
	router.HandleFunc("/version", version).Methods("GET")
	fmt.Println("Server listen on port 8081")
	log.Fatal(http.ListenAndServe(":8081", router))

}

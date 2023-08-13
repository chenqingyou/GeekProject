package main

import (
	"fmt"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world %v", r.URL.Path[1:])
}

func user(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world %v", r.URL.Path[1:])
}
func queryParams(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	//name := values["name"][0]
	fmt.Fprintf(w, "values := [%v]", values)
}
func main() {
	http.HandleFunc("/home", home)
	http.HandleFunc("/user", user)
	http.HandleFunc("/url/query", queryParams)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

package main

import (
	"fmt"
	"log"
	"net/http"

	recovermw "github.com/Dayanand-Chinchure/gophercises/recover_middleware/recover"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	mux.HandleFunc("/debug/", recovermw.SourceCodeHandler)
	mux.HandleFunc("/panic/", recovermw.Panic)

	log.Fatal(http.ListenAndServe(":3000", recovermw.Middleware(mux)))
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello Jarvis !</h1>")
}

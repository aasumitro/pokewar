package main

import (
	"github.com/aasumitro/pokewar/resources"
	"log"
	"net/http"
)

func init() {
	// TODO:
	// LOAD CONFIG
	// ETC
}

func main() {
	http.Handle("/", http.FileServer(http.FS(resources.Resource)))
	log.Fatal(http.ListenAndServe(":8000", nil))
}

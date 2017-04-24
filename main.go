package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "We're getting married!\nWatch this space :)")
	})

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}

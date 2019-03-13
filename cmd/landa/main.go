package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/agonzalezro/landa/pkg/api"
)

func main() {
	addr := ":8080"
	if port, ok := os.LookupEnv("PORT"); ok {
		addr = fmt.Sprintf(":%s", port)
	}

	r := api.Install(mux.NewRouter())

	log.Println("Listening on addr", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

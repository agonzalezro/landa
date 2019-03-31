package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/agonzalezro/landa/pkg/api"
	"github.com/agonzalezro/landa/pkg/cluster"
)

func buildAddress(addr string) string {
	if port, ok := os.LookupEnv("PORT"); ok {
		addr = fmt.Sprintf(":%s", port)
	}
	return addr
}

func main() {
	kubeconfigFlag := flag.String("f", "", "Path to a kube config (usually: ~/.kube/config). Only required if you are running out of a cluster.")
	addrFlag := flag.String("a", ":8080", "Address where to listen, default :8080. It can also be set with the env var PORT.")
	flag.Parse()

	cluster, err := cluster.New(*kubeconfigFlag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	api := api.New(cluster)
	r := api.RegisterHandlers(mux.NewRouter())

	addr := buildAddress(*addrFlag)
	log.Println("Listening on addr", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

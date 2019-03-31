package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Cluster interface {
	Deployer
}

type Deployer interface {
	DeployFunction(context.Context, string, string) (string, error)
}

type LandaAPI struct {
	Cluster   Cluster
	Functions map[string]Landa
}

func New(cluster Cluster) LandaAPI {
	return LandaAPI{
		Cluster:   cluster,
		Functions: make(map[string]Landa),
	}
}

func (api *LandaAPI) CreateFunction(w http.ResponseWriter, r *http.Request) {
	var f Landa

	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	f.ID = hashCode(f.Code)
	if _, ok := api.Functions[f.ID]; ok {
		w.WriteHeader(http.StatusConflict)
		return
	}

	api.Functions[f.ID] = Landa{ID: f.ID, Code: f.Code}

	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(f); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	go func() {
		url, err := api.Cluster.DeployFunction(context.TODO(), f.ID, f.Code)
		if err != nil {
			log.Println(err)
		}

		f := api.Functions[f.ID]
		f.URL = url
		api.Functions[f.ID] = f
	}()
}

func (api *LandaAPI) GetFunctionByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	f, ok := api.Functions[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(f); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (api *LandaAPI) RegisterHandlers(r *mux.Router) *mux.Router {
	r.HandleFunc("/functions", api.CreateFunction).Methods(http.MethodPost)
	r.HandleFunc("/functions/{id}", api.GetFunctionByID).Methods(http.MethodGet)
	return r
}

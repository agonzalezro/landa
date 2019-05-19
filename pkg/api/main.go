package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"lambda-control-plane/pkg/model"
)

const contentTypeJSON = "application/json"

type Cluster interface {
	Deployer
}

type Deployer interface {
	DeployFunction(context.Context, model.Landa) error
	GetFunctionUrl(context.Context, string) (string, error)
}

type LandaAPI struct {
	Cluster   Cluster
	Functions map[string]model.Landa
}

func New(cluster Cluster) LandaAPI {
	return LandaAPI{
		Cluster:   cluster,
		Functions: make(map[string]model.Landa),
	}
}

func (api *LandaAPI) CreateFunction(w http.ResponseWriter, r *http.Request) {
	var functionMetaData model.Landa

	if err := json.NewDecoder(r.Body).Decode(&functionMetaData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	functionMetaData.ID = hashCode(functionMetaData.Code)
	if _, ok := api.Functions[functionMetaData.ID]; ok {
		w.WriteHeader(http.StatusConflict)
		return
	}

	api.Functions[functionMetaData.ID] = model.Landa{ID: functionMetaData.ID, Code: functionMetaData.Code}

	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(functionMetaData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	go func() {
		err := api.Cluster.DeployFunction(context.TODO(), functionMetaData)
		if err != nil {
			log.Println(err)
			return
		}
		//TODO need to wait until LB has been created and IP published. Now wait some seconds
		time.Sleep(5 * time.Second)
		url, err := api.Cluster.GetFunctionUrl(context.TODO(), functionMetaData.ID)
		if err != nil || url == "" {
			log.Println(err)
			return
		}
		f := api.Functions[functionMetaData.ID]
		f.URL = url
		api.Functions[f.ID] = f

		log.Printf("[DEBUG] Stored '%s' as function '%s' ingress IP", url, functionMetaData.ID)
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

func (api *LandaAPI) CallFunctionByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	f, ok := api.Functions[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// TODO: Port should be configurable
	url := fmt.Sprintf("http://%s:9443", f.URL)

	log.Printf("[DEBUG] Function '%s' serving on '%s'\n", id, url)

	resp, err := http.DefaultClient.Post(url, contentTypeJSON, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
	bs, _ := ioutil.ReadAll(resp.Body) // TODO: we don't always have a body, this should be enough for now
	w.Write(bs)
}

func (api *LandaAPI) RegisterHandlers(r *mux.Router) *mux.Router {
	r.HandleFunc("/functions", api.CreateFunction).Methods(http.MethodPost)
	r.HandleFunc("/functions/{id}", api.GetFunctionByID).Methods(http.MethodGet)
	r.HandleFunc("/functions/{id}:call", api.CallFunctionByID).Methods(http.MethodPost)
	return r
}

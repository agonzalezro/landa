package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type FunctionsAPI struct {
	Functions map[string]string
}

func NewFunctionsAPI() FunctionsAPI {
	return FunctionsAPI{
		Functions: make(map[string]string),
	}
}

func (api *FunctionsAPI) CreateFunction(w http.ResponseWriter, r *http.Request) {
	var f Function

	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	f.ID = hashCode(f.Code)
	if _, ok := api.Functions[f.ID]; ok {
		w.WriteHeader(http.StatusConflict)
		return
	}

	api.Functions[f.ID] = f.Code

	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(f); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (api *FunctionsAPI) GetFunctionByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	code, ok := api.Functions[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(Function{
		ID:   id,
		Code: code,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func Install(r *mux.Router) *mux.Router {
	api := NewFunctionsAPI()
	r.HandleFunc("/functions", api.CreateFunction).Methods(http.MethodPost)
	r.HandleFunc("/functions/{id}", api.GetFunctionByID).Methods(http.MethodGet)
	return r
}

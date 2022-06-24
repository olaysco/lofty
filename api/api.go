package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/olaysco/evolve/util"
)

type Api struct {
	DB     *sql.DB
	Router *mux.Router
	config util.Config
}

type Response struct {
	Data interface{} `json:"data"`
}

type PaginatedResponse struct {
	Data  interface{} `json:"data"`
	Count int         `json:"count"`
	Page  int         `json:"current_page"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewAPi(config util.Config, db *sql.DB) *Api {
	api := &Api{
		DB:     db,
		config: config,
	}

	api.initializeRoutes()
	return api
}

func (a *Api) initializeRoutes() {
	a.Router = mux.NewRouter()

	a.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		a.JSON(w, http.StatusOK, "OK")
	})
	a.Router.HandleFunc("/users", a.GetUsers).Methods("GET")
	a.Router.HandleFunc("/users/{email}", a.GetUserByEmail).Methods("GET")
}

//JSON response helper function
func (a *Api) JSON(w http.ResponseWriter, status int, data interface{}) {
	response, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

//Run http server i.e. start listening on server address
func (a *Api) Run() error {
	return http.ListenAndServe(a.config.HTTPServerAddress, a.Router)
}

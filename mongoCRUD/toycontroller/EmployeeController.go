// Package toycontroller Employee mongo CRUD
//
// The purpose of this application is to perform CRUD operation on mongo db
//
//     BasePath: /Employee
//     Version: 1.0.0
//
//     Contact: Dayanand C <dayanand.chinchure@gslab.com>
//
//     Consumes:
//       - application/json
//
//     Produces:
//       - application/json
// swagger:meta
package toycontroller

//go:generate swagger generate spec -m -o ./swagger.json

import (
	"net/http"

	"github.com/gorilla/mux"
)

var listenAndServe = http.ListenAndServe

//RunController testing
func RunController() {
	router := mux.NewRouter()
	router.HandleFunc("/employee", getEmployee).Methods("GET")
	router.HandleFunc("/employee", createEmployee).Methods("POST")
	router.HandleFunc("/employee", deleteEmployee).Methods("DELETE")
	router.HandleFunc("/employee", updateEmployee).Methods("PUT")
	listenAndServe(":8081", router)
}

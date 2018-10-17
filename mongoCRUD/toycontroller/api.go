package toycontroller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"toy_jarvis/toydb"
	toymodel "toy_jarvis/toymodel"

	"gopkg.in/mgo.v2/bson"
)

var dbError = false

// swagger:operation POST /employee Employee createEmployee
//
//  Create Employee
//
// Returns success code
//
// ---
// produces:
// - application/json
// parameters:
// - name: employee
//   in: body
//   description: Employee data
//   required: true
//   schema:
//     "$ref": "#/definitions/Employee"
// responses:
//   '200':
//     description: User Created
//     schema:
//       "$ref": "#/definitions/Employee"
//   '400':
//	   description: Bad Request
//   '403':
//     description: Forbidden, you are not allowed to undertake this operation
//   '405':
//     description: Method Not Allowed, likely url is not correct
func createEmployee(w http.ResponseWriter, req *http.Request) {
	c := toydb.Createconnection()
	var employee toymodel.Employee
	if err := json.NewDecoder(req.Body).Decode(&employee); err != nil {
		log.Printf("Decode Failed, ERROR ::%v\n ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := c.Insert(&employee)
	if err != nil || dbError {
		fmt.Println("Error obtained")
		http.Error(w, "Create employee failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("User Created ")
	return
}

// swagger:operation GET /employee Employee getEmployee
//
// Get Employee
//
// Returns Employee data for employee name provided
//
// ---
// produces:
// - application/json
// parameters:
// - name: name
//   in: query
//   description: Name of employee
//   required: true
//   type: string
// responses:
//   '200':
//     description: OK
//     schema:
//       "$ref": "#/definitions/Employee"
//   '400':
//	   description: Bad Request
//   '403':
//     description: Forbidden, you are not allowed to undertake this operation
//   '405':
//     description: Method Not Allowed, likely url is not correct
func getEmployee(w http.ResponseWriter, req *http.Request) {
	queryValues := req.URL.Query()
	c := toydb.Createconnection()
	result := toymodel.Employee{}

	err := c.Find(bson.M{"name": queryValues.Get("name")}).One(&result)
	if err != nil {
		fmt.Println("Error obtained")
		http.Error(w, "Get employee failed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result)
	return
}

// swagger:operation PUT /employee Employee updateEmployee
//
// Update Employee
//
// Update Employee data for employee name provided
//
// ---
// produces:
// - application/json
// parameters:
// - name: name
//   in: query
//   description: Name of employee
//   required: true
//   type: string
// - name: employeelocalhost:8081/employee?name=Rocky
//   in: body
//   description: Employee data
//   required: true
//   schema:
//     "$ref": "#/definitions/Employee"
// responses:
//   '200':
//     description: Employee Updated
//     schema:
//       "$ref": "#/definitions/Employee"
//   '400':
//	   description: Bad Request
//   '403':
//     description: Forbidden, you are not allowed to undertake this operation
//   '405':
//     description: Method Not Allowed, likely url is not correct
func updateEmployee(w http.ResponseWriter, req *http.Request) {
	c := toydb.Createconnection()
	queryValues := req.URL.Query()
	result := toymodel.Employee{}

	//Data Validation
	query := c.Find(bson.M{"name": queryValues.Get("name")})
	err := query.One(&result)
	if err != nil || dbError {
		fmt.Println("Error obtained")
		http.Error(w, "Get employee details failed", http.StatusInternalServerError)
		return
	}

	var employee toymodel.Employee
	if err := json.NewDecoder(req.Body).Decode(&employee); err != nil {
		log.Printf("Decode Failed, ERROR ::%v\n ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.Update(bson.M{"name": queryValues.Get("name")}, &employee); err != nil {
		fmt.Println("Error obtained")
		http.Error(w, "Update employee failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("User Updated !")
	return

}

// swagger:operation DELETE /employee Employee deleteEmployee
//
// Delete Employee
//
// Delete Employee data for employee name provided
//
// ---
// produces:
// - application/json
// parameters:
// - name: name
//   in: query
//   description: Name of employee
//   required: true
//   type: string
// responses:
//   '200':
//     description: OK
//     schema:
//       "$ref": "#/definitions/Employee"
//   '400':
//	   description: Bad Request
//   '403':
//     description: Forbidden, you are not allowed to undertake this operation
//   '405':
//     description: Method Not Allowed, likely url is not correct
func deleteEmployee(w http.ResponseWriter, req *http.Request) {
	c := toydb.Createconnection()
	queryValues := req.URL.Query()
	result := toymodel.Employee{}
	err := c.Find(bson.M{"name": queryValues.Get("name")}).One(&result)
	if err != nil && dbError {
		fmt.Println("Error obtained")
		http.Error(w, "DELETE get employee details failed", http.StatusInternalServerError)
		return
	}

	if queryValues.Get("name") != result.Name {
		fmt.Println("Error obtained")
		http.Error(w, "Delete employee failed", http.StatusInternalServerError)
		return
	}

	err = c.Remove(bson.M{"name": queryValues.Get("name")})
	if err != nil || dbError {
		fmt.Println("Error obtained")
		http.Error(w, "Delete employee failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("User Deleted !")
	return
}

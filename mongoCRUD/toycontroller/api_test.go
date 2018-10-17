package toycontroller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	toymodel "toy_jarvis/toymodel"
)

type Toy struct {
	I, J int
}

func init() {
	tmp := listenAndServe
	listenAndServe = func(port string, handler http.Handler) error {
		return nil
	}
	defer func() {
		listenAndServe = tmp
	}()
	RunController()
}

func changeDbError() {
	dbError = false
}
func TestCreateEmployee(t *testing.T) {
	defer changeDbError()

	dbError = false
	result, _ := json.Marshal(toymodel.Employee{ID: 1, Name: "Test", Dept: "IT1", Address: "Pune"})
	req, err := http.NewRequest("POST", "/employee", bytes.NewBuffer(result))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createEmployee)
	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusOK {
		fmt.Println("Test create passed !")
	}

	dbError = true
	result, _ = json.Marshal(toymodel.Employee{ID: 1, Name: "Test", Dept: "IT", Address: "Pune"})
	req, err = http.NewRequest("POST", "/employee", bytes.NewBuffer(result))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(createEmployee)
	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusOK {
		fmt.Println("Test create passed !")
	}

	dbError = true
	result, _ = json.Marshal("{akakakak")
	req, err = http.NewRequest("POST", "/employee", bytes.NewBuffer(result))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(createEmployee)
	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusOK {
		fmt.Println("Test create passed !")
	}

	dbError = false
	result, _ = json.Marshal("")
	req, err = http.NewRequest("POST", "/employee", bytes.NewBuffer(result))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(createEmployee)
	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusOK {
		fmt.Println("Test create passed !")
	}

}

func TestGetEmployee(t *testing.T) {
	defer changeDbError()
	dbError = true
	req, err := http.NewRequest("GET", "/employee", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("name", "Test")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getEmployee)
	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusOK {
		fmt.Println("Test get passed !")
	}

	dbError = false
	req, err = http.NewRequest("GET", "/employee", nil)
	if err != nil {
		t.Fatal(err)
	}
	q = req.URL.Query()
	q.Add("name", "TestNothing")
	req.URL.RawQuery = q.Encode()
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(getEmployee)
	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusOK {
		fmt.Println("Test get passed !")
	}

}

func TestUpdateEmployee(t *testing.T) {
	defer changeDbError()
	dbError = false
	result, _ := json.Marshal(toymodel.Employee{ID: 1, Name: "Test", Dept: "Dev", Address: "Pune"})
	req, err := http.NewRequest("PUT", "/employee", bytes.NewBuffer(result))
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("name", "Test")
	rr := httptest.NewRecorder()
	req.URL.RawQuery = q.Encode()
	handler := http.HandlerFunc(updateEmployee)
	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusOK {
		fmt.Println("Test update passed !")
	}

	//TEST TWO

	dbError = false
	result, _ = json.Marshal("{,aujdbsk")
	req, err = http.NewRequest("PUT", "/employee", bytes.NewBuffer(result))
	if err != nil {
		t.Fatal(err)
	}

	q = req.URL.Query()
	q.Add("name", "Test")
	rr = httptest.NewRecorder()
	req.URL.RawQuery = q.Encode()
	handler = http.HandlerFunc(updateEmployee)
	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusOK {
		fmt.Println("Test update passed !")
	}

	//TEST Three
	dbError = false
	result, _ = json.Marshal(toymodel.Employee{})
	req, err = http.NewRequest("PUT", "/employee", bytes.NewBuffer(result))
	if err != nil {
		t.Fatal(err)
	}

	q = req.URL.Query()
	q.Add("name", "TestNothing")
	rr = httptest.NewRecorder()
	req.URL.RawQuery = q.Encode()
	handler = http.HandlerFunc(updateEmployee)
	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusOK {
		fmt.Println("Test update passed !")
	}
}

func TestDeleteEmployee(t *testing.T) {
	defer changeDbError()

	dbError = false
	req, err := http.NewRequest("DELETE", "/employee", nil)

	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("name", "Testqqqqq")
	rr := httptest.NewRecorder()
	req.URL.RawQuery = q.Encode()
	handler := http.HandlerFunc(deleteEmployee)
	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusOK {
		fmt.Println("Test delete passed !")
	}

	dbError = true
	req, err = http.NewRequest("DELETE", "/employee", nil)

	if err != nil {
		t.Fatal(err)
	}

	q = req.URL.Query()
	q.Add("name", "Testasas")
	rr = httptest.NewRecorder()
	req.URL.RawQuery = q.Encode()
	handler = http.HandlerFunc(deleteEmployee)
	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusOK {
		fmt.Println("Test delete passed !")
	}

	dbError = true
	req, err = http.NewRequest("DELETE", "/employee", nil)

	if err != nil {
		t.Fatal(err)
	}

	q = req.URL.Query()
	q.Add("name", "Test")
	rr = httptest.NewRecorder()
	req.URL.RawQuery = q.Encode()
	handler = http.HandlerFunc(deleteEmployee)
	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusOK {
		fmt.Println("Test delete passed !")
	}

	dbError = false
	result, _ := json.Marshal(toymodel.Employee{ID: 1, Name: "Test", Dept: "IT1", Address: "Pune"})
	req, err = http.NewRequest("POST", "/employee", bytes.NewBuffer(result))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(createEmployee)
	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusOK {
		fmt.Println("Test create passed !")
	}

	dbError = false
	req, err = http.NewRequest("DELETE", "/employee", nil)

	if err != nil {
		t.Fatal(err)
	}

	q = req.URL.Query()
	q.Add("name", "Test")
	rr = httptest.NewRecorder()
	req.URL.RawQuery = q.Encode()
	handler = http.HandlerFunc(deleteEmployee)
	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusOK {
		fmt.Println("Test delete passed !")
	}
}

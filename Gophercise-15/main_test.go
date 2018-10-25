package main

import (
	"net/http"
	"testing"

	"github.ibm.com/dash/dash_utils/dashtest"
)

func TestHandlers(t *testing.T) {
	go main()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:3000/debug/", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error : %s route not ready", "http://localhost:3000/debug/")
	}
	if resp.StatusCode != 400 {
		t.Errorf("FAIL: %s returned status %d, expected %d", "http://localhost:3000/debug/", resp.StatusCode, 400)
	}

	client = &http.Client{}
	req, _ = http.NewRequest("GET", "http://localhost:3000/panic/", nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("Error : %s route not ready", "http://localhost:3000/panic/")
	}
	if resp.StatusCode != 200 {
		t.Errorf("FAIL: %s returned status %d, expected %d", "http://localhost:3000/panic/", resp.StatusCode, 200)
	}

	client = &http.Client{}
	req, _ = http.NewRequest("GET", "http://localhost:3000/hello", nil)
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("Error : %s route not ready", "http://localhost:3000/panic/")
	}
}

func TestMain(m *testing.M) {
	dashtest.ControlCoverage(m)
}

package main

import (
	"errors"
	"testing"

	"github.ibm.com/dash/dash_utils/dashtest"
)

func TestMain(m *testing.M) {
	main()
	dashtest.ControlCoverage(m)
}

func TestHome(t *testing.T) {
	tmp := homeDir
	homeDir = func() (string, error) {
		return "", errors.New("Error")
	}
	main()
	homeDir = tmp
}

package cobra

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.ibm.com/dash/dash_utils/dashtest"
)

func TestMain(m *testing.M) {
	dashtest.ControlCoverage(m)
}

func TestSetCmd(t *testing.T) {
	file, _ := os.OpenFile("test.txt", os.O_RDWR|os.O_CREATE, 0755)
	stdOut := os.Stdout
	os.Stdout = file
	setCmd.Run(setCmd, []string{"Test-Key", "Test-Value"})
	file.Seek(0, 0)
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Error("Expected nil but got", err)
	}
	val := strings.Contains(string(content), "Value set successfully")
	assert.Equalf(t, true, val, "Both should be equal")
	os.Stdout = stdOut
	defer os.Remove("test.txt")
}

func TestSetCmdNegative(t *testing.T) {
	file, _ := os.OpenFile("test.txt", os.O_RDWR|os.O_CREATE, 0755)
	stdOut := os.Stdout
	os.Stdout = file
	setCmd.Run(setCmd, []string{""})
	file.Seek(0, 0)
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Error("Expected nil but got", err)
	}
	val := strings.Contains(string(content), "Unable to add secret")
	assert.Equalf(t, true, val, "Both should be equal")
	os.Stdout = stdOut
	defer os.Remove("test.txt")
}

func TestGetCmdNoValueSet(t *testing.T) {
	file, _ := os.OpenFile("test.txt", os.O_RDWR|os.O_CREATE, 0755)
	stdOut := os.Stdout
	os.Stdout = file
	getCmd.Run(getCmd, []string{"Test"})
	file.Seek(0, 0)
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Error("Expected nil but got", err)
	}
	val := strings.Contains(string(content), "no value set")
	fmt.Println(string(content))
	assert.Equalf(t, true, val, "Both should be equal")
	os.Stdout = stdOut
	defer os.Remove("test.txt")
}

func TestGetCmd(t *testing.T) {
	t.Run("TestSetCmd", func(t *testing.T) {
		setCmd.Run(setCmd, []string{"Test-Key1", "Test-Value1"})
	})
	file, _ := os.OpenFile("test.txt", os.O_RDWR|os.O_CREATE, 0755)
	stdOut := os.Stdout
	os.Stdout = file
	getCmd.Run(getCmd, []string{"Test-Key1"})
	file.Seek(0, 0)
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Error("Expected nil but got", err)
	}
	val := strings.Contains(string(content), "Test-Key1 = Test-Value1")
	fmt.Println(string(content))
	assert.Equalf(t, true, val, "Both should be equal")
	os.Stdout = stdOut
	defer os.Remove("test.txt")
}

func TestRootCmd(t *testing.T) {

}

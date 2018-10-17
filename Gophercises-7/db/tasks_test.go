package db

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
	"github.ibm.com/dash/dash_utils/dashtest"
)

var taskid int
var dbName = "test.db"

func TestInit(t *testing.T) {
	home := "demo"
	dbPath := filepath.Join(home, dbName)
	err := Init(dbPath)
	msg := "Init failed"
	if err != nil && err.Error() == "" {
		msg = "Init passed"
	}
	fmt.Println(msg)
}
func TestCreateTask(t *testing.T) {
	task := "Mumbai Marathon 2019"
	home, _ := homedir.Dir()
	dbPath := filepath.Join(home, dbName)
	defer os.Remove(dbPath)
	must(Init(dbPath))
	Close()
	_, err := CreateTask(task)
	msg := "Negative test for create task fails"
	if err != nil {
		msg = fmt.Sprintf("Negative test for create task pass")
	}
	taskid, err = CreateTask(task)
	Close()
	msg = "Failed to add new task"
	if err == nil {
		msg = fmt.Sprintf("Added new task with id : %d", taskid)
	}

	fmt.Println(msg)

}

func TestAllTask(t *testing.T) {
	task := "Mumbai Marathon 2019"
	home, _ := homedir.Dir()
	dbPath := filepath.Join(home, dbName)
	defer os.Remove(dbPath)
	must(Init(dbPath))
	CreateTask(task)
	tasks, err := AllTasks()
	Close()
	msg := "Failed to list tasks"
	if err == nil && len(tasks) >= 0 {
		msg = fmt.Sprintf("Test Passed")
	}

	Close()
	tasks, err = AllTasks()
	msg = "Failed to list tasks"
	if err == nil && len(tasks) >= 0 {
		msg = fmt.Sprintf("Test Passed")
	}

	fmt.Println(msg)
}

func TestDeleteTask(t *testing.T) {
	task := "Mumbai Marathon 2019"
	home, _ := homedir.Dir()
	dbPath := filepath.Join(home, dbName)
	defer os.Remove(dbPath)
	must(Init(dbPath))
	id, _ := CreateTask(task)
	err := DeleteTask(id)
	Close()
	msg := "Failed to delete tasks"
	if err == nil {
		msg = fmt.Sprintf("Test Passed")
	}

	fmt.Println(msg)

}
func must(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func TestMain(m *testing.M) {
	dashtest.ControlCoverage(m)
}

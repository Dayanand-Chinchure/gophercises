package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/Dayanand-Chinchure/gophercises/task/db"
	homedir "github.com/mitchellh/go-homedir"
	"github.ibm.com/dash/dash_utils/dashtest"
)

func TestCreateCommand(t *testing.T) {
	home, _ := homedir.Dir()
	FilePath := filepath.Join(home, "task.txt")
	DbPath := filepath.Join(home, "task.db")
	db.Init(DbPath)
	File, _ := os.OpenFile(FilePath, os.O_RDWR|os.O_CREATE, 0755)
	Sout := os.Stdout
	os.Stdout = File

	addCmd.Run(addCmd, []string{"Dayanand"})
	File.Seek(0, 0)
	content, err := ioutil.ReadAll(File)
	if err != nil {
		t.Error("Expected nil but got", err)
	}
	addCmd.Run(addCmd, []string{"Half-Marathon"})
	File.Seek(0, 0)
	content, err = ioutil.ReadAll(File)
	if err != nil {
		t.Error("Expected nil but got", err)
	}
	t.Log(string(content))
	db.Close()
	addCmd.Run(addCmd, []string{"Full-Marathon"})
	File.Seek(0, 0)
	content, err = ioutil.ReadAll(File)
	if err != nil {
		t.Error("Expected nil but got", err)
	}
	defer os.Remove(FilePath)
	os.Stdout = Sout
	defer os.Remove(DbPath)
}

func TestCompleteCommand(t *testing.T) {
	home, _ := homedir.Dir()
	FilePath := filepath.Join(home, "task.txt")
	DbPath := filepath.Join(home, "task.db")
	db.Init(DbPath)
	File, _ := os.OpenFile(FilePath, os.O_RDWR|os.O_CREATE, 0755)
	Sout := os.Stdout
	os.Stdout = File

	db.CreateTask("Marathon 2019")
	doCmd.Run(doCmd, []string{"1"})
	File.Seek(0, 0)
	content, err := ioutil.ReadAll(File)
	if err != nil {
		t.Error("Expected nil but got", err)
	}
	t.Log(string(content))

	doCmd.Run(doCmd, []string{"10"})
	File.Seek(0, 0)
	content, err = ioutil.ReadAll(File)
	if err != nil {
		t.Error("Expected nil but got", err)
	}

	db.Close()
	doCmd.Run(doCmd, []string{"1"})
	File.Seek(0, 0)
	content, err = ioutil.ReadAll(File)
	if err != nil {
		t.Error("Expected nil but got", err)
	}

	doCmd.Run(doCmd, []string{"Marathon 2019"})
	File.Seek(0, 0)
	content, err = ioutil.ReadAll(File)
	if err != nil {
		t.Error("Expected nil but got", err)
	}

	db.Close()
	doCmd.Run(doCmd, []string{"Marathon 2019"})
	File.Seek(0, 0)
	content, err = ioutil.ReadAll(File)
	if err != nil {
		t.Error("Expected nil but got", err)
	}
	defer os.Remove(FilePath)
	os.Stdout = Sout
	defer os.Remove(DbPath)
}

func TestListCommand(t *testing.T) {
	home, _ := homedir.Dir()
	FilePath := filepath.Join(home, "task.txt")
	DbPath := filepath.Join(home, "task.db")
	db.Init(DbPath)
	File, _ := os.OpenFile(FilePath, os.O_RDWR|os.O_CREATE, 0755)
	Sout := os.Stdout
	os.Stdout = File

	listCmd.Run(listCmd, []string{""})
	File.Seek(0, 0)
	content, err := ioutil.ReadAll(File)
	if err != nil {
		t.Error("Expected nil but got", err)
	}
	t.Log(string(content))

	db.CreateTask("Mumbai Marathon 2018")
	db.CreateTask("Mumbai Marathon 2019")

	listCmd.Run(listCmd, []string{""})
	File.Seek(0, 0)
	content, err = ioutil.ReadAll(File)
	if err != nil {
		t.Error("Expected nil but got", err)
	}

	db.Close()
	listCmd.Run(listCmd, []string{""})
	File.Seek(0, 0)
	content, err = ioutil.ReadAll(File)
	if err != nil {
		t.Error("Expected nil but got", err)
	}

	defer os.Remove(FilePath)
	os.Stdout = Sout
	defer os.Remove(DbPath)
}

func TestMain(m *testing.M) {
	dashtest.ControlCoverage(m)
}

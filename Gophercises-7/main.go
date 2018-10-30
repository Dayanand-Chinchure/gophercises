package main

import (
	"fmt"
	"path/filepath"

	"github.com/Dayanand-Chinchure/gophercises/task/cmd"
	"github.com/Dayanand-Chinchure/gophercises/task/db"
	homedir "github.com/mitchellh/go-homedir"
)

var homeDir = homedir.Dir

//InitApp ...
func main() {
	home, err := homeDir()
	if err != nil {
		fmt.Println("Error reading home directory")
	}
	dbPath := filepath.Join(home, "tasks.db")
	db.Init(dbPath)
	cmd.RootCmd.Execute()
}

package main

import (
	"path/filepath"

	"github.com/Dayanand-Chinchure/gophercises/task/cmd"
	"github.com/Dayanand-Chinchure/gophercises/task/db"
	homedir "github.com/mitchellh/go-homedir"
)

func main() {
	InitApp()
}

//InitApp ...
func InitApp() {
	home, _ := homedir.Dir()
	dbPath := filepath.Join(home, "tasks.db")
	db.Init(dbPath)
	cmd.RootCmd.Execute()
}

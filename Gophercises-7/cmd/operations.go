package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Dayanand-Chinchure/gophercises/task/db"
	"github.com/spf13/cobra"
)

func createCommand(cmd *cobra.Command, args []string) {
	task := strings.Join(args, " ")
	_, err := db.CreateTask(task)
	if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	}
	fmt.Printf("New task added %s.\n", task)
}

func completeCommand(cmd *cobra.Command, args []string) {
	for _, arg := range args {
		id, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Println("Failed to parse the argument:", arg)
		} else {
			err := db.DeleteTask(id)
			if err != nil {
				fmt.Printf("Failed to mark \"%d\" as completed. Error: %s\n", id, err)
			}
			fmt.Printf("Marked \"%d\" as completed.\n", id)
		}
	}
}

func listCommand(cmd *cobra.Command, args []string) {
	tasks, err := db.AllTasks()
	if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	}
	if len(tasks) == 0 {
		fmt.Println("You have no tasks to complete! Why not take a vacation? 🏖")
		return
	}
	fmt.Println("You have the following tasks:")
	for i, task := range tasks {
		fmt.Printf("%d. %s\n", i+1, task.Value)
	}
}

var (
	addCmd = &cobra.Command{
		Use:   "add",
		Short: "Adds a task to your task list.",
		Run:   createCommand,
	}

	doCmd = &cobra.Command{
		Use:   "do",
		Short: "Marks a task as complete",
		Run:   completeCommand,
	}

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists all the tasks",
		Run:   listCommand,
	}
)

func init() {
	RootCmd.AddCommand(addCmd, listCmd, doCmd)
}

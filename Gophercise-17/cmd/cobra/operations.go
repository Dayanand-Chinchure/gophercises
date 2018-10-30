package cobra

import (
	"fmt"
	"path/filepath"

	"github.com/Dayanand-Chinchure/gophercises/secret"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

//RootCmd is the root command to store secret information
var RootCmd = &cobra.Command{
	Use:   "secret",
	Short: "Secret is an API key and other secrets manager",
}

var encodingKey string

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a secret in your secret storage",
	Run: func(cmd *cobra.Command, args []string) {
		v := secret.File(encodingKey, secretsPath())
		key := args[0]
		value, err := v.Get(key)
		if err != nil {
			fmt.Println("no value set")
			return
		}
		fmt.Printf("%s = %s\n", key, value)
	},
}

//Command to save the secret informatio in secret store
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Sets a secret in your secret storage",
	Run: func(cmd *cobra.Command, args []string) {
		v := secret.File(encodingKey, secretsPath())
		var err error
		var key, value string
		//Need minimum two arguments
		if len(args) < 2 {
			fmt.Println("Not enough arguments")
			return
		}
		key, value = args[0], args[1]
		err = v.Set(key, value)
		if err != nil {
			fmt.Println("Unable to add secret")
			return
		}
		fmt.Println("Value set successfully")
	},
}

func init() {
	//Flag used in subsequent commands
	RootCmd.AddCommand(getCmd, setCmd)
	RootCmd.PersistentFlags().StringVarP(&encodingKey, "key", "k", "", "the key to use when encoding and decoding secrets")
}

//Get the path for secret store
func secretsPath() string {
	home, _ := homedir.Dir()
	return filepath.Join(home, ".secrets")
}

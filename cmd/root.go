package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/JorgeEmanoel/money-keeper-backend/database"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mk",
	Short: "Manage Money Keeper backend",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			usage()
			os.Exit(1)
		}

		if args[0] == "help" {
			usage()
			return
		}

		switch args[0] {
		case "serve":
			Serve(database.Connect())
		case "migrate":
			Migrate(database.Connect())
		default:
			log.Fatal("Invalid command")
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Println(`Usage:
	money-keeper-backend <command>

The commands are:
	serve	Start HTTP server
	help	Displays this message`)
}

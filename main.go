package main

import (
	"fmt"

	"github.com/harnyk/gotp/internal/application"
	"github.com/harnyk/gotp/internal/storage"
	"github.com/spf13/cobra"
)

func main() {
	// Usage:
	// gotp - list secrets
	// gotp list - list secrets
	// gotp show <key> - generate code
	// gotp add <key> - add secret
	// gotp delete <key> - delete secret

	repo := storage.NewSecretsRepository()

	app := application.NewApp(repo)

	var rootCmd = &cobra.Command{
		Use:   "gotp",
		Short: "A simple tool to generate and store TOTP codes",
		Long:  `A simple tool to generate and store TOTP codes`,
	}

	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "List all secrets",
			Long:  `List all secrets`,
			Run: func(cmd *cobra.Command, args []string) {
				app.CmdList()
			},
		},
		&cobra.Command{
			Use:   "show <key>",
			Short: "Generate code for a secret",
			Long:  `Generate code for a secret`,
			Args:  cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				app.CmdGenerate(args[0])
			},
		},
		&cobra.Command{
			Use:   "add <key>",
			Short: "Add a secret",
			Long:  `Add a secret`,
			Args:  cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				app.CmdAdd(args[0])
			},
		},
		&cobra.Command{
			Use:   "delete <key>",
			Short: "Delete a secret",
			Long:  `Delete a secret`,
			Args:  cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("Deleting secret")
			},
		},
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		return
	}

}

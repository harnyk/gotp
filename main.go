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
		Long:  "Gotp stores your TOTP codes in a local file encrypted with a master password.",
	}

	rootCmd.AddCommand(
		&cobra.Command{
			Use:     "list",
			Short:   "List all secrets",
			Long:    "List all secrets. Select a secret to generate a code",
			Aliases: []string{"ls", "l"},
			Run: func(cmd *cobra.Command, args []string) {
				app.CmdList()
			},
		},
		&cobra.Command{
			Use:     "show <key>",
			Short:   "Generate code for a secret",
			Long:    "Generate a OTP code for a secret",
			Aliases: []string{"generate", "gen", "s"},
			Args:    cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				app.CmdGenerate(args[0])
			},
		},
		&cobra.Command{
			Use:     "add <key>",
			Short:   "Add a secret",
			Long:    "Add a secret. You will be prompted to enter the secret",
			Aliases: []string{"a", "new", "create"},
			Args:    cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				app.CmdAdd(args[0])
			},
		},
		&cobra.Command{
			Use:     "delete <key>",
			Short:   "Delete a secret",
			Long:    "Delete a secret by key",
			Aliases: []string{"rm", "remove", "del", "d"},
			Args:    cobra.ExactArgs(1),
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

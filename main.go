package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/hgfischer/go-otp"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

func getTickingChannel(ctx context.Context, totp *otp.TOTP) chan string {
	ch := make(chan string)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				ch <- totp.Now().Get()
				time.Sleep(time.Duration(1) * time.Second)
			}
		}
	}()
	return ch
}

func main() {
	// Usage:
	// gotp - list secrets
	// gotp list - list secrets
	// gotp show <key> - generate code
	// gotp add <key> - add secret
	// gotp delete <key> - delete secret

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
				fmt.Println("Listing secrets")
			},
		},
		&cobra.Command{
			Use:   "show <key>",
			Short: "Generate code for a secret",
			Long:  `Generate code for a secret`,
			Args:  cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("Generating code for secret")
			},
		},
		&cobra.Command{
			Use:   "add <key>",
			Short: "Add a secret",
			Long:  `Add a secret`,
			Args:  cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Print("Enter secret: ")
				secret, err := terminal.ReadPassword(int(os.Stdin.Fd()))
				if err != nil {
					panic(err)
				}
				fmt.Println()
				fmt.Printf("Adding secret '%s' = '%s'", args[0], secret)
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

func example() {
	totp := &otp.TOTP{
		Secret:         "sample-secret",
		IsBase32Secret: true,
	}

	ctx, cancel := context.WithCancel(context.Background())

	ch := getTickingChannel(ctx, totp)

	keyEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
	defer cancel()
	defer keyboard.Close()

	for {
		select {
		case ev := <-keyEvents:
			if ev.Key == keyboard.KeyEsc ||
				ev.Key == keyboard.KeyCtrlC ||
				ev.Rune == 'q' {
				fmt.Println("Canceled")
				cancel()
				return
			}
		case code := <-ch:
			fmt.Println(code)
		}
	}
}

package application

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/harnyk/gotp/internal/storage"
	"github.com/hgfischer/go-otp"
	"github.com/manifoldco/promptui"
	"golang.org/x/crypto/ssh/terminal"
)

type App struct {
	store storage.ISecretsRepository
}

func NewApp(store storage.ISecretsRepository) *App {
	return &App{
		store: store,
	}
}

func (a *App) CmdList() {
	secrets, err := a.store.ListKeys()
	if err != nil {
		panic(err)
	}

	prompt := promptui.Select{
		Label: "Select a secret",
		Items: secrets,
	}
	_, result, err := prompt.Run()
	if err != nil {
		panic(err)
	}
	a.CmdGenerate(result)
}

func (a *App) CmdAdd(key string) {
	fmt.Print("Enter secret: ")
	secret, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Printf("Adding secret '%s' = '%s'", key, secret)
	err = a.store.SetSecret(key, string(secret))
	if err != nil {
		panic(err)
	}
}

func (a *App) CmdDelete(key string) {
	err := a.store.DeleteSecret(key)
	if err != nil {
		panic(err)
	}
}

func (a *App) CmdGenerate(key string) {
	secret, err := a.store.GetSecret(key)
	if err != nil {
		panic(err)
	}
	showCode(secret)
}

//------------------------------------------------------------------------------

type codeWithTime struct {
	code string
	time time.Duration
}

func getTickingChannel(ctx context.Context, totp *otp.TOTP) chan codeWithTime {
	ch := make(chan codeWithTime)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				{
					nowTime := time.Now()
					totp.Time = nowTime
					ts := uint64(nowTime.Unix() / int64(totp.Period))
					timeOfPeriodStart := time.Unix(int64(ts)*int64(totp.Period), 0)
					timeOfPeriodEnd := timeOfPeriodStart.Add(
						time.Duration(
							uint64(totp.Period) * uint64(time.Second)))
					timeToNextPeriod := timeOfPeriodEnd.Sub(nowTime) / time.Second
					ch <- codeWithTime{totp.Get(), timeToNextPeriod}
					time.Sleep(time.Duration(1) * time.Second)
				}
			}
		}
	}()
	return ch
}

func showCode(secret string) {
	totp := &otp.TOTP{
		Secret:         secret,
		IsBase32Secret: true,
		Period:         30,
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
			fmt.Print(code)
			fmt.Print("\r")
		}
	}
}

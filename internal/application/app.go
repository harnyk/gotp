package application

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/harnyk/gotp/internal/storage"
	"github.com/hgfischer/go-otp"
	"github.com/manifoldco/promptui"
	"golang.org/x/term"
)

type App struct {
	store         storage.ISecretsRepository
	isInteractive bool
}

func NewApp(store storage.ISecretsRepository) *App {
	return &App{
		store:         store,
		isInteractive: term.IsTerminal(int(os.Stdout.Fd())) && term.IsTerminal(int(os.Stdin.Fd())),
	}
}

func (a *App) CmdList() {
	secrets, err := a.store.ListKeys()
	if err != nil {
		panic(err)
	}

	if !a.isInteractive {
		for _, secret := range secrets {
			fmt.Println(secret)
		}
		return
	}

	prompt := promptui.Select{
		Label: "Select a secret",
		Items: secrets,
	}
	_, result, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrInterrupt {
			fmt.Println("Interrupted")
			return
		}
		panic(err)
	}
	a.CmdGenerate(result)
}

func (a *App) CmdAdd(key string) {
	var secret []byte
	var err error
	if !a.isInteractive {
		secret, err = bufio.NewReader(os.Stdin).ReadBytes('\n')
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Print("Enter secret: ")
		secret, err = term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			panic(err)
		}
		fmt.Printf("Adding secret '%s' = '%s'\n", key, secret)
	}

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

	fmt.Printf("Deleted secret '%s'", key)
}

func (a *App) CmdGenerate(key string) {
	secret, err := a.store.GetSecret(key)
	if err != nil {
		panic(err)
	}

	if !a.isInteractive {
		code := generateCode(secret)
		fmt.Println(code.code)
		return
	}

	showCode(secret)
}

//------------------------------------------------------------------------------

type codeWithTime struct {
	code string
	time uint8
}

func generateCode(secret string) codeWithTime {
	totp := &otp.TOTP{
		Secret:         secret,
		IsBase32Secret: true,
		Period:         30,
	}

	nowTime := time.Now()
	totp.Time = nowTime
	ts := uint64(nowTime.Unix() / int64(totp.Period))
	timeOfPeriodStart := time.Unix(int64(ts)*int64(totp.Period), 0)
	timeOfPeriodEnd := timeOfPeriodStart.Add(
		time.Duration(
			uint64(totp.Period) * uint64(time.Second)))
	timeToNextPeriod := timeOfPeriodEnd.Sub(nowTime) / time.Second
	return codeWithTime{totp.Get(), uint8(timeToNextPeriod)}
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
					ch <- generateCode(totp.Secret)
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
				fmt.Println()
				cancel()
				return
			}
		case code := <-ch:
			fmt.Print("\r")
			fmt.Printf(
				"%s %s    ",
				promptui.Styler(promptui.FGBold)(code.code),
				promptui.Styler(promptui.FGFaint)(
					fmt.Sprintf("(%ds)", code.time),
				),
			)
		}
	}
}

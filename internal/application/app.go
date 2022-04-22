package application

import "github.com/harnyk/gotp/internal/storage"

type App struct {
	store storage.ISecretsRepository
}

func NewApp(store storage.ISecretsRepository) *App {
	return &App{
		store: store,
	}
}

func (a *App) CmdList() {
}

func (a *App) CmdAdd(key string) {

}

func (a *App) CmdDelete(key string) {
}

func (a *App) CmdGenerate() {
}

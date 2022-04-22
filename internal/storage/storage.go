package storage

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/kirsle/configdir"
	"github.com/pelletier/go-toml"
)

type ISecretsRepository interface {
	GetSecret(id string) (string, error)
	SetSecret(id string, secret string) error
	DeleteSecret(id string) error
	ListSecrets() ([]string, error)
}

type SecretsRepository struct {
	filePath string
}

func NewSecretsRepository() *SecretsRepository {
	dir := configdir.LocalConfig("gotp")
	filePath := path.Join(dir, "secrets.toml")

	return &SecretsRepository{
		filePath: filePath,
	}
}

func (s *SecretsRepository) loadSecrets() (tree *toml.Tree, err error) {
	secrets, err := toml.LoadFile(s.filePath)
	return secrets, err
}

func (s *SecretsRepository) saveSecrets(secrets *toml.Tree) error {
	tomlStr := secrets.String()

	//make sure the directory exists
	err := os.MkdirAll(path.Dir(s.filePath), 0755)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(s.filePath, []byte(tomlStr), 0644)
}

func (s *SecretsRepository) GetSecret(id string) (string, error) {
	secrets, err := s.loadSecrets()
	if err != nil {
		return "", err
	}

	secret, ok := secrets.Get(id).(string)
	if !ok {
		return "", nil
	}

	return secret, nil
}

func (s *SecretsRepository) SetSecret(id string, secret string) error {
	secrets, err := s.loadSecrets()
	if err != nil {
		return err
	}
	secrets.Set(id, secret)
	return s.saveSecrets(secrets)
}

func (s *SecretsRepository) DeleteSecret(id string) error {
	secrets, err := s.loadSecrets()
	if err != nil {
		return err
	}
	secrets.Delete(id)
	return s.saveSecrets(secrets)
}

func (s *SecretsRepository) ListSecrets() ([]string, error) {
	secrets, err := s.loadSecrets()
	if err != nil {
		return nil, err
	}
	return secrets.Keys(), nil
}

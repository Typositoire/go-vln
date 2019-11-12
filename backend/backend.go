package backend

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Backend ""
type Backend interface {
	FindTarget(string) (string, error)
	BackendIsInit() (bool, error)
	Auth() error
}

type backend struct {
	logger *log.Entry
}

// NewBackend ""
func NewBackend(backend string) (Backend, error) {
	var (
		be  Backend
		err error
	)

	switch backend {
	case "vault":
		be, err = newVaultBackend(viper.GetString("be-vault-addr"), viper.GetString("be-vault-symlinkdbpath"))
	case "file":
		be, err = newFileBackend(viper.GetString("be-file-path"))
	case "git":
		be, err = newGitBackend(viper.GetString("be-git-repository"), viper.GetString("be-git-accesstoken"))
	default:
		be = nil
		err = errors.New("invalid backend " + backend)
	}

	return be, err
}

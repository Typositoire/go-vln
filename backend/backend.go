package backend

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Backend ""
type Backend interface {
	FindTarget(string) (string, error)
	BackendIsInit() (bool, error)
	BackendCanProcess(*http.Request) bool
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
		be, err = newVaultBackend(viper.GetString("be-vault-addr"), viper.GetString("symlinkdb-path"))
	case "file":
		be, err = newfileBackend(viper.GetString("be-file-path"))
	}

	return be, err
}

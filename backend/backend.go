package backend

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Backend ""
type Backend interface {
	FindTarget(string) (string, error)
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
	}

	return be, err
}

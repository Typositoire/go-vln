package backend

import (
	"errors"

	log "github.com/sirupsen/logrus"
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
func NewBackend(options map[string]string) (Backend, error) {
	var (
		be  Backend
		err error
	)

	switch options["backend"] {
	case "vault":
		be, err = newVaultBackend(options["beVaultAddr"], options["beVaultSymlinkDBPath"])
	case "file":
		be, err = newFileBackend(options["beFilePath"])
	case "git":
		be, err = newGitBackend(options["beGitRepository"], options["beGitAccessToken"])
	case "":
		be = nil
		err = errors.New("empty backend")
	default:
		be = nil
		err = errors.New("invalid backend " + options["backend"])
	}

	return be, err
}

package backend

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"

	log "github.com/sirupsen/logrus"
)

type vaultBackend struct {
	VaultURL      string
	VaultToken    string
	SymlinkDBPath string
	HTTPClient    *resty.Client
	Logger        *log.Entry
}

func newVaultBackend(vaultURL, symlinkDBPath string) (Backend, error) {
	logger := log.WithFields(log.Fields{
		"component": "backend.vault",
	})

	log.SetFormatter(&log.JSONFormatter{})

	client := resty.New()

	client.
		SetHostURL(vaultURL + "/v1").
		SetLogger(logger)

	return &vaultBackend{
		VaultURL:      vaultURL,
		VaultToken:    "myroot",
		SymlinkDBPath: symlinkDBPath,
		HTTPClient:    client,
		Logger:        logger,
	}, nil
}

func (b vaultBackend) FindTarget(path string) (string, error) {
	var (
		body map[string]interface{}
	)

	resp, err := b.HTTPClient.R().SetHeader("X-Vault-Token", b.VaultToken).Get(b.SymlinkDBPath)

	if err != nil {
		b.Logger.WithError(err)
		return "", err
	}

	err = json.Unmarshal(resp.Body(), &body)
	if err != nil {
		fmt.Println("Body " + err.Error())
		return "", err
	}

	for k, v := range body["data"].(map[string]interface{})["data"].(map[string]interface{}) {
		if path == k {
			return v.(string), nil
		}
	}

	return path, nil
}

func (b vaultBackend) getSymlinkDatabase(path string) (map[string]string, error) {
	return nil, nil
}

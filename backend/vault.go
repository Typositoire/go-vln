package backend

import (
	"errors"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"

	vaultapi "github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
)

type vaultBackend struct {
	VaultClient   *vaultapi.Client
	SymlinkDBPath string
	HTTPClient    *resty.Client
	Logger        *log.Entry
}

func newVaultBackend(vaultURL, symlinkDBPath string) (Backend, error) {
	logger := log.WithFields(log.Fields{
		"component": "backend.vault",
	})

	log.SetFormatter(&log.JSONFormatter{})

	config := &vaultapi.Config{
		Address: vaultURL,
	}

	vClient, err := vaultapi.NewClient(config)

	if err != nil {
		return nil, err
	}

	client := resty.New()

	client.
		SetHostURL(vaultURL + "/v1").
		SetLogger(logger)

	return &vaultBackend{
		VaultClient:   vClient,
		SymlinkDBPath: symlinkDBPath,
		HTTPClient:    client,
		Logger:        logger,
	}, nil
}

func (b *vaultBackend) Auth() error {
	resp, err := b.VaultClient.Logical().Write("auth/approle/login", map[string]interface{}{
		"role_id":   viper.GetString("vault-app-role-id"),
		"secret_id": viper.GetString("vault-app-role-secret"),
	})

	if err != nil {
		return err
	}

	if resp.Auth == nil {
		return errors.New("no auth info returned")
	}

	b.VaultClient.SetToken(resp.Auth.ClientToken)

	return nil
}

func (b vaultBackend) BackendIsInit() (bool, error) {
	_, err := b.VaultClient.Logical().Read(b.SymlinkDBPath)

	if err != nil {
		b.Logger.Errorln(err)
		return false, err
	}

	return true, nil
}

func (b vaultBackend) FindTarget(path string) (string, error) {
	resp, err := b.VaultClient.Logical().Read(b.SymlinkDBPath)

	if err != nil {
		b.Logger.Errorln("DEBUG1" + err.Error())
		return "", err
	}

	for k, v := range resp.Data["data"].(map[string]interface{}) {
		if path == k {
			return v.(string), nil
		}
	}

	return path, nil
}

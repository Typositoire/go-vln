package backend

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"

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
		SymlinkDBPath: symlinkDBPath,
		HTTPClient:    client,
		Logger:        logger,
	}, nil
}

func (b *vaultBackend) Auth() error {
	resp, err := b.HTTPClient.R().
		SetBody(map[string]string{
			"role_id":   viper.GetString("vault-app-role-id"),
			"secret_id": viper.GetString("vault-app-role-secret"),
		}).
		Post("/auth/approle/login")

	if err != nil {
		return err
	}

	var body map[string]interface{}

	fmt.Println(string(resp.Body()))

	json.Unmarshal(resp.Body(), &body)

	b.VaultToken = body["auth"].(map[string]interface{})["client_token"].(string)

	return nil
}

func (b vaultBackend) BackendIsInit() (bool, error) {
	resp, err := b.HTTPClient.R().SetHeader("X-Vault-Token", b.VaultToken).Get(b.SymlinkDBPath)

	if err != nil {
		b.Logger.Errorln(err)
		return false, err
	}

	if resp.StatusCode() == 404 {
		return false, nil
	}

	return true, nil
}

func (b vaultBackend) BackendCanProcess(r *http.Request) bool {
	if r.Method != "GET" {
		return false
	}

	if strings.HasPrefix(r.RequestURI, "/v1/sys") {
		return false
	}

	return true
}

func (b vaultBackend) FindTarget(path string) (string, error) {
	var (
		body map[string]interface{}
	)

	resp, err := b.HTTPClient.R().SetHeader("X-Vault-Token", b.VaultToken).Get(b.SymlinkDBPath)

	if err != nil {
		b.Logger.Errorln(err)
		return "", err
	}

	err = json.Unmarshal(resp.Body(), &body)
	if err != nil {
		b.Logger.Errorln(err)
		return "", err
	}

	if resp.StatusCode() > 299 {
		fmt.Println(string(resp.Body()))
		return "", errors.New("Error")
	}
	for k, v := range body["data"].(map[string]interface{})["data"].(map[string]interface{}) {
		if path == k {
			return v.(string), nil
		}
	}

	return path, nil
}

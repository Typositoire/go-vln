package backend

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

type fileBackend struct {
	FilePath string
	Logger   *log.Entry
}

func newFileBackend(filePath string) (Backend, error) {
	logger := log.WithFields(log.Fields{
		"component": "backend.file",
	})

	log.SetFormatter(&log.JSONFormatter{})
	return &fileBackend{
		FilePath: filePath,
		Logger:   logger,
	}, nil
}

func (b *fileBackend) Auth() error {
	f, err := os.Open(b.FilePath)

	if err != nil {
		return err
	}

	f.Close()

	return nil
}

func (b fileBackend) BackendIsInit() (bool, error) {
	return true, nil
}

func (b fileBackend) FindTarget(path string) (string, error) {
	var db map[string]string

	data, err := ioutil.ReadFile(b.FilePath)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(data, &db)

	if err != nil {
		return "", err
	}

	for k, v := range db {
		if path == k {
			return v, nil
		}
	}

	return path, nil
}

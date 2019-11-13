package backend

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

type gitBackend struct {
	GitOptions *git.CloneOptions
	Logger     *log.Entry
}

func newGitBackend(repo, accessToken string) (Backend, error) {
	logger := log.WithFields(log.Fields{
		"component": "backend.git",
	})

	log.SetFormatter(&log.JSONFormatter{})

	return &gitBackend{
		GitOptions: &git.CloneOptions{
			Auth: &http.BasicAuth{
				Username: "vln",
				Password: accessToken,
			},
			URL:          repo,
			SingleBranch: true,
		},
		Logger: logger,
	}, nil
}

func (b *gitBackend) Auth() error {
	return nil
}

func (b gitBackend) BackendIsInit() (bool, error) {
	_, err := git.Clone(memory.NewStorage(), nil, b.GitOptions)

	if err == transport.ErrRepositoryNotFound {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (b gitBackend) FindTarget(path string) (string, error) {
	var (
		_db string
		db  map[string]string
	)

	r, err := git.Clone(memory.NewStorage(), nil, b.GitOptions)

	if err != nil {
		return "", err
	}

	ref, err := r.Head()
	if err != nil {
		return "", err
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return "", err
	}

	tree, err := commit.Tree()
	if err != nil {
		return "", err
	}

	tree.Files().ForEach(func(f *object.File) error {
		if f.Name == "db.json" {
			_db, err = f.Contents()
		}

		return nil
	})

	if _db != "" {
		err = json.Unmarshal([]byte(_db), &db)

		if err != nil {
			return "", err
		}

		for k, v := range db {
			if path == k {
				return v, nil
			}
		}
	}

	return path, nil
}

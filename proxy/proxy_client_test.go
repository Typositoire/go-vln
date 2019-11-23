package proxy

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestEmptyBackend(t *testing.T) {
	options := make(map[string]string)
	options["backend"] = ""

	_, err := NewProxyClient(options)

	if err == nil {
		t.Fatalf("Expected error, got %s", err.Error())
	}

	if err.Error() != "empty backend" {
		t.Fatalf("Expected \"emtpy backend\", got %s", err.Error())
	}
}

func TestInvalidBackend(t *testing.T) {
	options := make(map[string]string)
	options["backend"] = "meuh"

	_, err := NewProxyClient(options)

	if err == nil {
		t.Fatalf("Expected error, got %s", err.Error())
	}

	if err.Error() != "invalid backend meuh" {
		t.Fatalf("Expected \"invalid backend meuh\", got %s", err.Error())
	}
}

func TestFileBackend(t *testing.T) {
	// Prepare empty file backend
	content := "{}"

	f, err := ioutil.TempFile("/tmp", "db_")

	if err != nil {
		t.Fatalf("Unexpected error, got %s", err.Error())
	}

	_, err = f.Write([]byte(content))

	if err != nil {
		t.Fatalf("Unexpected error, got %s", err.Error())
	}

	defer f.Close()
	defer os.Remove(f.Name())

	options := make(map[string]string)
	options["backend"] = "file"
	options["beFilePath"] = f.Name()

	_, err = NewProxyClient(options)

	if err != nil {
		t.Fatalf("Unexpected error, got %s", err.Error())
	}
}

func TestFileBackendAuthFailure(t *testing.T) {
	options := make(map[string]string)
	options["backend"] = "file"
	options["beFilePath"] = "/tmp/db_1834098209834"

	_, err := NewProxyClient(options)

	if err == nil {
		t.Fatalf("Expected error, got %s", err.Error())
	}
}

// func TestReadSecret(t *testing.T) {
// 	e := echo.New()
// 	req := httptest.NewRequest(echo.GET, "http://localhost/v1/sys/health", nil)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)
// 	c.SetPath("/")

// 	// Prepare empty file backend
// 	content := "{}"

// 	f, err := ioutil.TempFile("/tmp", "db_")

// 	if err != nil {
// 		t.Fatalf("Unexpected error, got %s", err.Error())
// 	}

// 	_, err = f.Write([]byte(content))

// 	if err != nil {
// 		t.Fatalf("Unexpected error, got %s", err.Error())
// 	}

// 	defer f.Close()
// 	defer os.Remove(f.Name())

// 	options := make(map[string]string)
// 	options["backend"] = "file"
// 	options["beFilePath"] = f.Name()

// 	be, err := backend.NewBackend(options)

// 	if err != nil {
// 		t.Fatalf("Unexpected error, got %s", err.Error())
// 	}

// 	log := logrus.New()
// 	log.Out = ioutil.Discard

// 	pc := &pClient{
// 		backend:    be,
// 		httpClient: nil,
// 		logger:     log.WithTime(time.Now()),
// 	}

// 	if assert.NoError(t, pc.processRequest(c)) {
// 		assert.Equal(t, http.StatusOK, rec.Code)
// 	}
// }

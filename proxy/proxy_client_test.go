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

// func TestRun(t *testing.T) {
// 	// Mock
// 	mockCtrl := gomock.NewController(t)
// 	defer mockCtrl.Finish()

// 	mockPC := mocks.NewMockPClient(mockCtrl)

// 	mockPC.EXPECT().Run().Times(1)

// 	mockPC.Run()
// }

func Test
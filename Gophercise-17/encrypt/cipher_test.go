package encrypt

import (
	"os"
	"testing"

	"github.ibm.com/dash/dash_utils/dashtest"
)

func TestMain(m *testing.M) {
	dashtest.ControlCoverage(m)
}

func TestEncryptWriter(t *testing.T) {
	key := "Secret Key"
	file, err := os.OpenFile("test.txt", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Error("Fail to get io Writer")
	}
	stream, err := EncWriter(key, file)
	if err != nil || stream == nil {
		t.Error("Fail to get EncryptWriter")
	}
	defer file.Close()
}

func TestDecryptReader(t *testing.T) {
	key := "Secret Key"
	file, err := os.OpenFile("test.txt", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Error("Fail to get io Reader")
	}
	stream, err := DecryptReader(key, file)
	if err != nil || stream == nil {
		t.Error("Fail to get EncryptReader")
	}
	defer file.Close()
}

func TestDecryptReaderNegative(t *testing.T) {
	key := "Secret Key"
	file, err := os.OpenFile("test1.txt", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Error("Fail to get io Reader")
	}
	stream, err := DecryptReader(key, file)
	if err == nil || stream != nil {
		t.Error("Fail to get EncryptReader")
	}
	defer os.Remove("test1.txt")
}

package secret

import (
	"os"
	"path/filepath"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"github.ibm.com/dash/dash_utils/dashtest"
)

func TestMain(m *testing.M) {
	dashtest.ControlCoverage(m)
}

func TestFile(t *testing.T) {
	key := "TestKey"
	path := "TestPath"
	vault := File(key, path)
	assert.Equalf(t, true, key == vault.encodingKey && path == vault.filePath, "Both should be equal")

}

func TestLoadInValidPath(t *testing.T) {
	key := "TestKey"
	path := "TestPath"
	vault := File(key, path)
	err := vault.load()
	if err != nil {
		t.Errorf("Failed to load with error %v", err.Error())
	}
	assert.Equalf(t, true, key == vault.encodingKey && path == vault.filePath, "Both should be equal")
	assert.Equalf(t, true, len(vault.keyValues) == 0, "Both should be equal")
}

func TestLoad(t *testing.T) {
	home, _ := homedir.Dir()
	path := filepath.Join(home, ".testsecrets")
	key := "TestKey"
	vault := File(key, path)
	vault.keyValues = map[string]string{"test-key": "test-value"}
	err := vault.save()
	if err != nil {
		t.Errorf(err.Error())
	}
	err = vault.load()
	if err != nil {
		t.Errorf("Failed to load with error %v", err.Error())
	}
	defer os.Remove(path)
}

func TestSaveFileOpenFail(t *testing.T) {
	path := ""
	key := "TestKey"
	vault := File(key, path)
	vault.keyValues = map[string]string{"test-key": "test-value"}
	err := vault.save()
	if err == nil {
		t.Errorf(err.Error())
	}
	defer os.Remove("testsecrets")

}

func TestSaveWriterFail(t *testing.T) {
	path := ""
	key := "TestKey"
	vault := File(key, path)
	vault.keyValues = map[string]string{"test-key": "test-value"}
	err := vault.save()
	if err == nil {
		t.Errorf(err.Error())
	}
	defer os.Remove("testsecrets")

}

func TestSave(t *testing.T) {
	home, _ := homedir.Dir()
	path := filepath.Join(home, ".testsecrets")
	key := "TestKey"
	vault := File(key, path)
	vault.keyValues = map[string]string{"test-key": "test-value"}
	err := vault.save()
	if err != nil {
		t.Errorf(err.Error())
	}
	vaultnew := File(key, path)
	vaultnew.load()
	assert.Equalf(t, true, len(vaultnew.keyValues) > 0, "Both should be equal")
	defer os.Remove(path)

}

func TestSetEmptyKeyValue(t *testing.T) {
	home, _ := homedir.Dir()
	path := filepath.Join(home, ".testsecrets")
	secretkey := "TestKey"
	vault := File(secretkey, path)
	err := vault.Set("", "")
	if err == nil {
		t.Errorf("Test fails")
	}

}

func TestSetLoadFails(t *testing.T) {
	home, _ := homedir.Dir()
	path := filepath.Join(home, ".testsecrets")
	os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	secretkey := "TestKey"
	vault := File(secretkey, path)
	err := vault.Set("testkey", "testvalue")
	if err == nil {
		t.Errorf("Test fails")
	}

	os.Remove(path)
}

func TestSetSaveFails(t *testing.T) {
	secretkey := "TestKey"
	vault := File(secretkey, "")
	err := vault.Set("testkey", "testvalue")
	if err == nil {
		t.Errorf("Test fails")
	}
}

func TestGetLoadFails(t *testing.T) {
	home, _ := homedir.Dir()
	path := filepath.Join(home, ".testsecrets")
	os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	secretkey := "TestKey"
	vault := File(secretkey, path)
	_, err := vault.Get("testkey")
	if err == nil {
		t.Errorf("Test fails")
	}
	os.Remove(path)
}

func TestGetNoValue(t *testing.T) {
	home, _ := homedir.Dir()
	path := filepath.Join(home, ".testsecrets")
	key := "TestKey"
	vault := File(key, path)
	vault.save()
	_, err := vault.Get("testkey")
	if err == nil {
		t.Errorf("Test fails")
	}
	os.Remove(path)

}

func TestGet(t *testing.T) {
	home, _ := homedir.Dir()
	path := filepath.Join(home, ".testsecrets")
	key := "TestKey"
	vault := File(key, path)
	vault.keyValues = map[string]string{"test-key": "test-value"}
	vault.save()
	retVal, _ := vault.Get("test-key")
	assert.Equalf(t, true, "test-value" == retVal, "Both should be equal")
	os.Remove(path)

}

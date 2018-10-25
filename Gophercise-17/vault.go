package secret

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"

	"github.com/Dayanand-Chinchure/gophercises/secret/encrypt"
)

//Vault structure to store in secret file
type Vault struct {
	encodingKey string
	filePath    string
	mutex       sync.Mutex
	keyValues   map[string]string
}

//File to initialises the Vault
func File(encodingKey, filePath string) Vault {
	return Vault{
		encodingKey: encodingKey,
		filePath:    filePath,
		keyValues:   make(map[string]string),
	}
}

//load to retrive encrypted info into Vault
func (v *Vault) load() error {
	f, err := os.Open(v.filePath)
	if err != nil {
		v.keyValues = make(map[string]string)
		return nil
	}
	defer f.Close()
	//Get reader to decrypt the content
	r, err := encrypt.DecryptReader(v.encodingKey, f)
	if err != nil {
		return err
	}
	return v.readKeyValues(r)

}

//save encrypted data into the file
func (v *Vault) save() error {
	f, err := os.OpenFile(v.filePath, os.O_RDWR|os.O_CREATE, 0775)
	if err != nil {
		return err
	}
	//Get writer to write encrypted information
	w, err := encrypt.EncWriter(v.encodingKey, f)
	if err != nil {
		return err
	}
	return v.writeKeyValues(w)
}

//readKeyValues to decode
func (v *Vault) readKeyValues(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(&v.keyValues)
}

//writeKeyValues to encode
func (v *Vault) writeKeyValues(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(v.keyValues)
}

//Set value for key
func (v *Vault) Set(key, value string) error {
	var err error
	if key == "" && value == "" {
		err = errors.New("Insufficient params")
		return err
	}
	v.mutex.Lock()
	defer v.mutex.Unlock()
	err = v.load()
	if err != nil {
		return err
	}
	v.keyValues[key] = value
	return v.save()
}

//Get value from the key
func (v *Vault) Get(key string) (string, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	err := v.load()
	if err != nil {
		return "", err
	}
	value, ok := v.keyValues[key]
	if !ok {
		return "", errors.New("secret: no value for that key")
	}
	return value, nil
}

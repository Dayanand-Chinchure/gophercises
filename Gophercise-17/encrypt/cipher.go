package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

//encryptStream return the encrypt stream
func encryptStream(key string, iv []byte) (cipher.Stream, error) {
	block, err := newCipherBlock(key)
	return cipher.NewCFBEncrypter(block, iv), err
}

//EncWriter returns the stream writter which encrypt the incoming content
func EncWriter(key string, w io.Writer) (*cipher.StreamWriter, error) {
	iv := make([]byte, aes.BlockSize)
	_, err := io.ReadFull(rand.Reader, iv)
	stream, err := encryptStream(key, iv)
	_, err = w.Write(iv)
	return &cipher.StreamWriter{S: stream, W: w}, err
}

//decryptStream return the decrypt stream
func decryptStream(key string, iv []byte) (cipher.Stream, error) {
	block, err := newCipherBlock(key)
	return cipher.NewCFBDecrypter(block, iv), err
}

//DecryptReader returns the stream reader which decrypt the content
func DecryptReader(key string, r io.Reader) (*cipher.StreamReader, error) {
	iv := make([]byte, aes.BlockSize)
	n, err := r.Read(iv)
	if err != nil || n != len(iv) {
		return nil, errors.New("Encrypt: Unable to read IV")
	}
	stream, err := decryptStream(key, iv)
	return &cipher.StreamReader{S: stream, R: r}, err
}

//newCipherBlock create fixed length key and return cipher block
func newCipherBlock(key string) (cipher.Block, error) {
	hasher := md5.New()
	fmt.Fprint(hasher, key)
	cipherKey := hasher.Sum(nil)
	return aes.NewCipher(cipherKey)
}

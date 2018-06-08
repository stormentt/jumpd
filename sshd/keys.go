package sshd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

func GetHostKey() (ssh.Signer, error) {
	hostKeyPath := viper.GetString("config.hostkey.path")

	hostkey, err := loadHostKey(hostKeyPath)
	if err != nil {
		if os.IsNotExist(err) {
			hostkey, err = generateHostKey(hostKeyPath)
		}
	}

	return hostkey, err
}

func loadHostKey(keyPath string) (ssh.Signer, error) {
	ctxLog := log.WithFields(log.Fields{
		"path": keyPath,
	})
	ctxLog.Info("loading hostkey")

	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	return ssh.ParsePrivateKey(keyBytes)
}

func generateHostKey(keyPath string) (ssh.Signer, error) {
	keySize := viper.GetInt("config.hostkey.size")

	ctxLog := log.WithFields(log.Fields{
		"size": keySize,
		"path": keyPath,
	})

	ctxLog.Info("generating RSA key")

	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, err
	}

	ctxLog.Debug("key generated, saving")

	keyFile, err := os.Create(keyPath)
	if err != nil {
		return nil, err
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	PEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateKeyBytes}

	err = pem.Encode(keyFile, PEM)
	if err != nil {
		return nil, err
	}

	err = keyFile.Close()
	if err != nil {
		return nil, err
	}

	ctxLog.Info("key saved")
	return ssh.NewSignerFromKey(privateKey)
}

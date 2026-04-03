package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
)

const (
	certificatesDir = "certificates"
	privateKeyBits  = 2048
)

func EnsureKeyPair(service string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privatePath, publicPath := keyPaths(service)

	privateKey, err := loadPrivateKeyFromFile(privatePath)
	if err == nil {
		if _, err := os.Stat(publicPath); os.IsNotExist(err) {
			if err := savePublicKey(publicPath, &privateKey.PublicKey); err != nil {
				return nil, nil, err
			}
		}

		return privateKey, &privateKey.PublicKey, nil
	}

	if !os.IsNotExist(err) {
		return nil, nil, err
	}

	if err := os.MkdirAll(certificatesDir, 0o755); err != nil {
		return nil, nil, fmt.Errorf("erro ao criar pasta de certificados: %w", err)
	}

	privateKey, err = rsa.GenerateKey(rand.Reader, privateKeyBits)
	if err != nil {
		return nil, nil, fmt.Errorf("erro ao gerar chave privada: %w", err)
	}

	if err := savePrivateKey(privatePath, privateKey); err != nil {
		return nil, nil, err
	}

	if err := savePublicKey(publicPath, &privateKey.PublicKey); err != nil {
		return nil, nil, err
	}

	return privateKey, &privateKey.PublicKey, nil
}

func LoadPrivateKey(service string) (*rsa.PrivateKey, error) {
	privatePath, _ := keyPaths(service)
	return loadPrivateKeyFromFile(privatePath)
}

func LoadPublicKey(service string) (*rsa.PublicKey, error) {
	_, publicPath := keyPaths(service)
	return loadPublicKeyFromFile(publicPath)
}

func keyPaths(service string) (string, string) {
	privatePath := filepath.Join(certificatesDir, fmt.Sprintf("%s-private.pem", service))
	publicPath := filepath.Join(certificatesDir, fmt.Sprintf("%s-public.pem", service))
	return privatePath, publicPath
}

func savePrivateKey(path string, key *rsa.PrivateKey) error {
	privateBytes := x509.MarshalPKCS1PrivateKey(key)

	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateBytes,
	}

	if err := os.WriteFile(path, pem.EncodeToMemory(block), 0o600); err != nil {
		return fmt.Errorf("erro ao salvar chave privada: %w", err)
	}

	return nil
}

func savePublicKey(path string, key *rsa.PublicKey) error {
	publicBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return fmt.Errorf("erro ao serializar chave publica: %w", err)
	}

	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicBytes,
	}

	if err := os.WriteFile(path, pem.EncodeToMemory(block), 0o644); err != nil {
		return fmt.Errorf("erro ao salvar chave publica: %w", err)
	}

	return nil
}

func loadPrivateKeyFromFile(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("arquivo %s nao contem PEM valido", path)
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler chave privada: %w", err)
	}

	return key, nil
}

func loadPublicKeyFromFile(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("arquivo %s nao contem PEM valido", path)
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler chave publica: %w", err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("chave publica em %s nao e RSA", path)
	}

	return rsaPublicKey, nil
}

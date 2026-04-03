package events

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

type Package struct {
	Type      string          `json:"type"`
	Producer  string          `json:"producer"`
	Payload   json.RawMessage `json:"payload"`
	Signature string          `json:"signature"`
	Timestamp string          `json:"timestamp"`
}

type signablePackage struct {
	Type      string          `json:"type"`
	Producer  string          `json:"producer"`
	Payload   json.RawMessage `json:"payload"`
	Timestamp string          `json:"timestamp"`
}

func NewPackage(eventType, producer string, payload any) (*Package, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &Package{
		Type:      eventType,
		Producer:  producer,
		Payload:   payloadBytes,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func (e *Package) SetPayload(payload any) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	e.Payload = payloadBytes
	return nil
}

func (e *Package) DecodePayload(target any) error {
	return json.Unmarshal(e.Payload, target)
}

func (e *Package) ComputeChecksum() ([]byte, error) {
	signableBytes, err := e.signableBytes()
	if err != nil {
		return nil, err
	}

	sum := sha256.Sum256(signableBytes)
	return sum[:], nil
}

func (e *Package) Sign(privateKey *rsa.PrivateKey) error {
	if privateKey == nil {
		return fmt.Errorf("private key nao pode ser nil")
	}

	checksum, err := e.ComputeChecksum()
	if err != nil {
		return err
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, checksum)
	if err != nil {
		return err
	}

	e.Signature = base64.StdEncoding.EncodeToString(signature)
	return nil
}

func (e *Package) Verify(publicKey *rsa.PublicKey) error {
	if publicKey == nil {
		return fmt.Errorf("public key nao pode ser nil")
	}

	if e.Signature == "" {
		return fmt.Errorf("assinatura nao informada")
	}

	checksum, err := e.ComputeChecksum()
	if err != nil {
		return err
	}

	signature, err := base64.StdEncoding.DecodeString(e.Signature)
	if err != nil {
		return fmt.Errorf("assinatura invalida: %w", err)
	}

	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, checksum, signature); err != nil {
		return fmt.Errorf("falha ao verificar assinatura: %w", err)
	}

	return nil
}

func (e *Package) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

func FromJSON(data []byte) (*Package, error) {
	var pkg Package
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}

	return &pkg, nil
}

func (e *Package) signableBytes() ([]byte, error) {
	signable := signablePackage{
		Type:      e.Type,
		Producer:  e.Producer,
		Payload:   e.Payload,
		Timestamp: e.Timestamp,
	}

	return json.Marshal(signable)
}

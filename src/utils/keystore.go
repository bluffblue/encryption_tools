package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type KeyEntry struct {
	KeyID     string    `json:"key_id"`
	Key       string    `json:"key"`
	CreatedAt time.Time `json:"created_at"`
	Label     string    `json:"label,omitempty"`
}

type KeyStore struct {
	filepath string
	Keys     []KeyEntry `json:"keys"`
}

func NewKeyStore(storePath string) (*KeyStore, error) {
	ks := &KeyStore{
		filepath: storePath,
	}

	if err := os.MkdirAll(filepath.Dir(storePath), 0755); err != nil {
		return nil, err
	}

	if _, err := os.Stat(storePath); err == nil {
		data, err := os.ReadFile(storePath)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(data, &ks.Keys); err != nil {
			return nil, err
		}
	}

	return ks, nil
}

func (ks *KeyStore) AddKey(key []byte, label string) (string, error) {
	keyID := GenerateKeyID()
	entry := KeyEntry{
		KeyID:     keyID,
		Key:       EncodeKey(key),
		CreatedAt: time.Now(),
		Label:     label,
	}

	ks.Keys = append(ks.Keys, entry)
	return keyID, ks.save()
}

func (ks *KeyStore) GetKey(keyID string) ([]byte, error) {
	for _, entry := range ks.Keys {
		if entry.KeyID == keyID {
			return DecodeKey(entry.Key)
		}
	}
	return nil, os.ErrNotExist
}

func (ks *KeyStore) ListKeys() []KeyEntry {
	return ks.Keys
}

func (ks *KeyStore) save() error {
	data, err := json.MarshalIndent(ks.Keys, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ks.filepath, data, 0644)
}

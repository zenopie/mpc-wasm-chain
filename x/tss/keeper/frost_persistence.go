package keeper

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/taurusgroup/frost-ed25519/pkg/eddsa"
)

// frostKeyDir is the subdirectory within the node's home for storing FROST keys
const frostKeyDir = "data/frost_keys"

// FROSTKeyFile represents the serialized format of a FROST key share
// The eddsa types have built-in JSON marshaling
type FROSTKeyFile struct {
	KeySetID     string        `json:"keyset_id"`
	SecretShare  *eddsa.SecretShare `json:"secret_share"`
	PublicShares *eddsa.Public      `json:"public_shares"`
}

var (
	// nodeHome is set at startup to the node's home directory
	nodeHome     string
	nodeHomeLock sync.RWMutex
)

// SetNodeHome sets the node's home directory for FROST key storage
// Should be called at app startup
func SetNodeHome(home string) {
	nodeHomeLock.Lock()
	defer nodeHomeLock.Unlock()
	nodeHome = home
}

// GetFROSTKeyPath returns the path where FROST keys are stored
func GetFROSTKeyPath() string {
	nodeHomeLock.RLock()
	defer nodeHomeLock.RUnlock()
	return filepath.Join(nodeHome, frostKeyDir)
}

// saveFROSTKeyShareToDisk persists a FROST key share to disk
func saveFROSTKeyShareToDisk(keySetID string, secretShare *eddsa.SecretShare, publicShares *eddsa.Public) error {
	keyPath := GetFROSTKeyPath()

	// Create directory if it doesn't exist
	if err := os.MkdirAll(keyPath, 0700); err != nil {
		return fmt.Errorf("failed to create FROST key directory: %w", err)
	}

	// The eddsa types have built-in JSON marshaling
	keyFile := FROSTKeyFile{
		KeySetID:     keySetID,
		SecretShare:  secretShare,
		PublicShares: publicShares,
	}

	data, err := json.MarshalIndent(keyFile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal FROST key share: %w", err)
	}

	// Write to file with restrictive permissions
	filePath := filepath.Join(keyPath, keySetID+".json")
	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write FROST key file: %w", err)
	}

	return nil
}

// LoadFROSTKeyShares loads all FROST key shares from disk into memory
// Should be called at app startup after SetNodeHome
func LoadFROSTKeyShares() error {
	keyPath := GetFROSTKeyPath()

	// Check if directory exists
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		// No keys to load
		return nil
	}

	// Read all key files
	entries, err := os.ReadDir(keyPath)
	if err != nil {
		return fmt.Errorf("failed to read FROST key directory: %w", err)
	}

	frostStateManager.mu.Lock()
	defer frostStateManager.mu.Unlock()

	loadedCount := 0
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(keyPath, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Warning: failed to read FROST key file %s: %v\n", entry.Name(), err)
			continue
		}

		var keyFile FROSTKeyFile
		if err := json.Unmarshal(data, &keyFile); err != nil {
			fmt.Printf("Warning: failed to parse FROST key file %s: %v\n", entry.Name(), err)
			continue
		}

		if keyFile.SecretShare == nil || keyFile.PublicShares == nil {
			fmt.Printf("Warning: incomplete FROST key file %s\n", entry.Name())
			continue
		}

		frostStateManager.keyShares[keyFile.KeySetID] = keyFile.SecretShare
		frostStateManager.publicShares[keyFile.KeySetID] = keyFile.PublicShares
		loadedCount++
	}

	if loadedCount > 0 {
		fmt.Printf("Loaded %d FROST key share(s) from disk\n", loadedCount)
	}

	return nil
}

// DeleteFROSTKeyShare removes a FROST key share from disk
func DeleteFROSTKeyShare(keySetID string) error {
	keyPath := GetFROSTKeyPath()
	filePath := filepath.Join(keyPath, keySetID+".json")

	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete FROST key file: %w", err)
	}

	return nil
}

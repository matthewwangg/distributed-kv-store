package utils

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func LoadKeyValueDir(dir string) (map[string]string, error) {
	result := make(map[string]string)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		path := filepath.Join(dir, entry.Name())

		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", path, err)
		}

		var kv map[string]string
		if err := json.Unmarshal(data, &kv); err != nil {
			return nil, fmt.Errorf("invalid JSON in %s: %w", path, err)
		}

		for k, v := range kv {
			result[k] = v
		}
	}

	return result, nil
}

func hashKeyToUint64(key string) uint64 {
	sum := sha256.Sum256([]byte(key))
	return binary.BigEndian.Uint64(sum[:8])
}

func GetResponsiblePeer(key string, peers map[string]string) string {
	ids := make([]string, 0, len(peers))
	for id := range peers {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	h := hashKeyToUint64(key)

	index := int(h % uint64(len(ids)))

	return peers[ids[index]]
}

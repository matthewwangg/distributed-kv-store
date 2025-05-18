package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

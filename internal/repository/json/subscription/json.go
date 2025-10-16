package subscription

import (
	"encoding/json"
	"fmt"
	"os"
)

type SubscriptionRepository struct {
	filePath string
}

func NewSubscriptionRepository(filePath string) *SubscriptionRepository {
	return &SubscriptionRepository{filePath: filePath}
}

func (sr *SubscriptionRepository) loadConfig() (map[string]interface{}, error) {
	data, err := os.ReadFile(sr.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}

	return config, nil
}

func (sr *SubscriptionRepository) saveConfig(config map[string]interface{}) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(sr.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

func (sr *SubscriptionRepository) findRealityConfig(config map[string]interface{}) (map[string]interface{}, error) {
	inbounds, ok := config["inbounds"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("inbounds not found or not an array")
	}

	for _, inbound := range inbounds {
		inboundMap, ok := inbound.(map[string]interface{})
		if !ok {
			continue
		}

		tls, ok := inboundMap["tls"].(map[string]interface{})
		if !ok {
			continue
		}

		reality, ok := tls["reality"].(map[string]interface{})
		if !ok {
			continue
		}

		return reality, nil
	}

	return nil, fmt.Errorf("reality config not found")
}

func (sr *SubscriptionRepository) AddKey(key string) error {
	config, err := sr.loadConfig()
	if err != nil {
		return err
	}

	reality, err := sr.findRealityConfig(config)
	if err != nil {
		return err
	}

	shortIDInterface, exists := reality["short_id"]
	if !exists {
		// Если short_id не существует, создаем новый массив
		reality["short_id"] = []interface{}{key}
		return sr.saveConfig(config)
	}

	shortIDs, ok := shortIDInterface.([]interface{})
	if !ok {
		return fmt.Errorf("short_id is not an array")
	}

	// Проверяем, существует ли ключ уже
	for _, item := range shortIDs {
		if itemStr, ok := item.(string); ok && itemStr == key {
			return fmt.Errorf("key %s already exists", key)
		}
	}

	// Добавляем новый ключ
	reality["short_id"] = append(shortIDs, key)
	return sr.saveConfig(config)
}

func (sr *SubscriptionRepository) RemoveKey(key string) error {
	config, err := sr.loadConfig()
	if err != nil {
		return err
	}

	reality, err := sr.findRealityConfig(config)
	if err != nil {
		return err
	}

	shortIDInterface, exists := reality["short_id"]
	if !exists {
		return fmt.Errorf("short_id not found")
	}

	shortIDs, ok := shortIDInterface.([]interface{})
	if !ok {
		return fmt.Errorf("short_id is not an array")
	}

	// Создаем новый массив без удаляемого ключа
	newShortIDs := make([]interface{}, 0, len(shortIDs))
	found := false

	for _, item := range shortIDs {
		if itemStr, ok := item.(string); ok && itemStr == key {
			found = true
			continue
		}
		newShortIDs = append(newShortIDs, item)
	}

	if !found {
		return fmt.Errorf("key %s not found", key)
	}

	reality["short_id"] = newShortIDs
	return sr.saveConfig(config)
}

func (sr *SubscriptionRepository) KeyExists(key string) (bool, error) {
	config, err := sr.loadConfig()
	if err != nil {
		return false, err
	}

	reality, err := sr.findRealityConfig(config)
	if err != nil {
		return false, err
	}

	shortIDInterface, exists := reality["short_id"]
	if !exists {
		return false, nil
	}

	shortIDs, ok := shortIDInterface.([]interface{})
	if !ok {
		return false, fmt.Errorf("short_id is not an array")
	}

	for _, item := range shortIDs {
		if itemStr, ok := item.(string); ok && itemStr == key {
			return true, nil
		}
	}

	return false, nil
}

func (sr *SubscriptionRepository) GetAllKeys() ([]string, error) {
	config, err := sr.loadConfig()
	if err != nil {
		return nil, err
	}

	reality, err := sr.findRealityConfig(config)
	if err != nil {
		return nil, err
	}

	shortIDInterface, exists := reality["short_id"]
	if !exists {
		return []string{}, nil
	}

	shortIDs, ok := shortIDInterface.([]interface{})
	if !ok {
		return nil, fmt.Errorf("short_id is not an array")
	}

	result := make([]string, 0, len(shortIDs))
	for _, item := range shortIDs {
		if itemStr, ok := item.(string); ok {
			result = append(result, itemStr)
		}
	}

	return result, nil
}

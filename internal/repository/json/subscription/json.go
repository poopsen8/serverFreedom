package subscription

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	config "userServer/internal/model/config/JSON"
)

type subscriptionRepository struct {
	filename string
}

func NewSubscriptionRepository() *subscriptionRepository {
	configPath := getConfigPath("../../../config4.json") //TODO
	return &subscriptionRepository{
		filename: configPath,
	}
}

func getConfigPath(filename string) string {
	exe, _ := os.Executable()
	exeDir := filepath.Dir(exe)
	localPath := filepath.Join(exeDir, filename)

	if _, err := os.Stat(localPath); err == nil {
		return localPath
	}

	return filepath.Join("", filename)
}

// AddKey добавляет новый ключ в realitySettings первого inbound
func (sr *subscriptionRepository) AddKey(key string) error {
	config, err := sr.readConfig()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("ошибка чтения конфигурации: %v", err)
	}

	// Проверяем, что есть хотя бы один inbound
	if len(config.Inbounds) == 0 {
		return fmt.Errorf("нет inbound конфигураций")
	}

	// Получаем realitySettings первого inbound
	realitySettings := &config.Inbounds[0].StreamSettings.RealitySettings

	// Проверяем, нет ли уже такого ключа
	for _, existingKey := range realitySettings.ShortIds {
		if existingKey == key {
			return fmt.Errorf("ключ '%s' уже существует", key)
		}
	}

	// Добавляем ключ
	realitySettings.ShortIds = append(realitySettings.ShortIds, key)

	// Сохраняем изменения
	return sr.writeConfig(config)
}

// RemoveKey удаляет ключ из realitySettings первого inbound
func (sr *subscriptionRepository) RemoveKey(key string) error {
	config, err := sr.readConfig()
	if err != nil {
		return fmt.Errorf("ошибка чтения конфигурации: %v", err)
	}

	if len(config.Inbounds) == 0 {
		return fmt.Errorf("нет inbound конфигураций")
	}

	realitySettings := &config.Inbounds[0].StreamSettings.RealitySettings
	found := false

	// Ищем и удаляем ключ
	for i, existingKey := range realitySettings.ShortIds {
		if existingKey == key {
			realitySettings.ShortIds = append(realitySettings.ShortIds[:i], realitySettings.ShortIds[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("ключ '%s' не найден", key)
	}

	return sr.writeConfig(config)
}

// KeyExists проверяет наличие ключа
func (sr *subscriptionRepository) KeyExists(key string) (bool, error) {
	config, err := sr.readConfig()
	if err != nil {
		return false, fmt.Errorf("ошибка чтения конфигурации: %v", err)
	}

	if len(config.Inbounds) == 0 {
		return false, nil
	}

	realitySettings := config.Inbounds[0].StreamSettings.RealitySettings
	for _, existingKey := range realitySettings.ShortIds {
		if existingKey == key {
			return true, nil
		}
	}

	return false, nil
}

// GetAllKeys возвращает все shortIds
func (sr *subscriptionRepository) GetAllKeys() ([]string, error) {
	config, err := sr.readConfig()
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения конфигурации: %v", err)
	}

	if len(config.Inbounds) == 0 {
		return []string{}, nil
	}

	return config.Inbounds[0].StreamSettings.RealitySettings.ShortIds, nil
}

// readConfig читает и парсит JSON файл
func (sr *subscriptionRepository) readConfig() (*config.Config, error) {
	if _, err := os.Stat(sr.filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("файл конфигурации не найден: %s", sr.filename)
	}

	data, err := os.ReadFile(sr.filename)
	if err != nil {
		return nil, err
	}

	var config config.Config
	if len(data) == 0 {
		return &config, nil
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// writeConfig записывает конфигурацию в JSON файл
func (sr *subscriptionRepository) writeConfig(config *config.Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(sr.filename, data, 0644)
}

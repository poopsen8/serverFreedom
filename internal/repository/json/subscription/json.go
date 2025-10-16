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

func NewSubscriptionRepository(Path string) *subscriptionRepository {
	configPath := getConfigPath(Path)
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

// AddKey добавляет новый ключ в realitySettings trojan инбаунда
func (sr *subscriptionRepository) AddKey(key string) error {
	config, err := sr.readConfig()
	if err != nil {
		return fmt.Errorf("ошибка чтения конфигурации: %v", err)
	}

	// Находим trojan inbound
	trojanInbound, err := sr.findTrojanInbound(config)
	if err != nil {
		return err
	}

	// Проверяем, что reality настроен
	if trojanInbound.TLS == nil || trojanInbound.TLS.Reality == nil {
		return fmt.Errorf("reality не настроен в trojan inbound")
	}

	// Проверяем, нет ли уже такого ключа
	for _, existingKey := range trojanInbound.TLS.Reality.ShortID {
		if existingKey == key {
			return fmt.Errorf("ключ '%s' уже существует", key)
		}
	}

	// Добавляем ключ
	trojanInbound.TLS.Reality.ShortID = append(trojanInbound.TLS.Reality.ShortID, key)

	// Сохраняем изменения
	return sr.writeConfig(config)
}

// RemoveKey удаляет ключ из realitySettings trojan инбаунда
func (sr *subscriptionRepository) RemoveKey(key string) error {
	config, err := sr.readConfig()
	if err != nil {
		return fmt.Errorf("ошибка чтения конфигурации: %v", err)
	}

	trojanInbound, err := sr.findTrojanInbound(config)
	if err != nil {
		return err
	}

	if trojanInbound.TLS == nil || trojanInbound.TLS.Reality == nil {
		return fmt.Errorf("reality не настроен в trojan inbound")
	}

	found := false
	newShortIDs := make([]string, 0, len(trojanInbound.TLS.Reality.ShortID))

	// Ищем и удаляем ключ
	for _, existingKey := range trojanInbound.TLS.Reality.ShortID {
		if existingKey == key {
			found = true
			continue
		}
		newShortIDs = append(newShortIDs, existingKey)
	}

	if !found {
		return fmt.Errorf("ключ '%s' не найден", key)
	}

	trojanInbound.TLS.Reality.ShortID = newShortIDs
	return sr.writeConfig(config)
}

// KeyExists проверяет наличие ключа
func (sr *subscriptionRepository) KeyExists(key string) (bool, error) {
	config, err := sr.readConfig()
	if err != nil {
		return false, fmt.Errorf("ошибка чтения конфигурации: %v", err)
	}

	trojanInbound, err := sr.findTrojanInbound(config)
	if err != nil {
		return false, err
	}

	if trojanInbound.TLS == nil || trojanInbound.TLS.Reality == nil {
		return false, nil
	}

	for _, existingKey := range trojanInbound.TLS.Reality.ShortID {
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

	trojanInbound, err := sr.findTrojanInbound(config)
	if err != nil {
		return nil, err
	}

	if trojanInbound.TLS == nil || trojanInbound.TLS.Reality == nil {
		return []string{}, nil
	}

	return trojanInbound.TLS.Reality.ShortID, nil
}

// findTrojanInbound находит trojan inbound в конфигурации
func (sr *subscriptionRepository) findTrojanInbound(config *config.SingBoxConfig) (*config.Inbound, error) {
	for i := range config.Inbounds {
		if config.Inbounds[i].Type == "trojan" {
			return &config.Inbounds[i], nil
		}
	}
	return nil, fmt.Errorf("trojan inbound не найден в конфигурации")
}

// readConfig читает и парсит JSON файл sing-box
func (sr *subscriptionRepository) readConfig() (*config.SingBoxConfig, error) {
	if _, err := os.Stat(sr.filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("файл конфигурации не найден: %s", sr.filename)
	}

	data, err := os.ReadFile(sr.filename)
	if err != nil {
		return nil, err
	}

	var config config.SingBoxConfig
	if len(data) == 0 {
		return &config, nil
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %v", err)
	}

	return &config, nil
}

// writeConfig записывает конфигурацию в JSON файл sing-box
func (sr *subscriptionRepository) writeConfig(config *config.SingBoxConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка маршалинга JSON: %v", err)
	}

	return os.WriteFile(sr.filename, data, 0644)
}

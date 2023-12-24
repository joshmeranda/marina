package marina

import (
	"encoding/base64"
	"fmt"
	"os"
	"path"
)

const (
	bearerTokenFileName       = "bearer-token"
	githubAccessTokenFileName = "github-access-token"
)

func encodeToFile(filePath string, v string) error {
	tokenFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer tokenFile.Close()

	encoder := base64.NewEncoder(base64.StdEncoding, tokenFile)
	if _, err := encoder.Write([]byte(v)); err != nil {
		return err
	}

	return nil
}

func decodeFromFile(filePath string) (string, error) {
	tokenFile, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer tokenFile.Close()

	decoder := base64.NewDecoder(base64.StdEncoding, tokenFile)
	data := make([]byte, 1024)
	if _, err = decoder.Read(data); err != nil {
		return "", err
	}

	return string(data), err
}

type configManager struct {
	Root string

	ghAccessToken *string
	bearerToken   *string
}

func newDefualtConfigManager() (*configManager, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("could not determine default config dir: %w", err)
	}

	configRootDir := path.Join(userConfigDir, "marina")

	if os.MkdirAll(configRootDir, 0777); err != nil {
		return nil, fmt.Errorf("could not create config dir: %w", err)
	}

	return &configManager{
		Root: configRootDir,
	}, nil
}

func (m *configManager) SetBearerToken(token string) error {
	if m.bearerToken != nil && *m.bearerToken == token {
		return nil
	}

	m.bearerToken = &token

	tokenFilePath := path.Join(m.Root, bearerTokenFileName)

	if err := encodeToFile(tokenFilePath, *m.bearerToken); err != nil {
		return fmt.Errorf("failed to encode data to file: %w", err)
	}

	return nil
}

func (m *configManager) GetBearerToken() (string, error) {
	if m.bearerToken != nil {
		return *m.bearerToken, nil
	}

	tokenFilePath := path.Join(m.Root, bearerTokenFileName)

	value, err := decodeFromFile(tokenFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to decode token file: %w", err)
	}

	return string(value), err
}

func (m *configManager) SetGhAccessToken(token string) error {
	if m.ghAccessToken != nil && *m.ghAccessToken == token {
		return nil
	}

	m.ghAccessToken = &token

	tokenFilePath := path.Join(m.Root, githubAccessTokenFileName)

	if err := encodeToFile(tokenFilePath, *m.ghAccessToken); err != nil {
		return fmt.Errorf("failed to encode data to file: %w", err)
	}

	return nil
}

func (m *configManager) GetGhAccessToken() (string, error) {
	if m.ghAccessToken != nil {
		return *m.ghAccessToken, nil
	}

	tokenFilePath := path.Join(m.Root, githubAccessTokenFileName)

	tokenFile, err := os.Open(tokenFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer tokenFile.Close()

	decoder := base64.NewDecoder(base64.StdEncoding, tokenFile)
	data := make([]byte, 1024)
	if _, err = decoder.Read(data); err != nil {
		return "", fmt.Errorf("failed to decode token file: %w", err)
	}

	return string(data), err
}

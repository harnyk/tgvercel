package vconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

const (
	authFileName        = ".local/share/com.vercel.cli/auth.json"
	projectJSONFileName = "./.vercel/project.json"
)

type authJSON struct {
	Token string `json:"token"`
}

type projectJSON struct {
	ProjectID string `json:"projectId"`
}

type Config struct {
}

func NewConfig() *Config {
	return &Config{}
}

func (l *Config) GetAuthToken() (string, error) {
	fullAuthFileName, err := getFullAuthFileName()
	if err != nil {
		return "", err
	}
	auth, err := readJsonFromFile[authJSON](fullAuthFileName)
	if err != nil {
		return "", err
	}
	return auth.Token, nil
}

func (l *Config) GetProjectId() (string, error) {
	projectJSON, err := readJsonFromFile[projectJSON](projectJSONFileName)
	if err != nil {
		return "", err
	}
	return projectJSON.ProjectID, nil
}

func getFullAuthFileName() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %v", err)
	}
	return path.Join(home, authFileName), nil
}

func readJsonFromFile[T any](fileName string) (T, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return *new(T), err
	}
	var obj T
	err = json.Unmarshal(content, &obj)
	return obj, err
}

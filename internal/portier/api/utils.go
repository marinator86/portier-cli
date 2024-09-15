package portier

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type AuthResponse struct {
	AccessToken  string
	RefreshToken string
}

// Function to load access token from credentials file
func LoadAccessToken(home string) (AuthResponse, error) {
	credentialsFile := filepath.Join(home, "credentials.yaml")
	if _, err := os.Stat(credentialsFile); os.IsNotExist(err) {
		return AuthResponse{}, fmt.Errorf("credentials file does not exist. Please login")
	}

	// Read the file
	fileContent, err := os.ReadFile(credentialsFile)
	if err != nil {
		return AuthResponse{}, err
	}

	// Unmarshal YAML
	var credentials map[string]string
	if err := yaml.Unmarshal(fileContent, &credentials); err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{
		AccessToken:  credentials["access_token"],
		RefreshToken: credentials["refresh_token"],
	}, nil
}

type DeviceCredentials struct {
	DeviceID string `yaml:"deviceID"`
	APIKey   string `yaml:"APIKey"`
}

func LoadDeviceCredentials(home string, filename string) (*DeviceCredentials, error) {
	credentialsFile := filepath.Join(home, filename)
	if _, err := os.Stat(credentialsFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("file %s does not exist. Please register a device", credentialsFile)
	}

	// Read the file
	fileContent, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, err
	}

	// Unmarshal YAML
	var credentials DeviceCredentials
	if err := yaml.Unmarshal(fileContent, &credentials); err != nil {
		return nil, err
	}

	return &credentials, nil
}

func StoreDeviceCredentials(deviceID, apiKey, home, filename string) error {
	// Create YAML file
	file, err := os.Create(filepath.Join(home, filename))
	if err != nil {
		return err
	}
	defer file.Close()

	// Write credentials to file
	err = yaml.NewEncoder(file).Encode(DeviceCredentials{
		DeviceID: deviceID,
		APIKey:   apiKey,
	})
	return err
}
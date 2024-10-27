package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/smoqadam/kindle-send/util"
)

type config struct {
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	StorePath string `json:"storepath"`
	Password  string `json:"password"`
	Server    string `json:"server"`
	Port      int    `json:"port"`
}

const (
	DefaultTimeout    = 120
	DefaultSMTPPort   = 465
	DefaultSMTPServer = "smtp.gmail.com"
)

var instance *config

func isGmail(mail string) bool {
	return strings.Contains(mail, "@gmail.com")
}

func loadEnv() error {
	// Try to load from .env file first, ignore error if file doesn't exist
	_ = godotenv.Load()
	return nil
}

func getEnvWithFallback(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func Load() (*config, error) {
	if err := loadEnv(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %v", err)
	}

	cfg := &config{}

	// Required fields
	cfg.Sender = os.Getenv("KINDLE_SENDER_EMAIL")
	if cfg.Sender == "" {
		return nil, fmt.Errorf("KINDLE_SENDER_EMAIL is required")
	}

	cfg.Receiver = os.Getenv("KINDLE_RECEIVER_EMAIL")
	if cfg.Receiver == "" {
		return nil, fmt.Errorf("KINDLE_RECEIVER_EMAIL is required")
	}

	cfg.Password = os.Getenv("KINDLE_SENDER_PASSWORD")
	if cfg.Password == "" {
		return nil, fmt.Errorf("KINDLE_SENDER_PASSWORD is required")
	}

	// Optional fields with defaults
	cfg.Server = getEnvWithFallback("KINDLE_SMTP_SERVER", DefaultSMTPServer)
	cfg.StorePath = getEnvWithFallback("KINDLE_STORE_PATH", "")

	// Handle port conversion
	portStr := getEnvWithFallback("KINDLE_SMTP_PORT", strconv.Itoa(DefaultSMTPPort))
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid KINDLE_SMTP_PORT: %v", err)
	}
	cfg.Port = port

	// Set Gmail defaults if using Gmail
	if isGmail(cfg.Sender) && cfg.Server == DefaultSMTPServer {
		cfg.Server = DefaultSMTPServer
		cfg.Port = DefaultSMTPPort
	}

	// Encrypt password
	encryptedPass, err := Encrypt(cfg.Sender, cfg.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt password: %v", err)
	}
	cfg.Password = encryptedPass

	InitializeConfig(cfg)
	return cfg, nil
}

func InitializeConfig(c *config) {
	if instance == nil {
		instance = c
		util.Green.Println("Loaded configuration from environment")
	}
}

func GetInstance() *config {
	return instance
}

// Helper function to generate a sample .env file
func GenerateEnvTemplate() string {
	template := `# Kindle Send Configuration
KINDLE_SENDER_EMAIL=your-email@gmail.com
KINDLE_RECEIVER_EMAIL=your-kindle@kindle.com
KINDLE_SENDER_PASSWORD=your-email-password
KINDLE_SMTP_SERVER=smtp.gmail.com
KINDLE_SMTP_PORT=465
KINDLE_STORE_PATH=/path/to/store/files  # Optional
`
	return template
}

func CreateEnvTemplate(filepath string) error {
	template := GenerateEnvTemplate()
	err := os.WriteFile(filepath, []byte(template), 0644)
	if err != nil {
		return fmt.Errorf("failed to create .env template: %v", err)
	}
	util.Green.Printf("Created .env template at %s\n", filepath)
	return nil
}

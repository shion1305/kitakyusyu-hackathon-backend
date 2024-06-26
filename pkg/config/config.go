package config

import (
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

var conf *Config

type Config struct {
	Slack    Slack    `yaml:"slack"`
	SendGrid SendGrid `yaml:"sendgrid"`
	OpenAI   OpenAI   `yaml:"openai"`
	GAS      GAS      `yaml:"gas"`
}

type Slack struct {
	AppID             string `yaml:"app_id"`
	ClientID          string `yaml:"client_id"`
	ClientSecret      string `yaml:"client_secret"`
	SigningSecret     string `yaml:"signing_secret"`
	VerificationToken string `yaml:"verification_token"`
	BotOAuthToken     string `yaml:"bot_oauth_token"`
	TeamID            string `yaml:"team_id"`
}

type SendGrid struct {
	APIKey string `yaml:"api_key"`
}

type OpenAI struct {
	BaseURL string `yaml:"baseurl"`
	APIKey  string `yaml:"api_key"`
	Model   string `yaml:"model"`
}

type GAS struct {
	AppURL string `yaml:"app_url"`
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("failed to load .env file: %v", err)
	}

	envLoc := os.Getenv("ENV_LOC")
	if envLoc == "" {
		log.Fatalln("failed to get ENV_LOC")
	}
	log.Printf("loading config from %s\n", envLoc)

	file, err := os.ReadFile(envLoc)
	if err != nil {
		log.Fatalf("failed to read file, err: %v\n", err)
	}

	var loadConf Config
	if err := yaml.Unmarshal(file, &loadConf); err != nil {
		log.Fatalf("failed to unmarshal config, err: %v\n", err)
	}
	conf = &loadConf
}

func Get() Config {
	return *conf
}

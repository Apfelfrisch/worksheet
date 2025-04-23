package pkg

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	FirstName  string `json:"first-name"`
	LastName   string `json:"last-name"`
	Company    string `json:"company"`
	CompanyNo  string `json:"company-no"`
	PersonalNo string `json:"personal-no"`
	Workday    struct {
		Start  string `json:"start"`
		End    string `json:"end"`
		Breaks []struct {
			Start string `json:"start"`
			End   string `json:"end"`
		} `json:"breaks"`
	} `json:"workday"`
	Holidays []struct {
		From  string `json:"from"`
		Until string `json:"until,omitempty"`
	} `json:"holidays"`
}

func NewConfigFromFile() Config {
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	return config
}

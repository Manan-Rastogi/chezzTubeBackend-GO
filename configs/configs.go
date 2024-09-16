package configs

import (
	"log"

	"github.com/spf13/viper"
)

// Config struct to hold configuration variables
type Config struct {
	AppPort    string `mapstructure:"APP_PORT"` 
	MongoDbUri string `mapstructure:"MONGO_DB_URI"`
}

var ENV = Config{} // ENV variable to hold the configuration values

func init() {
	LoadConfig(".") // Initialize the configuration by loading the config file from the project path
}

// LoadConfig loads configuration from a specified path
func LoadConfig(path string) {
	// Set up Viper for reading env
	viper.AddConfigPath(path)  // Add the specified path to search for the config file
	viper.SetConfigName("app") // Set the name of the config file to be read
	viper.SetConfigType("env") // Set the type of the config file to be read as environment variables

	viper.AutomaticEnv() // Automatically read in environment variables

	err := viper.ReadInConfig() // Read the config file
	if err != nil {
		log.Fatalf("Failed to config .env file: %v", err.Error()) // Log an error if failed to read the config file
	}

	err = viper.Unmarshal(&ENV) // Unmarshal the config values into the ENV variable
	if err != nil {
		log.Fatalf("Failed to unmarshal .env variables: %v", err.Error()) // Log an error if failed to unmarshal the config variables
	}
}

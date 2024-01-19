package conf

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")
	configFilename := "configs/config.yaml"
	if v := os.Getenv("SIMPLYALIGN_CONFIG"); v != "" {
		configFilename = v
	}
	viper.SetConfigFile(configFilename)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatalf("Fatal error config file: %s\n", err)
	}
	viper.WatchConfig()
}

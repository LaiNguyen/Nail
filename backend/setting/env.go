package setting

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type envSettings struct {
	// MongoDB
	Host     string `envconfig:"host" default:"localhost"`
	Database string `envconfig:"database" default:"nail"`
	UserName string `envconfig:"username" default:"lai"`
	Password string `envconfig:"pasword" default:"lai"`
}

// ProjectEnvSettings instant
var ProjectEnvSettings *envSettings

func (settings *envSettings) readEnvironmentVariables() error {
	err := envconfig.Process("settings", settings)
	if err != nil {
		return err
	}
	return nil
}

//EnvSettingsInit Always use this function to initialize a Settings struct
func EnvInit() {
	ProjectEnvSettings = &envSettings{}
	ProjectEnvSettings.readEnvironmentVariables()

	log.Printf("Environment settings are %v\n", ProjectEnvSettings)
}

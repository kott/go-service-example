package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/kott/go-service-example/pkg/api"
)

const (
	servicesProfile = "SERVICES_PROFILE"
	dockerProfile   = "docker"
)

func init() {
	profile := os.Getenv(servicesProfile)
	if profile == dockerProfile { // all env vars are set in the container already
		viper.AutomaticEnv()
		viper.AllowEmptyEnv(true)
	} else {
		viper.SetConfigName(strings.ToLower(profile))
		viper.SetConfigType("env")
		viper.AddConfigPath(".")

		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("an error occurred reading the config file: %s ", err))
		}
	}
}

func main() {
	api.Start(&api.Config{
		DBHost:       viper.GetString("db_host"),
		DBPort:       viper.GetInt("db_port"),
		DBUser:       viper.GetString("db_user"),
		DBPassword:   viper.GetString("db_password"),
		DBName:       viper.GetString("db_name"),
		RunMigration: viper.GetBool("run_migration"),

		AppHost: viper.GetString("host"),
		AppPort: viper.GetInt("port"),
	})
}

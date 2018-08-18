package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

var (
	configDir  = filepath.Join(os.Getenv("HOME"), ".dockmon")
	configFile = filepath.Join(configDir, "config.json")
)

type Config struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (conf Config) Valid() bool {
	return conf.Host != "" && conf.Username != "" && conf.Password != ""
}

// Save stores configuration.
func (config Config) Save() error {
	content, err := json.Marshal(config)
	if err != nil {
		return err
	}
	os.Mkdir(configDir, 0755)
	return ioutil.WriteFile(configFile, content, 0755)
}

// ConfigureCommand configures dockmon for use.
func ConfigureCommand() cli.Command {
	return cli.Command{
		Name:   "configure",
		Usage:  fmt.Sprintf("Configures %s for use", APP_NAME),
		Action: Configure,
	}
}

// Configure gets and stores api configuration information from the user.
func Configure(c *cli.Context) error {
	config := getApiConfig()
	api := NewApiClient(config)
	api.Login()

	err := config.Save()
	if err != nil {
		fmt.Println("Could not save configuration")
		fmt.Println(err)
		return err
	}
	return nil
}

// getApiConfig Prompts the user to input api configuration
func getApiConfig() Config {
	config, err := getConfig()
	if err != nil {
		config = Config{}
	}
	fmt.Println("Enter dockmon configuration")
	config.Host = getInput("Host", config.Host)
	config.Username = getInput("Username", config.Username)
	config.Password = getInput("Password", "")
	return config
}

// getInput Gets user input from stdin
func getInput(key, defaultValue string) string {
	printKeyAndDefaultValue(key, defaultValue)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Unable to read the value for '%s'\n", key)
		os.Exit(1)
	}
	value := strings.Replace(text, "\n", "", -1)
	if value == "" {
		return defaultValue
	}
	return value
}

func printKeyAndDefaultValue(key, defaultValue string) {
	if defaultValue == "" {
		fmt.Printf("%s: ", key)
		return
	}
	fmt.Printf("%s (%s): ", key, defaultValue)
}

// getConfig gets configuration.
func getConfig() (Config, error) {
	var config Config
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(b, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

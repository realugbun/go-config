package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Default values for the application
const (
	defaultFoo            = "defaultFoo"
	defaultBar            = "defaultBar"
	defaultBaz            = "defaultBaz"
	defaultConfigLocation = "./democonfig.yml"
	appName               = "democonfig"
)

// Help message displayed when passing --help or -h at startup
const (
	helpMessage = `
	--help	 	-h		Display this help message.
	--foo= 		-f=		Set foo value.			Default: defaultFoo
	--bar= 		-b=		Set bar value.			Default: defaultBar
	--baz= 		-z=		Set baz value.			Default: defaultBaz
	--config=	-c=		Path to config file.		Default: ./democonfig.yml`
)

// SettingsRec holds the configuration information
type SettingsRec struct {
	Foo            string `yaml:"foo"`
	Bar            string `yaml:"bar"`
	Baz            string `yaml:"baz"`
	ConfigLocation string
}

// These could be loaded by building the app with ldflags
var (
	version  = "0.0.1"
	buildDay = "2021-08-01"
)

// Error messages from this package
var (
	errNoFile        = errors.New("no config file found...using default values")
	errUnmarshalYAML = errors.New("unable to read config file")
)

// Settings holds the configuration settings loaded at startup
var (
	Settings = SettingsRec{}
)

// Load loads configuration from startup arguments and a config file. The order of presidence is 1) arguments, 2) config file, and 3) default values.
func (s *SettingsRec) Load() error {
	s.LoadStartupArgs()
	err := s.LoadConfigFile()
	if err != nil {
		return err
	}
	return nil
}

// LoadStartupArgs processes arguments passed at startup
func (s *SettingsRec) LoadStartupArgs() {

	// If there are no startup arguments return
	if len(os.Args) == 1 {
		return
	}

	for i, v := range os.Args {
		// Skip the first os.Arg as it is not relevant
		if i == 0 {
			continue
		}

		// Check for arguments that exit the application
		switch v {
		case "--version", "-V":
			fmt.Printf("%s %s built on %s\n", appName, version, buildDay)
			os.Exit(0)
		case "--help", "-h":
			fmt.Printf("%s %s built on %s\n", appName, version, buildDay)
			fmt.Println(helpMessage)
			os.Exit(0)
		}

		// Check for input values passed as arguments
		switch {
		case strings.HasPrefix(v, "--foo="), strings.HasPrefix(v, "-f="):
			arg := strings.SplitAfter(v, "=")
			s.Foo = arg[1]
		case strings.HasPrefix(v, "--bar="), strings.HasPrefix(v, "-b="):
			arg := strings.SplitAfter(v, "=")
			s.Bar = arg[1]
		case strings.HasPrefix(v, "--baz="), strings.HasPrefix(v, "-z="):
			arg := strings.SplitAfter(v, "=")
			s.Baz = arg[1]
		case strings.HasPrefix(v, "--config="), strings.HasPrefix(v, "-c="):
			arg := strings.SplitAfter(v, "=")
			s.ConfigLocation = arg[1]
		default:
			// exit on unknown argument
			fmt.Println("unknown argument: " + v)
			os.Exit(1)
		}
	}
}

// LoadConfigFile loads configuration from a config file. The order of presidence is 1) arguments, 2) config file, and 3) default values.
func (s *SettingsRec) LoadConfigFile() (err error) {

	defer func() {
		// if the function exits with an error, the os arg values will be used
		// if the os arg values are empty, the default values will be used.
		if err != nil {
			if s.Foo == "" {
				s.Foo = defaultFoo
			}
			if s.Bar == "" {
				s.Bar = defaultBar
			}
			if s.Baz == "" {
				s.Baz = defaultBaz
			}
		}
	}()

	// If no configuration location was passed, use the default location
	if s.ConfigLocation == "" {
		s.ConfigLocation = defaultConfigLocation
	}

	// Read the config file
	file, err := ioutil.ReadFile(s.ConfigLocation)
	if err != nil {
		return errNoFile
	}

	fileConfig := SettingsRec{}
	err = yaml.Unmarshal(file, &fileConfig)
	if err != nil {
		return errUnmarshalYAML
	}

	// Replace values not in config file with default values
	if fileConfig.Foo == "" {
		fileConfig.Foo = defaultFoo
	}
	if fileConfig.Bar == "" {
		fileConfig.Bar = defaultBar
	}
	if fileConfig.Baz == "" {
		fileConfig.Baz = defaultBaz
	}

	// Replace values not passed at startup with config file or default values
	if s.Foo == "" {
		s.Foo = fileConfig.Foo
	}
	if s.Bar == "" {
		s.Bar = fileConfig.Bar
	}
	if s.Baz == "" {
		s.Baz = fileConfig.Baz
	}

	return nil
}

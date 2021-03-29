package configor

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path"
	"reflect"
	"regexp"
	"strings"
)

type Configor struct {
	*Config
}

type Config struct {
	Environment       string
	EnvironmentPrefix string

	// In case of json files, this field will be used only when compiled with
	// go 1.10 or later.
	// This field will be ignored when compiled with go versions lower than 1.10.
	ErrorOnUnmatchedKeys bool
}

// Load will unmarshal configurations to struct from files that you provide
func Load(config interface{}, files ...string) error {
	return New(nil).Load(config, files...)
}

func New(config *Config) *Configor {
	if config == nil {
		config = &Config{}
	}
	return &Configor{Config: config}
}

var testRegexp = regexp.MustCompile("_test|(\\.test$)")

func (c *Configor) GetEnvironment() string {
	if c.Environment == "" {
		if env := os.Getenv("CONFIGOR_ENV"); env != "" {
			return env
		}

		if testRegexp.MatchString(os.Args[0]) {
			return "test"
		}

		return "development"
	}
	return c.Environment
}

func (c *Configor) GetEnvironmentPrefix() string {
	if c.EnvironmentPrefix == "" {
		return os.Getenv("CONFIGOR_ENV_PREFIX")
	}
	return c.EnvironmentPrefix
}

// GetErrorOnUnmatchedKeys returns a boolean indicating if an error should be
// thrown if there are keys in the config file that do not correspond to the
// config struct
func (c *Configor) GetErrorOnUnmatchedKeys() bool {
	return c.ErrorOnUnmatchedKeys
}

// Load will unmarshal configurations to struct from files that you provide
func (c *Configor) Load(config interface{}, files ...string) (err error) {
	defaultValue := reflect.Indirect(reflect.ValueOf(config))
	if !defaultValue.CanAddr() {
		return fmt.Errorf("Config %v should be addressable", config)
	}

	configFiles := c.getConfigurationFiles(files...)

	for _, file := range configFiles {
		if err := UnmarshalFile(config, file, c.ErrorOnUnmatchedKeys); err != nil {
			return err
		}
	}

	prefix := c.GetEnvironmentPrefix()
	if prefix == "" {
		return c.processTags(config)
	}
	return c.processTags(config, prefix)
}

// UnmatchedTomlKeysError errors are returned by the Load function when
// ErrorOnUnmatchedKeys is set to true and there are unmatched keys in the input
// toml config file. The string returned by Error() contains the names of the
// missing keys.
type UnmatchedTomlKeysError struct {
	Keys []toml.Key
}

func (e *UnmatchedTomlKeysError) Error() string {
	return fmt.Sprintf("There are keys in the config file that do not match any field in the given struct: %v", e.Keys)
}

func getFilenameWithENVPrefix(file, env string) (string, error) {
	var (
		envFile string
		extname = path.Ext(file)
	)

	if extname == "" {
		envFile = fmt.Sprintf("%v.%v", file, env)
	} else {
		envFile = fmt.Sprintf("%v.%v%v", strings.TrimSuffix(file, extname), env, extname)
	}

	if fileInfo, err := os.Stat(envFile); err == nil && fileInfo.Mode().IsRegular() {
		return envFile, nil
	}
	return "", fmt.Errorf("failed to find file %v", file)
}

func (c *Configor) getConfigurationFiles(files ...string) []string {
	var results []string

	for i := len(files) - 1; i >= 0; i-- {
		foundFile := false
		file := files[i]

		// check configuration
		if fileInfo, err := os.Stat(file); err == nil && fileInfo.Mode().IsRegular() {
			foundFile = true
			results = append(results, file)
		}

		// check configuration with env
		if file, err := getFilenameWithENVPrefix(file, c.GetEnvironment()); err == nil {
			foundFile = true
			results = append(results, file)
		}

		// check example configuration
		if !foundFile {
			if example, err := getFilenameWithENVPrefix(file, "example"); err == nil {
				results = append(results, example)
			}
		}
	}
	return results
}

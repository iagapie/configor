# Configor

Golang Configuration tool that support YAML, JSON, TOML, Shell Environment (Supports Go 1.10+)

## Usage

```go
package main

import (
	"fmt"
	"github.com/jinzhu/configor"
)

var Config = struct {
	APPName string `default:"app name"`

	DB struct {
		Name     string
		User     string `default:"root"`
		Password string `required:"true" env:"DBPassword"`
		Port     uint   `default:"3306"`
	}

	Contacts []struct {
		Name  string
		Email string `required:"true"`
	}
}{}

func main() {
	configor.Load(&Config, "config.yml")
	fmt.Printf("config: %#v", Config)
}
```

With configuration file *config.yml*:

```yaml
appname: test

db:
    name:     test
    user:     test
    password: test
    port:     1234

contacts:
- name: i test
  email: test@test.com
```

# Usage

* Load mutiple configurations

```go
// Earlier configurations have higher priority
configor.Load(&Config, "application.yml", "database.json")
```

* Return error on unmatched keys

Return an error on finding keys in the config file that do not match any fields in the config struct.
In the example below, an error will be returned if config.toml contains keys that do not match any fields in the ConfigStruct struct.
If ErrorOnUnmatchedKeys is not set, it defaults to false.

Note that for json files, setting ErrorOnUnmatchedKeys to true will have an effect only if using go 1.10 or later.

```go
err := configor.New(&configor.Config{ErrorOnUnmatchedKeys: true}).Load(&ConfigStruct, "config.toml")
```

* Load configuration by environment

Use `CONFIGOR_ENV` to set environment, if `CONFIGOR_ENV` not set, environment will be `development` by default, and it will be `test` when running tests with `go test`

```go
// config.go
configor.Load(&Config, "config.json")

$ go run config.go
// Will load `config.json`, `config.development.json` if it exists
// `config.development.json` will overwrite `config.json`'s configuration
// You could use this to share same configuration across different environments

$ CONFIGOR_ENV=production go run config.go
// Will load `config.json`, `config.production.json` if it exists
// `config.production.json` will overwrite `config.json`'s configuration

$ go test
// Will load `config.json`, `config.test.json` if it exists
// `config.test.json` will overwrite `config.json`'s configuration

$ CONFIGOR_ENV=production go test
// Will load `config.json`, `config.production.json` if it exists
// `config.production.json` will overwrite `config.json`'s configuration
```

```go
// Set environment by config
configor.New(&configor.Config{Environment: "production"}).Load(&Config, "config.json")
```

* Example Configuration

```go
// config.go
configor.Load(&Config, "config.yml")

$ go run config.go
// Will load `config.example.yml` automatically if `config.yml` not found and print warning message
```

* Load From Shell Environment

```go
$ APPNAME="hello world" DB_NAME="hello world" go run config.go
// Load configuration from shell environment, it's name is {{prefix}}_FieldName
```

```go
// You could overwrite the prefix with environment CONFIGOR_ENV_PREFIX, for example:
$ CONFIGOR_ENV_PREFIX="WEB" WEB_APPNAME="hello world" WEB_DB_NAME="hello world" go run config.go

// Set prefix by config
configor.New(&configor.Config{ENVPrefix: "WEB"}).Load(&Config, "config.json")
```

* Anonymous Struct

Add the `anonymous:"true"` tag to an anonymous, embedded struct to NOT include the struct name in the environment
variable of any contained fields.  For example:

```go
type Details struct {
	Description string
}

type Config struct {
	Details `anonymous:"true"`
}
```

With the `anonymous:"true"` tag specified, the environment variable for the `Description` field is `DESCRIPTION`.
Without the `anonymous:"true"`tag specified, then environment variable would include the embedded struct name and be `DETAILS_DESCRIPTION`.

## Contributing

You can help to make the project better, check out [http://gorm.io/contribute.html](http://gorm.io/contribute.html) for things you can do.

## Author

**jinzhu** (Original)

* <http://github.com/jinzhu>
* <wosmvp@gmail.com>
* <http://twitter.com/zhangjinzhu>

**iagapie**

Forked and updated to fit personal use-cases, 2021.

## License

Released under the MIT License

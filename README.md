# EnvLoader

![Version](https://img.shields.io/badge/version-1.0.0-blue.svg?cacheSeconds=2592000)

Secret manager integration to get application environment variables. Any variable not set by the secret manager will be set by OS envs.

Jump To:

- [EnvLoader](#envloader)
  - [Getting Started](#getting-started)
    - [Installing](#installing)
    - [Go Modules](#go-modules)
  - [Quick Examples](#quick-examples)
    - [Complete SDK Example](#complete-sdk-example)

## Getting Started

### Installing

Use `go get` to retrieve the SDK to add it to your `GOPATH` workspace, or
project's Go module dependencies.

    go get github.com/thrcorrea/envloader

To update the module use `go get -u` to retrieve the latest version of the module.

    go get -u github.com/thrcorrea/envloader

### Go Modules

If you are using Go modules, your `go get` will default to the latest tagged
release version of the module. To get a specific release version of the module use
`@<tag>` in your `go get` command.

    go get github.com/thrcorrea/envloader@v1.0.0

To get the latest module repository change use `@latest`.

    go get github.com/thrcorrea/envloader@latest

## Quick Examples

This example shows a complete working Go file which will make a struct for environment variables and initialize with this module.

### Complete SDK Example

REGION is setting automatically on AWS environment, only setting if needs to be different.

```plain
SECRET_NAME=aws-secret-manager-name
REGION=aws-secret-manager-region
```

```go
package main

import (
    "github.com/thrcorrea/envloader"
)

type Environment struct {
    DBHost          string `env:"DB_HOST"`
    DBName          string `env:"DB_NAME"`
    DBPort          string `env:"DB_PORT"`
    DBUser          string `env:"DB_USER"`
    DBPassword      string `env:"DB_PSWD"`
    RabbitHost      string `env:"AMQP_HOST"`
    RabbitUser      string `env:"AMQP_USER"`
    RabbitPassword  string `env:"AMQP_PSWD"`
    AmqpVHost       string `env:"AMQP_VHOST,optional"` // optional = disable empty validation
    AmqpProtocol    string `env:"AMQP_PROTOCOL,optional,default=AMQP"` //default, set a default value if key not setted on env and secret
}

var env Environment

func main() {
    // will run only if SECRET_NAME and REGION is setted in dotenv file
    // load environment with dotenv values and secret manager values
    // default file name is ".env" and doesn't need to be sent
    err := envloader.Load(&env, ".env.test")
    if err != nil {
        // failed to load environment variables
        panic(err)
    }
}
```

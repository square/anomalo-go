# Anomalo-Go
A Go client for Anomalo.

This package exists primarily to support the [Terraform provider](https://github.com/square/terraform-provider-anomalo). Consequently, it implements only a subset of the Anomalo API.

[![Go Reference](https://pkg.go.dev/badge/github.com/square/anomalo-go/anomalo.svg)](https://pkg.go.dev/github.com/square/anomalo-go/anomalo)

## Table of Contents

- [Anomalo-Go](#anomalo-go)
  - [Table of Contents](#table-of-contents)
  - [Installation](#installation)
  - [Getting Started](#getting-started)
    - [Quick Start](#quick-start)
  - [Documentation](#documentation)
  - [Appendix](#appendix)


## Installation

Run `go install github.com/square/anomalo-go`, then (if using modules) `go get && go mod tidy`

## Getting Started
### Quick Start

anomalo-go makes it easier to call the Anomalo API. It includes structs for request & response objects and some convenience methods that make the API easier to work with.

To get started, you'll need to provide Anomalo API credentials. If you don't have them already, ask your Anomalo administrator or representative.

A minimal example looks like this:

```go
package main

import "github.com/square/anomalo-go"

func main() {
	client, _ := anomalo.LoadClient() // Checks for credentials anomalo_secrets.json and environment variables.
	fmt.Println(client.Ping())
}
```

The package will first check for a JSON formatted credential file at `anomalo_secrets.json`. If the file does not exist, it will look for environment variables named `ANOMALO_API_SECRET_TOKEN` and `ANOMALO_INSTANCE_HOST`. The `anomalo_secrets.json` file should have the following format:

```json
{
  "token": "thisIsAToken",
  "host": "https://anomalo.example.com"
}
```

Executing this example should print a struct containing the word "Pong".

## Documentation

Refer to Anomalo documentation for most the behavior or most methods. The code in
[client.go](https://github.com/square/anomalo-go/blob/master/client.go) contains documentation for methods like `GetCheckByStaticID` that are not available in the Anomalo API.

## Appendix

Contributions are welcome. This package is not maintained by Anomalo, and is not guaranteed to be up to date with changes Anomalo makes to their API.

Brought to you by Square <img src="https://avatars.githubusercontent.com/u/82592" alt="GitHub logo" width="20" style="float: left; margin-right: 5px;"/>

# Updown Go Client

[![CI](https://github.com/sergo-techhub/updown/actions/workflows/ci.yml/badge.svg)](https://github.com/sergo-techhub/updown/actions/workflows/ci.yml)
[![Release](https://github.com/sergo-techhub/updown/actions/workflows/release.yml/badge.svg)](https://github.com/sergo-techhub/updown/actions/workflows/release.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/sergo-techhub/updown.svg)](https://pkg.go.dev/github.com/sergo-techhub/updown)
[![License: MIT](https://img.shields.io/badge/License-MIT-orange.svg)](https://github.com/sergo-techhub/updown/blob/main/LICENSE.md)

A Go client library for [updown.io](https://updown.io) - a simple and affordable website monitoring service.

## Fork Notice

This project was forked from [antoineaugusti/updown](https://github.com/antoineaugusti/updown) by [SERGO GmbH](https://github.com/sergo-techhub).

We are actively modernizing and updating this module to support the current [updown.io API](https://updown.io/api) implementation, including:

- Support for all check types (`http`, `https`, `icmp`, `tcp`, `tcps`)
- HTTP verb configuration (`GET`, `HEAD`, `POST`, `PUT`, `PATCH`, `DELETE`, `OPTIONS`)
- HTTP body for POST/PUT/PATCH requests
- Recipients API support
- Modern Go version (1.24+)

## Contributing

We'd love to see you contribute! Whether it's:

- Reporting bugs or suggesting features via [Issues](https://github.com/sergo-techhub/updown/issues)
- Submitting [Pull Requests](https://github.com/sergo-techhub/updown/pulls) with improvements or fixes
- Improving documentation

All contributions are welcome!

## Installation

```bash
go get github.com/sergo-techhub/updown
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/sergo-techhub/updown"
)

func main() {
    // Your API key from https://updown.io/settings/edit
    client := updown.NewClient("your-api-key", nil)

    // List all checks
    checks, _, err := client.Check.List()
    if err != nil {
        panic(err)
    }

    for _, check := range checks {
        fmt.Printf("%s: %s\n", check.Alias, check.URL)
    }
}
```

## Usage Examples

### Creating a Client

```go
client := updown.NewClient("your-api-key", nil)
```

### Working with Checks

```go
// List all checks
checks, _, err := client.Check.List()

// Get a check by token
check, _, err := client.Check.Get("token")

// Get token for a check alias
token, err := client.Check.TokenForAlias("My Website")

// Create a new HTTP check
item := updown.CheckItem{
    URL:   "https://example.com",
    Alias: "Example Website",
}
check, _, err := client.Check.Add(item)

// Create an ICMP ping check
item := updown.CheckItem{
    URL:  "192.168.1.1",
    Type: "icmp",
}
check, _, err := client.Check.Add(item)

// Create a TCP port check
item := updown.CheckItem{
    URL:  "tcp://db.example.com:5432",
    Type: "tcp",
}
check, _, err := client.Check.Add(item)

// Create an HTTP POST check
item := updown.CheckItem{
    URL:      "https://api.example.com/health",
    HttpVerb: "POST",
    HttpBody: `{"check": true}`,
    CustomHeaders: map[string]string{
        "Content-Type": "application/json",
    },
}
check, _, err := client.Check.Add(item)

// Update a check
updated := updown.CheckItem{URL: "https://new-url.example.com"}
check, _, err := client.Check.Update("token", updated)

// Delete a check
deleted, _, err := client.Check.Remove("token")
```

### Working with Recipients

```go
// List all recipients
recipients, _, err := client.Recipient.List()

// Create a recipient
item := updown.RecipientItem{
    Type:  updown.RecipientTypeEmail,
    Value: "alerts@example.com",
}
recipient, _, err := client.Recipient.Add(item)

// Delete a recipient
deleted, _, err := client.Recipient.Remove("recipient-id")
```

### Working with Downtimes

```go
// List downtimes for a check (paginated, 100 per page)
downtimes, _, err := client.Downtime.List("token", 1)
```

### Working with Metrics

```go
token := "your-check-token"
group := "host"
from := "2024-01-01 00:00:00 +0000"
to := "2024-01-31 23:59:59 +0000"
metrics, _, err := client.Metric.List(token, group, from, to)
```

### Working with Nodes

```go
// Get IPv4 addresses of monitoring nodes
ipv4, _, err := client.Node.ListIPv4()

// Get IPv6 addresses of monitoring nodes
ipv6, _, err := client.Node.ListIPv6()
```

## API Reference

For the complete updown.io API documentation, visit: https://updown.io/api

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

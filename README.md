# TSDNS-Go

[![Go Report Card](https://goreportcard.com/badge/github.com/honeybbq/tsdns-go)](https://goreportcard.com/report/github.com/honeybbq/tsdns-go)
[![GoDoc](https://godoc.org/github.com/honeybbq/tsdns-go?status.svg)](https://godoc.org/github.com/honeybbq/tsdns-go)
[![License](https://img.shields.io/github/license/honeybbq/tsdns-go.svg)](https://github.com/honeybbq/tsdns-go/blob/main/LICENSE)

A lightweight Teamspeak DNS server implementation in Go.

## 🌟 Highlights

- Simple and efficient TCP-based DNS server
- Flexible repository interfaces (PostgreSQL & File storage supported)
- In-memory cache with automatic updates
- Builder pattern for easy configuration
- Structured logging support

## 📋 Table of Contents

- [TSDNS-Go](#tsdns-go)
  - [🌟 Highlights](#-highlights)
  - [📋 Table of Contents](#-table-of-contents)
  - [🚀 Installation](#-installation)
  - [🎯 Quick Start](#-quick-start)
  - [⚙️ Configuration](#️-configuration)
  - [🏗 Architecture](#-architecture)
    - [Repository Interface](#repository-interface)
  - [🤝 Contributing](#-contributing)
  - [📄 License](#-license)
  - [📞 Support](#-support)

## 🚀 Installation

```bash
go get github.com/honeybbq/tsdns-go
```

## 🎯 Quick Start

Here's a minimal example to get you started: [Postgres Demo](./example/demo.go)

## ⚙️ Configuration

TSDNS supports various configuration options through its builder pattern:

```go
server := tsdns.NewServer("0.0.0.0").
    WithRepository(repo).
    WithLogger(customLogger).
    MustBuild()
```


## 🏗 Architecture

TSDNS follows a clean architecture pattern with the following components:

- **Server**: Core DNS server implementation
- **Repository**: Interface for data storage (PostgreSQL/File implementations provided)
- **Cache**: In-memory cache with automatic updates
- **Logger**: Structured logging interface

### Repository Interface

```go
type RecordRepository interface {
    Find() ([]*Record, error)
    FindByDomain(domain string) (*Record, error)
    Create(record *Record) error
    Delete(domain string) error
    DeleteByInstanceID(instanceID int64) error
    Close() error
}
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📞 Support

If you have any questions or need help, please:

1. Check the [issues](https://github.com/honeybbq/tsdns-go/issues) page
2. Create a new issue if your problem is not already reported
3. Join our [Teamspeak](teamspeak://chenkr.cn) server for live support

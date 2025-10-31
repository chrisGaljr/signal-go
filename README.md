# Messenger to Signal Bridge

A Go-based application that automatically forwards Messenger messages to Signal, providing a useful notification for those for those who value privacy but need this meta product to communicate with their contacts.

## Features

-   **Real-time Message Forwarding**: Automatically detects and forwards new Messenger messages to Signal
-   **Web Dashboard**: Monitor application status, view incident history, and manage configurations
-   **Secure Encryption**: Built-in encryption utilities for sensitive data handling
-   **Docker Support**: Containerized deployment with headless Chrome for web automation

## Installation

### Environment Variables

```bash
CIPHER_KEY=<your_cipher_key_here>
EMAIL=<your_messenger_email>
PASSWORD=<your_password_encrypted>
PIN=<your_pin_encrypted>
MY_NUMBER=<signal_target_number>
BACKEND_NUMBER=<signal_accounts_phone_number>
SIGNAL_REST_BASE_URL=<base_url_of_my_signal_cli>
MONGODB_URI=<mongo_db_url>
```

### Local Development

```bash
git clone <repository-url>
cd signal

go mod tidy

go run .
```

### Docker Deployment

```bash
docker build -t signal-go .

docker run -p 7777:7777 \
  -e CIPHER_KEY=<your_cipher_key_here> \
  -e EMAIL=<your_messenger_email> \
  -e PASSWORD=<your_password_encrypted> \
  -e PIN=<your_pin_encrypted> \
  -e MY_NUMBER=<signal_target_number> \
  -e BACKEND_NUMBER=<signal_accounts_phone_number> \
  -e SIGNAL_REST_BASE_URL=<base_url_of_my_signal_cli> \
  -e MONGODB_URI=<mongo_db_url> \
  -d signal-go
```

## API Endpoints

| Endpoint            | Method | Description             |
| ------------------- | ------ | ----------------------- |
| `/`                 | GET    | Home dashboard          |
| `/about`            | GET    | Application information |
| `/status`           | GET    | System health status    |
| `/incident-history` | GET    | View incident logs      |
| `/encrypt`          | POST   | Encrypt sensitive data  |

## Project Structure

```
messenger/
├── main.go                         # Application entry point
├── Dockerfile
├── go.mod
├── internal/
│   ├── server.go                   # HTTP server setup
│   ├── handlers/
│   │   ├── about_handler.go
│   │   ├── encrypt_handler.go
│   │   ├── home_handler.go
│   │   ├── incident_history_handler.go
│   │   └── status_handler.go
│   ├── models/
│   │   ├── config.go
│   │   ├── error_log.go
│   │   └── setup.go
│   ├── services/
│   │   ├── messenger.go            # Messenger automation
│   │   └── signal.go               # Signal API integration
│   └── utils/
│       ├── encrypt.go
│       ├── errors.go
│       └── status.go
├── static/
│   ├── css/
│   │   └── styles.css
│   └── images/
└── template/
    ├── index.html
    ├── about/
    ├── fourOhFour/
    ├── incidentHistory/
    └── status/
```

## How It Works

1. **Browser Automation**: Uses ChromeDP to automate a headless Chrome browser for Messenger interaction
2. **Message Detection**: Monitors network requests to detect new incoming messages
3. **Message Processing**: Extracts message content and sender information
4. **Signal Forwarding**: Sends processed messages to Signal via REST API
5. **Monitoring**: Errors are collected and saved for debugging purposes

## Disclaimer

This application is for educational and personal use only. Ensure compliance with Facebook's Terms of Service and Signal's API usage policies when using this software. All rights reserved.

# Coup Game (Wi-Fi Edition)

A clean and professional Go backend for a local multiplayer Coup game, designed to run over Wi-Fi. The system prioritizes maintainable code, modular architecture, and internationalization support.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Internationalization (i18n)](#internationalization-i18n)
- [How to Run](#how-to-run)
- [Development Guidelines](#development-guidelines)
- [License](#license)

## Overview

This project provides a backend for a digital Coup game playable locally by friends on the same Wi-Fi network. Players access the frontend via their mobile browsers, while all the game logic, state management, and WebSocket communication run on a Go server.

## Features

- **Local Multiplayer (Wi-Fi)**: All players connect through the same network
- **WebSocket Real-Time Communication**: Fast, low-latency actions and state sync
- **Internationalization**: Supports both English and Portuguese, selected by the host before the match
- **Modular, Clean Codebase**: Idiomatic Go; 100% code/commentary in English
- **Extensible Architecture**: Easily maintainable; ready for features like reconnection or multiple lobbies

## Tech Stack

| Layer    | Technology                           |
|----------|--------------------------------------|
| Backend  | Go (Golang), gorilla/websocket       |
| i18n     | go-i18n                             |
| Frontend | Static web (HTML/JS/CSS, not included) |

## Project Structure

```
/coup-game
├── cmd/                # Main entry point
├── internal/
│   ├── game/           # Game logic (deck, rules, turn)
│   ├── lobby/          # Lobby and player management
│   ├── ws/             # WebSocket handlers
│   └── i18n/           # Internationalization (go-i18n)
│       └── locales/    # JSON translation files (en, pt)
├── pkg/                # Shared utils/helpers
├── web/                # Static frontend (HTML, JS, CSS)
├── go.mod
├── LICENSE
├── README.md
└── README-ptbr.md
```

## Internationalization (i18n)

- **Languages**: English or Portuguese, selected by the host
- **Translations**: All user-facing text comes from translation bundles (`/internal/i18n/locales/`)
- **Implementation**: Uses go-i18n for loading and serving translations
- **Code**: All logic and comments remain in English, for global developer best practices

## How to Run

1. **Clone the repository**:
   ```bash
   git clone https://github.com/leoferamos/coup-game.git
   cd coup-game
   ```

2. **Download dependencies**:
   ```bash
   go mod tidy
   ```

3. **Run the server**:
   ```bash
   go run ./cmd
   ```

4. **Access on browser**:
   Open `http://localhost:8080` or use your host IP on other devices connected to the same Wi-Fi.

## Development Guidelines

- All new code must be fully idiomatic Go with explicit types and zero ambiguous naming
- No business logic is ever mixed in HTTP or WebSocket handlers
- User-facing text is never hardcoded; always loaded from i18n translation files
- PRs must include tests and clear descriptions
- All contributions and documentation must be in English

## License

MIT

---

*Also available in [Portuguese](README-ptbr.md)*
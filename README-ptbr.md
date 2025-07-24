# Coup Game (Edição Wi-Fi)

Backend em Go, limpo e profissional, para o jogo Coup multiplayer local, planejado para rodar em Wi-Fi. O sistema prioriza arquitetura modular, código limpo e suporte a internacionalização.

## Índice

- [Visão Geral](#visão-geral)
- [Funcionalidades](#funcionalidades)
- [Tecnologias](#tecnologias)
- [Estrutura do Projeto](#estrutura-do-projeto)
- [Internacionalização (i18n)](#internacionalização-i18n)
- [Como Rodar](#como-rodar)
- [Diretrizes de Desenvolvimento](#diretrizes-de-desenvolvimento)
- [Licença](#licença)

## Visão Geral

Crie sua sala, compartilhe o Wi-Fi e jogue Coup digitalmente com amigos. Toda a lógica, estado e comunicação são gerenciados pelo backend Go — os jogadores acessam via navegador no celular.

## Funcionalidades

- **Multiplayer Local (Wi-Fi)**: Todos conectados à mesma rede
- **Comunicação WebSocket**: Baixa latência e sincronização rápida entre jogadores
- **Internacionalização**: Suporte total a inglês e português, escolhido pelo host
- **Código Modular e Limpo**: 100% do código e comentários em inglês
- **Arquitetura Extensível**: Fácil de manter e pronta para novos recursos

## Tecnologias

| Camada   | Tecnologia                           |
|----------|--------------------------------------|
| Backend  | Go (Golang), gorilla/websocket       |
| i18n     | go-i18n                             |
| Frontend | Web estático (HTML/JS/CSS, não incluso) |

## Estrutura do Projeto

```
/coup-game
├── cmd/                # Ponto de entrada principal
├── internal/
│   ├── game/           # Regras do jogo
│   ├── lobby/          # Salas e jogadores
│   ├── ws/             # WebSocket
│   └── i18n/           # Internacionalização (go-i18n)
│       └── locales/    # Arquivos de tradução JSON
├── pkg/                # Utilitários/Helpers compartilhados
├── web/                # Frontend (HTML, JS, CSS)
├── go.mod
├── LICENSE
├── README.md
└── README-ptbr.md
```

## Internacionalização (i18n)

- **Idiomas**: Inglês ou português, definido no início pelo anfitrião da sala
- **Traduções**: Nenhum texto exposto ao usuário é fixo; tudo vem dos arquivos JSON em `/internal/i18n/locales/`
- **Implementação**: Utiliza go-i18n para carregamento e distribuição das traduções
- **Código**: Código, comentários e documentação sempre em inglês

## Como Rodar

1. **Clone o repositório**:
   ```bash
   git clone https://github.com/leoferamos/coup-game.git
   cd coup-game
   ```

2. **Instale as dependências**:
   ```bash
   go mod tidy
   ```

3. **Rode o servidor**:
   ```bash
   go run ./cmd
   ```

4. **Acesse pelo navegador**:
   Acesse `http://localhost:8080` ou o IP da sua máquina pelo Wi-Fi nos celulares.

## Diretrizes de Desenvolvimento

- TODO o código e documentação deve estar em inglês
- Nenhuma lógica de negócio direto nos handlers HTTP/WebSocket
- Textos visíveis ao usuário nunca são hardcoded; use sempre i18n
- Commits claros, pull requests com testes e explicações
- Use as melhores práticas do ecossistema Go

## Licença

MIT

---

*Também disponível em [inglês](README.md)*

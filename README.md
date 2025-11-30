# Crispay Backend

Bem-vindo ao Crispay Backend — um serviço backend em Go que gerencia verificações KYC/AML utilizando a plataforma **Didit**. Este README descreve a estrutura do projeto, configuração, execução, endpoints, testes, migrações e instruções de desenvolvimento.

- Importante: Este backend oferece o KYC mas será integrado com diversas outras funcionalidades, se você é um dev e quer ter acesso ao projeto somente da verificação KYC/AML entre em contato. Obrigado!

---

## Sumário

- [Visão Geral](#visão-geral)
- [Principais Features](#principais-features)
- [Arquitetura do Projeto](#arquitetura-do-projeto)
- [Pré-requisitos](#pré-requisitos)
- [Variáveis de Ambiente](#variáveis-de-ambiente)
- [Execução Local (Sem Docker)](#execução-local-sem-docker)
- [Execução com Docker Compose](#execução-com-docker-compose)
- [Migrações de Banco de Dados](#migrações-de-banco-de-dados)
- [Geração de Código SQLC](#geração-de-código-sqlc)
- [Build & Deploy](#build--deploy)
- [Endpoints da API](#endpoints-da-api)
- [Testes](#testes)
- [Depuração e Troubleshooting](#depuração-e-troubleshooting)
- [Contribuição](#contribuição)
- [Licença](#licença)

---

## Visão Geral

O backend fornece endpoints para iniciar e gerenciar sessões de verificação KYC, integrando com a API da Didit para realizar o processo de verificação e recebendo status via webhooks. O projeto foi escrito em Go e usa:

- Gin como framework HTTP
- SQLC para geração de código seguro para acesso a banco Postgres
- Migrations SQL para evolução de schema
- HMAC para validar assinaturas de webhooks (segurança)

---

## Principais Features

- Início de verificação (Redirecionamento para Didit)
- Consulta de status de verificação
- Recebimento de webhooks para atualizar o estado da verificação
- Persistência com PostgreSQL
- Estrutura para testes e logs configuráveis

---

## Arquitetura do Projeto

Estrutura de diretórios principais:

- `cmd/`
	- `main.go` — ponto de entrada da aplicação
	- `api.go` — configuração das rotas e servidor
- `config/`
	- `config.go` — leitura do `.env` e carregamento de variáveis
	- `logger.go` — (se presente) abstração para log
- `internal/`
	- `database/` — setup do DB, migrations, utilidades sqlc
	- `didit/` — cliente Didit e lógica de integração
	- `handlers/` — handlers HTTP (ex: `webhook.go`)
	- `models/` — modelos de domínio
	- `repository/` — repository pattern com interfaces que usam código gerado por `sqlc`
- `vendor/` — módulos vendorizados (opcional, se `go mod vendor` usado)
- `docker-compose.yml` — ambiente de desenvolvimento com Postgres
- `internal/database/migrations` — scripts SQL de migração
- `makefile` — tarefas para criar e aplicar migrações
- `sqlc.yaml` — configuração do `sqlc`

---

## Pré-requisitos

- Go (recomendado >= 1.20) — verifique `go.mod`
- PostgreSQL (local ou Docker)
- `sqlc` (para gerar queries) — https://docs.sqlc.dev
- `migrate` (opcional, para executar migrações) — https://github.com/golang-migrate/migrate
- (Opcional) `docker` + `docker-compose` para subir um Postgres localmente

---

## Variáveis de Ambiente

Crie um arquivo `.env` na raiz com pelo menos as variáveis abaixo (exemplo):

```ini
# Didit
DIDIT_API_KEY=your_didit_api_key
DIDIT_WEBHOOK_SECRET_KEY=your_webhook_secret
DIDIT_WEBHOOK_URL=https://yourhost/api/webhooks/didit
DIDIT_WORKFLOW_ID=xxxxx
DIDIT_BASE_URL=https://api.didit.me

# Server & Database
SERVER_ADDR=6000
DATABASE_URL=postgresql://user:pass@localhost:5432/kyc_app?sslmode=disable

# Postgres (quando usar docker-compose)
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=kyc_app
POSTGRES_PORT=5432
POSTGRES_HOST=localhost
```

Observações:
- `SERVER_ADDR` é apenas a porta (por exemplo `6000`) ou já pode conter o `:` (ex.: `:6000`). A aplicação normaliza o formato ao iniciar.
- `DATABASE_URL` deve ser um DSN válido (Postgres) compatível com `sqlc` e bibliotecas `pgx`.

---

## Execução Local (Sem Docker)

1. Exporte as variáveis ou crie `.env`:

```bash
cp .env.example .env
# editar .env com seus valores
```

2. Instale dependências e gere vendor (opcional):

```bash
go mod download
# opcional
go mod vendor
```

3. Gere o código `sqlc` (se for necessário):

```bash
sqlc generate
```

4. Execute o servidor:

```bash
# usando a porta 6000 (padrão config)
export SERVER_ADDR=6001 # se quiser outra porta
go run cmd/*.go
```

---

## Execução com Docker Compose

O oferecemos um `docker-compose.yml` para subir o PostgreSQL localmente:

```bash
# Subindo Postgres
docker compose up -d

# verificar logs
docker compose logs -f postgres
```

Depois que o Postgres estiver pronto, configure a `DATABASE_URL` em `.env` com as credenciais do container (ou com host `localhost` e porta exposta `5432`).

---

## Migrações de Banco de Dados

As migrações estão em `internal/database/migrations`. Utilize o `makefile` para criar, aplicar e resetar as migrations:

```bash
# criar migração
make create_migration

# aplicar migrações
make migrate_up

# rollback 1
make migrate_down

# forçar estado
make migrate_force

# reset (drop + apply)
make reset
```

As strings de conexão usadas são montadas a partir das variáveis de ambiente no `Makefile`.

---

## Geração de Código SQLC

Este projeto usa `sqlc` para gerar tipos e queries tipadas em Go.

```bash
# gere o código (executar sempre que SQLs de queries forem alterados)
sqlc generate
```

Arquivos gerados ficam em `internal/database/sqlc/`.

---

## Build & Deploy

Para compilar a aplicação:

```bash
go build -o bin/crispay ./cmd
```

Para executar o binário:

```bash
./bin/crispay
```

Para containers e CI, crie uma imagem Docker (exemplo básico):

```Dockerfile
# Exemplo: Dockerfile
FROM golang:1.20-alpine as builder
WORKDIR /app
COPY . .
RUN go build -o digest ./cmd

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/digest /app/digest
EXPOSE 6000
CMD ["/app/digest"]
```

---

## Endpoints da API

- GET `/health` — status e timestamp (healthcheck)
- POST `/api/verification/start` — inicia verificação (redireciona para Didit)
- GET `/api/verification/status/{sessionId}` — status de sessão
- GET `/api/verification/user/{userId}` — verificação(s) do usuário
- POST `/api/webhooks/didit` — webhook que Didit usará para enviar resultados

Exemplo: iniciar verificação (curl, payload e headers dependem do contrato Didit):

```bash
curl -X POST http://localhost:6000/api/verification/start \
	-H "Content-Type: application/json" \
	-d '{"user_id": "user123"}'
```

Webhook expects HMAC header `X-Signature` (configurável) and body. O handler valida a assinatura com `DIDIT_WEBHOOK_SECRET_KEY`.

---

## Testes

Rode testes unitários/integrados:

```bash
go test ./...
```

Para testes com banco de dados, configure um Postgres de teste local e a variável `DATABASE_URL` apontando para ele.

---

## Depuração e Troubleshooting

- `address already in use`: outra aplicação já está rodando na porta configurada. Pare-a ou escolha outra porta:

```bash
export SERVER_ADDR=6001
go run cmd/*.go
```

- `too many colons in address`: A variável `SERVER_ADDR` veio com algo como `::6000` — a aplicação normaliza `SERVER_ADDR` mas se você quiser suportar IPv6 com host, utilize formato `[::1]:6000`.

- `.env not found`: a aplicação tentará ler `.env` com `godotenv`, mas não é obrigatório — se não existir, usa variáveis do ambiente do sistema.

---

## Logging

O projeto utiliza `slog`. Em `main.go` inicializamos um `slog.Logger` que é passado para a aplicação. Ajuste o `slog.New...` (text/json) conforme a necessidade.

---

## Security

- Webhook: HMAC-SHA256 com `DIDIT_WEBHOOK_SECRET_KEY`. Confira o handler `internal/handlers/webhook.go`.
- Cuide para não versionar secrets; use `./.env` e `.gitignore`.

---

## Contribuição

- Abra issues para bugs ou requests de features
- Faça PRs com testes
- Mantenha o padrão `gofmt` e `go vet`

Checklist para PRs:
- Código testado (`go test ./...`)
- Formatação (`gofmt`) e análise (`go vet`)
- Atualizar `sql`/`migrations` e `sqlc` se necessário

---

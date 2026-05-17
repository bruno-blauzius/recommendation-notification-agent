# project-go-sender-recommendation-agent

Serviço de recomendação de remetentes desenvolvido em **Go**, seguindo os princípios de **Arquitetura Hexagonal** (Ports & Adapters) e **Clean Architecture**.

---

## Sumário

- [Arquitetura](#arquitetura)
- [Estrutura de Diretórios](#estrutura-de-diretórios)
- [Pré-requisitos](#pré-requisitos)
- [Como executar](#como-executar)
- [Testes](#testes)
- [Variáveis de Ambiente](#variáveis-de-ambiente)
- [Banco de Dados](#banco-de-dados)
- [Migrations](#migrations)
- [RabbitMQ Consumer](#rabbitmq-consumer)

---

## Arquitetura

O projeto adota dois padrões arquiteturais complementares:

### Arquitetura Hexagonal (Ports & Adapters)

A aplicação é dividida em três zonas:

```
┌────────────────────────────────────────────────────────┐
│                    Adapters (Primary)                  │
│              HTTP handlers, CLI, gRPC, etc.            │
└───────────────────────┬────────────────────────────────┘
                        │  Input Ports (interfaces)
┌───────────────────────▼────────────────────────────────┐
│                    CORE (Domain)                        │
│         Entidades · Use Cases · Regras de Negócio       │
└───────────────────────┬────────────────────────────────┘
                        │  Output Ports (interfaces)
┌───────────────────────▼────────────────────────────────┐
│                   Adapters (Secondary)                  │
│          PostgreSQL, Redis, mensageria, APIs ext.       │
└────────────────────────────────────────────────────────┘
```

- **Ports de entrada (Input Ports):** interfaces que o núcleo expõe para ser chamado (ex.: `RecommendationUseCase`, `MessageHandler`).
- **Ports de saída (Output Ports):** interfaces que o núcleo define e que os adaptadores implementam (ex.: `RecommendationRepository`).
- **Adapters primários:** recebem requisições externas e chamam os use cases — inclui o consumer RabbitMQ.
- **Adapters secundários:** implementam os output ports — no caso, PostgreSQL.

### Clean Architecture

As dependências sempre apontam para dentro:

```
Infrastructure → Adapters → Core (Domain + Use Cases)
```

O núcleo (`core/`) não depende de nenhum framework, banco de dados ou detalhe técnico externo.

---

## Estrutura de Diretórios

```
project-go-sender-recommendation-agent/
├── cmd/
│   └── api/
│       └── main.go                  # Entrypoint da aplicação
│
├── internal/
│   ├── core/                        # Núcleo da aplicação (independente de frameworks)
│   │   ├── domain/
│   │   │   └── recommendation.go   # Entidade principal
│   │   ├── ports/
│   │   │   ├── input/
│   │   │   │   ├── message_handler.go          # Interface MessageHandler (driving port)
│   │   │   │   └── recommendation_usecase.go   # Interface RecommendationUseCase (driving port)
│   │   │   └── output/
│   │   │       └── recommendation_repository.go # Interface (driven port)
│   │   └── usecases/
│   │       └── recommendation_service.go       # Implementação dos use cases
│   │
│   ├── adapters/
│   │   ├── primary/                 # Adaptadores de entrada (HTTP, gRPC, CLI...)
│   │   │   └── rabbitmq/
│   │   │       ├── consumer.go             # Consumer RabbitMQ (adapter primário)
│   │   │       └── hello_world_handler.go  # Implementação de MessageHandler
│   │   └── secondary/
│   │       └── postgres/
│   │           └── recommendation_repository.go # Adaptador de saída (PostgreSQL)
│   │
│   └── infrastructure/              # Detalhes técnicos e configurações
│       ├── config/
│       │   └── config.go            # Leitura de variáveis de ambiente
│       └── database/
│           └── postgres.go          # Conexão com o PostgreSQL
│
├── migrations/                      # Scripts SQL de migração
│   ├── 000001_create_recommendations_table.up.sql
│   └── 000001_create_recommendations_table.down.sql
│
├── Dockerfile                       # Build multi-stage da imagem
├── docker-compose.yml               # Orquestração local (API + PostgreSQL + RabbitMQ)
├── .dockerignore
├── .gitignore
├── go.mod
└── README.md
```

---

## Pré-requisitos

| Ferramenta     | Versão mínima |
|----------------|---------------|
| Go             | 1.22          |
| Docker         | 24.x          |
| Docker Compose | 2.x           |

---

## Como executar

### Com Docker Compose (recomendado)

```bash
# Subir todos os serviços (API + PostgreSQL + RabbitMQ)
docker compose up --build

# Derrubar os serviços e remover volumes
docker compose down -v
```

### Localmente (sem Docker)

```bash
# Instalar dependências
go mod download

# Configurar as variáveis de ambiente
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=recommendation_db
export RABBITMQ_DSN=amqp://guest:guest@localhost:5672/
export RABBITMQ_QUEUE=hello_world

# Executar
go run ./cmd/api
```

---

## Testes

Os testes estão organizados em um diretório separado, divididos entre testes unitários e de integração. Nenhum teste requer infraestrutura real — bancos de dados são mockados.

```
tests/
├── unit/
│   ├── usecases/       # RecommendationService (mocks hand-written)
│   ├── rabbitmq/       # HelloWorldHandler
│   ├── redis/          # IdempotencyRepository (miniredis in-process)
│   └── postgres/       # RecommendationRepository (go-sqlmock)
└── integration/
    └── recommendation_service_integration_test.go
```

### Executar todos os testes

```bash
go test ./tests/...
```

### Executar apenas testes unitários

```bash
go test ./tests/unit/...
```

### Executar apenas testes de integração

```bash
go test ./tests/integration/...
```

### Executar com cobertura

```bash
go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Executar com output detalhado

```bash
go test ./tests/... -v
```

---

## Variáveis de Ambiente

| Variável      | Descrição                          | Padrão             |
|---------------|------------------------------------|--------------------|
| `DB_HOST`     | Host do banco de dados             | `localhost`        |
| `DB_PORT`     | Porta do banco de dados            | `5432`             |
| `DB_USER`     | Usuário do PostgreSQL              | `postgres`         |
| `DB_PASSWORD` | Senha do PostgreSQL                | `postgres`         |
| `DB_NAME`     | Nome do banco de dados             | `recommendation_db`|
| `DB_SSLMODE`  | Modo SSL da conexão                | `disable`          |
| `SERVER_PORT`    | Porta HTTP do servidor             | `8080`                              |
| `RABBITMQ_DSN`   | DSN de conexão com o RabbitMQ      | `amqp://guest:guest@localhost:5672/`|
| `RABBITMQ_QUEUE` | Nome da fila a ser consumida       | `hello_world`                       |

---

## Banco de Dados

### Tabela `recommendations`

| Coluna       | Tipo             | Descrição                          |
|--------------|------------------|------------------------------------|
| `id`         | `VARCHAR(36)`    | Identificador único (UUID)         |
| `sender_id`  | `VARCHAR(36)`    | ID do remetente                    |
| `payload`    | `TEXT`           | Conteúdo da recomendação           |
| `score`      | `NUMERIC(5,4)`   | Pontuação da recomendação (0–1)    |
| `created_at` | `TIMESTAMP`      | Data/hora de criação               |

### Exemplo de inserção

```sql
INSERT INTO recommendations (id, sender_id, payload, score, created_at)
VALUES (
    'a1b2c3d4-e5f6-7890-abcd-ef1234567890',
    'sender-uuid-001',
    '{"product_id": "prod-42", "reason": "frequently_bought_together"}',
    0.9500,
    NOW()
);
```

---

## Migrations

As migrations são aplicadas automaticamente pelo PostgreSQL ao iniciar via Docker Compose (pasta `migrations/` montada em `/docker-entrypoint-initdb.d`).

Para aplicar manualmente:

```bash
psql -h localhost -U postgres -d recommendation_db -f migrations/000001_create_recommendations_table.up.sql
```

Para reverter:

```bash
psql -h localhost -U postgres -d recommendation_db -f migrations/000001_create_recommendations_table.down.sql
```

---

## RabbitMQ Consumer

O serviço implementa um **consumer RabbitMQ** seguindo a Arquitetura Hexagonal:

```
RabbitMQ (broker)
      │
      ▼
consumer.go          ← adapter primário: conecta, declara fila, consome mensagens
      │  chama
      ▼
MessageHandler       ← input port (interface definida em core/ports/input)
      │  implementado por
      ▼
hello_world_handler.go  ← loga "Hello World — received message: <payload>"
```

### Biblioteca utilizada

[`github.com/rabbitmq/amqp091-go`](https://github.com/rabbitmq/amqp091-go) — cliente oficial RabbitMQ para Go (AMQP 0-9-1).

> **Documentação oficial:** [RabbitMQ Tutorials — Go](https://www.rabbitmq.com/tutorials/tutorial-one-go)

### Comportamento do consumer

| Situação | Ação |
|---|---|
| Mensagem processada com sucesso | `Ack` — mensagem removida da fila |
| Handler retorna erro | `Nack` com requeue — mensagem volta à fila |
| Payload vazio | Erro retornado pelo handler → Nack |

### Como publicar uma mensagem de teste

Acesse o **RabbitMQ Management UI** em `http://localhost:15672` (usuário `guest`, senha `guest`).

1. Acesse **Queues** → `hello_world`
2. Clique em **Publish message**
3. Preencha o campo **Payload** com qualquer texto
4. Clique em **Publish message**

O log da API exibirá:

```
Hello World — received message: <seu payload>
```

Ou via CLI com o container em execução:

```bash
docker exec recommendation-rabbitmq \
  rabbitmqadmin publish exchange=amq.default \
    routing_key=hello_world \
    payload="mensagem de teste"
```

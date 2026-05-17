# Engineering Instructions (Go + Hexagonal Architecture)

Você é um engenheiro de software sênior especializado em Go (Golang), com domínio em dois grupos de princípios:
- **Estrutura:** Clean Architecture e Arquitetura Hexagonal (Ports & Adapters)
- **Qualidade de código:** SOLID e Clean Code

Seu objetivo é sempre gerar código escalável, testável, desacoplado e de fácil manutenção.

> **Prioridade de aplicação dos princípios** (do mais para o menos prioritário):
> 1. Separação de responsabilidades e inward dependency rule (Clean Architecture)
> 2. Isolamento do domínio via Ports & Adapters (Arquitetura Hexagonal)
> 3. Inversão de dependência e interfaces segregadas (SOLID — D e I)
> 4. Demais princípios SOLID (S, O, L)
> 5. Diretrizes de Clean Code (nomenclatura, tamanho de funções, DRY)
>
> **Resolução de conflitos — aplique as perguntas na ordem:**
>
> | # | Pergunta | Resposta | Ação |
> |---|----------|----------|--------|
> | 1 | O princípio acoplaria `core/` a um detalhe externo (framework, banco, lib)? | Sim | **Recuse.** Mova a lógica para um adapter. |
> | 2 | O princípio quebraria a regra de dependência inward? | Sim | **Recuse.** Introduza uma interface (port) para inverter a dependência. |
> | 3 | O princípio de menor prioridade melhora legibilidade sem violar as regras acima? | Sim | **Aplique.** |
> | 4 | Nenhuma das condições acima se aplica? | — | Escolha o princípio com o número mais alto na lista de prioridades. |

---

## 🧠 Princípios obrigatórios

### Clean Architecture
- Separe claramente:
  - Domain (entidades e regras de negócio puras)
  - Application (casos de uso)
  - Infrastructure (DB, APIs externas)
  - Interfaces (handlers, controllers, CLI, HTTP)

- Dependências devem sempre apontar para dentro (inward dependency rule)

---

### Arquitetura Hexagonal (Ports and Adapters)
- Utilize Ports (interfaces) para comunicação entre camadas
- Adapters devem implementar Ports
- Nunca acople domínio a frameworks, banco ou detalhes externos

---

### SOLID
- S: Single Responsibility → funções e structs com responsabilidade única
- O: Open/Closed → código extensível sem modificação
- L: Liskov Substitution → interfaces respeitadas
- I: Interface Segregation → interfaces pequenas e específicas
- D: Dependency Inversion → dependa de abstrações, não implementações

---

### Clean Code
- Nomes claros e sem abreviações
- Funções pequenas (< 30 linhas)
- Evitar comentários redundantes ou que expliquem código óbvio, mas incluir comentários para lógica complexa ou decisões arquiteturais
- Tratamento de erro explícito (idiomático em Go)
- Evitar lógica duplicada (DRY)

---

## 🧱 Estrutura de projeto esperada

```bash
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
│   │   │   │   └── recommendation_usecase.go   # Interface (driving port)
│   │   │   └── output/
│   │   │       └── recommendation_repository.go # Interface (driven port)
│   │   └── usecases/
│   │       └── recommendation_service.go       # Implementação dos use cases
│   │
│   ├── adapters/
│   │   ├── primary/                 # Adaptadores de entrada (HTTP, gRPC, CLI...)
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
├── docker-compose.yml               # Orquestração local (API + PostgreSQL)
├── .dockerignore
├── .gitignore
├── go.mod
└── README.md
```

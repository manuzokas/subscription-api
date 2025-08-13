# 📦 API de Gerenciamento de Assinaturas Desenvolvida em GO

## 🔍 Visão Geral
Esta é uma API RESTful robusta para gerenciamento de assinaturas (Subscriptions), desenvolvida em Go (Golang). O projeto foi construído com foco em boas práticas de arquitetura de software, segurança e design de sistemas distribuídos, servindo como uma demonstração de habilidades de desenvolvimento backend para aplicações modernas, escaláveis e resilientes.

A API adota uma arquitetura assíncrona com um worker dedicado para processar tarefas em segundo plano, garantindo que a API principal permaneça rápida e responsiva.

---

## 🚀 Features Implementadas

- **Arquitetura Assíncrona com RabbitMQ**  
  Tarefas demoradas são desacopladas da API principal via mensageria.

- **Autenticação e Autorização com JWT**  
  Registro e login com emissão de tokens JWT.

- **Segurança de Senhas com bcrypt**  
  Senhas armazenadas de forma segura.

- **CRUD Completo para Assinaturas**  
  Gerenciamento completo do ciclo de vida.

- **Autorização Refinada**  
  Middleware garante acesso apenas aos recursos do usuário autenticado.

- **Validação de Entrada**  
  Validação robusta com feedback claro.

- **Configuração Externalizada (.env)**  
  Segredos e configurações sensíveis fora do código.

- **Persistência com PostgreSQL**  
  Banco relacional com consistência transacional.

---

## 🏛️ Design Arquitetural

O projeto segue os princípios da **Clean Architecture**, com separação clara de responsabilidades e suporte à mensageria para sistemas distribuídos.

cmd/
├── api/       → Ponto de entrada da API REST
└── worker/    → Ponto de entrada do Worker

internal/
├── domain/    → Entidades e regras de negócio
├── core/      → Casos de uso e interfaces
└── adapters/
    ├── web/       → Handlers HTTP
    ├── database/  → Acesso ao PostgreSQL
    └── messaging/ → Publicador de eventos RabbitMQ


### 🔄 Fluxo de Criação de Assinatura

1. API recebe `POST /subscriptions`
2. Valida requisição, salva no DB com status `PENDING`, responde `202 Accepted`
3. Publica evento `subscription.created` no RabbitMQ
4. Worker consome evento
5. Executa lógica (ex: ativa trial, envia e-mail), atualiza status para `TRIAL`

---

## 🛠️ Tecnologias Utilizadas

| Categoria            | Tecnologia             |
|----------------------|------------------------|
| Linguagem            | Go (Golang)            |
| Banco de Dados       | PostgreSQL             |
| Mensageria           | RabbitMQ               |
| Roteador HTTP        | Chi                    |
| Autenticação         | JWT for Go             |
| Driver PostgreSQL    | pgx                    |
| Driver RabbitMQ      | amqp091-go             |
| Validação            | go-playground/validator|
| Gerenciamento de Senhas | bcrypt              |
| Configuração         | godotenv               |
| Containerização      | Docker                 |

---

## ⚙️ Como Executar o Projeto Localmente

### ✅ Pré-requisitos

- Go (versão 1.20+)
- PostgreSQL
- Docker (para RabbitMQ)
- Git

### 📦 Passos para Configuração

# Clone o repositório
git clone https://github.com/seu-usuario/seu-repositorio.git
cd seu-repositorio

# Inicie o RabbitMQ via Docker
docker run -d --name meu-rabbit -p 5672:5672 -p 15672:15672 rabbitmq:3-management-alpine

A interface de gestão estará em: [http://localhost:15672](http://localhost:15672)  
Login: `guest` / Senha: `guest`

### 🧾 Crie o arquivo `.env`

DATABASE_URL="postgres://SEU_USUARIO:SUA_SENHA@localhost:5432/subscription_api"
JWT_SECRET="SEU_SEGREDO_SUPER_SEGURO_AQUI"
API_PORT="8080"
RABBITMQ_URL="amqp://guest:guest@localhost:5672/"

### 🗄️ Configure o Banco de Dados

CREATE DATABASE subscription_api;

CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE subscriptions (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL REFERENCES users(id),
    plan_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    cancelled_at TIMESTAMPTZ,
    trial_ends_at TIMESTAMPTZ
);

### 📦 Instale as dependências

go mod tidy

### ▶️ Execute a Aplicação


# Terminal 1 - API
go run ./cmd/api/main.go

# Terminal 2 - Worker
go run ./cmd/worker/main.go

---

## 📖 Documentação dos Endpoints

### 🔐 Autenticação

#### Registrar

- **POST** `/auth/register`


{
  "name": "Nome do Utilizador",
  "email": "utilizador@exemplo.com",
  "password": "umaPasswordForte"
}

#### Login

- **POST** `/auth/login`

{
  "email": "utilizador@exemplo.com",
  "password": "umaPasswordForte"
}

---

### 📬 Assinaturas

> Requer cabeçalho: `Authorization: Bearer <SEU_TOKEN_JWT>`

#### Criar Assinatura

- **POST** `/subscriptions`

{
  "planId": "plano_pro_mensal"
}

#### Buscar Assinatura

- **GET** `/subscriptions/{id}`

#### Cancelar Assinatura

- **DELETE** `/subscriptions/{id}`

---

## 🔮 Próximos Passos

- [x] Mensageria com RabbitMQ  
- [ ] Escrever Testes  
- [ ] Docker Compose  
- [ ] Logging Estruturado com zerolog  

---

> Projeto desenvolvido com foco em escalabilidade, segurança e boas práticas de arquitetura.

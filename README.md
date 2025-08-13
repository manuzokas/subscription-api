📦 API de Gerenciamento de Assinaturas Desenvolvida em GO.

🔍 Visão Geral
Esta é uma API RESTful robusta para gerenciamento de assinaturas (Subscriptions), desenvolvida em Go (Golang). O projeto foi construído com foco em boas práticas de arquitetura de software, segurança e design de sistemas distribuídos, servindo como uma demonstração de habilidades de desenvolvimento backend para aplicações modernas, escaláveis e resilientes.

A API adota uma arquitetura assíncrona com um worker dedicado para processar tarefas em segundo plano, garantindo que a API principal permaneça rápida e responsiva.

🚀 Features Implementadas
Arquitetura Assíncrona com RabbitMQ: Tarefas demoradas (como o processamento pós-criação de uma assinatura) são desacopladas da API principal. A API publica eventos numa fila e um worker independente consome e processa esses eventos, aumentando a resiliência e a performance percebida do sistema.

Autenticação e Autorização com JWT: Sistema completo de registro (/register) e login (/login) que emite JSON Web Tokens para autenticar requisições.

Segurança de Senhas: As senhas são seguramente "hasheadas" utilizando o algoritmo bcrypt.

CRUD Completo para Assinaturas: Gerenciamento completo do ciclo de vida das assinaturas.

Autorização Refinada: Um middleware garante que um usuário autenticado só possa visualizar ou modificar os recursos que lhe pertencem.

Validação de Entrada: Validação robusta dos dados de entrada com feedback claro para o cliente.

Configuração Externalizada: Uso de .env para segredos e configurações sensíveis.

Persistência com PostgreSQL: Banco de dados relacional com consistência transacional.

🏛️ Design Arquitetural
O projeto adota os princípios da Clean Architecture, com uma clara separação de responsabilidades. A introdução da mensageria expande a arquitetura para um modelo de sistema distribuído.

cmd/
├── api/       → Ponto de entrada da API REST (síncrona, rápida)
└── worker/    → Ponto de entrada do Worker (assíncrono, processa tarefas)

internal/
├── domain/    → Entidades e regras de negócio puras
├── core/      → Casos de uso e interfaces (contratos)
└── adapters/
    ├── web/       → Handlers e roteamento HTTP
    ├── database/  → Implementação do acesso ao PostgreSQL
    └── messaging/ → Implementação do publicador de eventos para RabbitMQ

Fluxo de Criação de Assinatura (Assíncrono):

API recebe POST /subscriptions.

API valida a requisição, salva a assinatura no DB com status PENDING e responde 202 Accepted imediatamente.

API publica um evento subscription.created na fila do RabbitMQ.

Worker consome o evento da fila.

Worker executa a lógica de negócio (ex: envia e-mail, ativa o trial) e atualiza o status da assinatura no DB para TRIAL.

🛠️ Tecnologias Utilizadas
Categoria

Tecnologia

Linguagem

Go (Golang)

Banco de Dados

PostgreSQL

Mensageria

RabbitMQ

Roteador HTTP

Chi

Autenticação

JWT for Go

Driver PostgreSQL

pgx

Driver RabbitMQ

amqp091-go

Validação

go-playground/validator

Gerenciamento de Senhas

bcrypt

Configuração

godotenv

Containerização

Docker

⚙️ Como Executar o Projeto Localmente
✅ Pré-requisitos
Go (versão 1.20+)

PostgreSQL

Docker (para o RabbitMQ)

Git

📦 Passos para Configuração
Clone o repositório:

git clone https://github.com/seu-usuario/seu-repositorio.git
cd seu-repositorio

Inicie o RabbitMQ via Docker:

docker run -d --name meu-rabbit -p 5672:5672 -p 15672:15672 rabbitmq:3-management-alpine

A interface de gestão fica disponível em http://localhost:15672 (login: guest/guest).

Crie o arquivo de configuração .env na raiz do projeto:

# Configurações do Banco de Dados
DATABASE_URL="postgres://SEU_USUARIO:SUA_SENHA@localhost:5432/subscription_api"

# Segredo para assinatura do JWT
JWT_SECRET="SEU_SEGREDO_SUPER_SEGURO_AQUI"

# Porta da API
API_PORT="8080"

# URL de conexão do RabbitMQ
RABBITMQ_URL="amqp://guest:guest@localhost:5672/"

Configure o Banco de Dados:
Conecte-se ao seu servidor PostgreSQL e execute os seguintes comandos SQL:

CREATE DATABASE subscription_api;

-- Após conectar ao banco:
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

Instale as dependências:

go mod tidy

Execute a Aplicação (API e Worker):
Abra dois terminais separados na raiz do projeto.

No Terminal 1, inicie a API:

go run ./cmd/api/main.go

No Terminal 2, inicie o Worker:

go run ./cmd/worker/main.go

📖 Documentação dos Endpoints
🔐 Autenticação (Endpoints Públicos)
Registar um Novo Utilizador

Método: POST

Endpoint: /auth/register

Descrição: Cria uma nova conta de utilizador no sistema.

Corpo da Requisição (application/json):

{
    "name": "Nome do Utilizador",
    "email": "utilizador@exemplo.com",
    "password": "umaPasswordForte"
}

Resposta de Sucesso (201 Created):

{
    "id": "...",
    "name": "Nome do Utilizador",
    "email": "utilizador@exemplo.com",
    "createdAt": "...",
    "updatedAt": "..."
}

Autenticar um Utilizador

Método: POST

Endpoint: /auth/login

Descrição: Verifica as credenciais de um utilizador e retorna um token JWT para ser usado em rotas protegidas.

Corpo da Requisição (application/json):

{
    "email": "utilizador@exemplo.com",
    "password": "umaPasswordForte"
}

Resposta de Sucesso (200 OK):

{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}

📬 Assinaturas (Endpoints Protegidos)
Requisito: Todas as requisições para estes endpoints devem incluir o cabeçalho de autorização.
Authorization: Bearer <SEU_TOKEN_JWT>

Criar uma Nova Assinatura

Método: POST

Endpoint: /subscriptions

Descrição: Inicia o processo de criação de uma nova assinatura para o utilizador autenticado. A requisição é processada de forma assíncrona. A API responde imediatamente e um worker processa a ativação em segundo plano.

Corpo da Requisição (application/json):

{
    "planId": "plano_pro_mensal"
}

Resposta de Sucesso (202 Accepted):

{
    "id": "...",
    "userId": "...",
    "planId": "plano_pro_mensal",
    "status": "PENDING",
    "createdAt": "...",
    "updatedAt": "..."
}

Buscar Detalhes de uma Assinatura

Método: GET

Endpoint: /subscriptions/{id}

Descrição: Retorna os detalhes de uma assinatura específica que pertença ao utilizador autenticado.

Parâmetros da URL:

id (string): O ID da assinatura a ser buscada.

Resposta de Sucesso (200 OK):

{
    "id": "...",
    "userId": "...",
    "planId": "plano_pro_mensal",
    "status": "TRIAL",
    "createdAt": "...",
    "updatedAt": "...",
    "trialEndsAt": "..."
}

Cancelar uma Assinatura

Método: DELETE

Endpoint: /subscriptions/{id}

Descrição: Cancela uma assinatura ativa que pertença ao utilizador autenticado. Esta operação é idempotente.

Parâmetros da URL:

id (string): O ID da assinatura a ser cancelada.

Resposta de Sucesso (204 No Content):

Nenhum corpo na resposta.

🔮 Próximos Passos
[x] Mensageria com RabbitMQ

[ ] Escrever Testes: Adicionar testes unitários e de integração.

[ ] Containerização com Docker Compose: Criar um docker-compose.yml para orquestrar a API, o DB, e o RabbitMQ com um único comando.

[ ] Logging Estruturado: Implementar um logger mais robusto (ex: zerolog) para logs em formato JSON.

Projeto desenvolvido com foco em escalabilidade, segurança e boas práticas de arquitetura.
# 📦 API de Gerenciamento de Assinaturas

## 🔍 Visão Geral
Esta é uma API RESTful robusta para gerenciamento de assinaturas (Subscriptions), desenvolvida em **Go (Golang)**. O projeto foi construído com foco em boas práticas de arquitetura de software, segurança e código limpo, servindo como uma demonstração de habilidades de desenvolvimento backend para aplicações modernas e escaláveis.

A API permite:
- Registro e autenticação de usuários
- CRUD completo de assinaturas
- Sistema de permissões que garante acesso apenas aos próprios recursos

---

## 🚀 Features Implementadas

- **Autenticação e Autorização com JWT**  
  Registro (`/register`) e login (`/login`) com emissão de JSON Web Tokens

- **Segurança de Senhas**  
  Senhas hasheadas com `bcrypt`, nunca armazenadas em texto plano

- **CRUD Completo para Assinaturas**
  - `POST /subscriptions`: Cria uma nova assinatura
  - `GET /subscriptions/{id}`: Busca detalhes de uma assinatura
  - `DELETE /subscriptions/{id}`: Cancela uma assinatura

- **Autorização Refinada**  
  Middleware garante acesso apenas aos recursos do próprio usuário

- **Validação de Entrada**  
  Validação robusta com feedback claro para o cliente

- **Configuração Externalizada**  
  Uso de `.env` para segredos e configurações sensíveis

- **Persistência com PostgreSQL**  
  Banco relacional com consistência transacional

---

## 🏛️ Design Arquitetural

Adota os princípios da **Clean Architecture**, com camadas bem definidas:

```
domain/      → Entidades e regras de negócio puras
core/        → Casos de uso e interfaces de repositórios
adapters/
├── web/     → Handlers e roteamento HTTP
└── database/→ Implementação do acesso ao PostgreSQL
```

**Regra de Dependência:** todas as dependências apontam para dentro (em direção ao `domain`), tornando o núcleo independente de frameworks e tecnologias.

---

## 🛠️ Tecnologias Utilizadas

| Categoria             | Tecnologia               |
|----------------------|--------------------------|
| Linguagem            | Go (Golang)              |
| Banco de Dados       | PostgreSQL               |
| Roteador HTTP        | Chi                      |
| Autenticação         | JWT for Go               |
| Driver PostgreSQL    | pgx                      |
| Validação            | go-playground/validator  |
| Gerenciamento de Senhas | bcrypt               |
| Configuração         | godotenv                 |

---

## ⚙️ Como Executar o Projeto Localmente

### ✅ Pré-requisitos

- Go (versão 1.20+)
- PostgreSQL
- Git

### 📦 Passos para Configuração

```bash
git clone https://github.com/seu-usuario/seu-repositorio.git
cd seu-repositorio
```

Crie o arquivo `.env` na raiz do projeto:

```env
DATABASE_URL="postgres://SEU_USUARIO:SUA_SENHA@localhost:5432/subscription_api"
JWT_SECRET="SEU_SEGREDO_SUPER_SEGURO_AQUI"
API_PORT="8080"
```

Configure o banco de dados:

```sql
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
```

Instale as dependências:

```bash
go mod tidy
```

Execute a API:

```bash
go run ./cmd/api/main.go
```

A API estará disponível em: `http://localhost:8080`

---

## 📖 Documentação dos Endpoints

### 🔐 Autenticação

| Método | Endpoint        | Descrição                  |
|--------|------------------|----------------------------|
| POST   | `/auth/register` | Registra um novo usuário   |
| POST   | `/auth/login`    | Autentica e retorna JWT    |

**Exemplo de corpo para `/auth/register`:**

```json
{
  "name": "Nome do Utilizador",
  "email": "utilizador@exemplo.com",
  "password": "umaPasswordForte"
}
```

---

### 📬 Assinaturas (Requer Token JWT)

**Cabeçalho obrigatório:**

```
Authorization: Bearer <SEU_TOKEN_JWT>
```

| Método | Endpoint              | Descrição                          |
|--------|------------------------|------------------------------------|
| POST   | `/subscriptions`       | Cria nova assinatura               |
| GET    | `/subscriptions/{id}`  | Busca detalhes da assinatura       |
| DELETE | `/subscriptions/{id}`  | Cancela assinatura específica      |

**Exemplo de corpo para `POST /subscriptions`:**

```json
{
  "planId": "plano_pro_mensal"
}
```

---

## 🔮 Próximos Passos

- [ ] **Mensageria com RabbitMQ**  
  Desacoplar tarefas demoradas com workers assíncronos

- [ ] **Escrever Testes**  
  Testes unitários e de integração

- [ ] **Containerização com Docker**  
  `docker-compose.yml` para orquestrar API, DB e RabbitMQ

---

> Projeto desenvolvido com foco em escalabilidade, segurança e boas práticas de arquitetura. Ideal para quem busca aprender ou demonstrar habilidades avançadas em backend com Go.

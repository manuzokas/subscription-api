Guia de Teste da API de Assinaturas

Este documento descreve o fluxo completo para testar as funcionalidades da API, incluindo o sistema de autenticação e o processamento assíncrono de assinaturas.

Pré-requisito: A API, o Worker e o RabbitMQ devem estar a correr conforme as instruções no README.md.

Ferramenta: Os comandos abaixo são para serem executados no PowerShell.

Fluxo de Teste Completo
1. Registar um Novo Utilizador
Vamos criar uma conta. Este endpoint é público.

Invoke-RestMethod -Method Post -Uri http://localhost:8080/auth/register -ContentType "application/json" -Body '{"name": "Utilizador de Teste", "email": "teste@exemplo.com", "password": "password123"}'

Resultado Esperado: Um objeto JSON com os detalhes do utilizador recém-criado.

2. Tentar Aceder a uma Rota Protegida (Sem Token)
Este passo verifica se o nosso middleware de segurança está a funcionar.

Invoke-RestMethod -Method Post -Uri http://localhost:8080/subscriptions -ContentType "application/json" -Body '{"planId": "plano-teste"}' -SkipHttpErrorCheck

Resultado Esperado: Um erro 401 Unauthorized com a mensagem "Authorization header required".

3. Fazer Login e Capturar o Token JWT
Agora, vamos autenticar e guardar o token para as próximas requisições.

# Executa o login e guarda a resposta completa
$loginResponse = Invoke-RestMethod -Method Post -Uri http://localhost:8080/auth/login -ContentType "application/json" -Body '{"email": "teste@exemplo.com", "password": "password123"}'

# Extrai o token da resposta
$token = $loginResponse.token

# Cria o cabeçalho de autorização para usar nas próximas chamadas
$headers = @{ "Authorization" = "Bearer $token" }

# (Opcional) Exibe o token para verificar
Write-Host "Token capturado: $token"

4. Criar uma Assinatura (Processamento Assíncrono)
Com o token em mãos, vamos criar uma assinatura.

# Guarda o resultado da criação para podermos usar o ID depois
$newSubscription = Invoke-RestMethod -Method Post -Uri http://localhost:8080/subscriptions -Headers $headers -ContentType "application/json" -Body '{"planId": "plano_pro_async"}'

# Exibe a resposta imediata da API
$newSubscription

Observações:

A resposta da API deve ser quase instantânea.

O status da assinatura retornada será PENDING.

No terminal onde o Worker está a correr, você verá logs a aparecer, indicando que uma mensagem foi recebida e está a ser processada.

5. Verificar o Processamento do Worker
Aguarde alguns segundos (o nosso worker tem um sleep de 3 segundos) e depois verifique o estado final da assinatura.

# Usa o ID da assinatura que guardámos no passo anterior
$subscriptionId = $newSubscription.id

# Faz a requisição para obter o estado atualizado
Invoke-RestMethod -Method Get -Uri "http://localhost:8080/subscriptions/$subscriptionId" -Headers $headers

Resultado Esperado: O status da assinatura agora deve ser TRIAL, e o campo trialEndsAt deve estar preenchido. Isto confirma que o worker processou o evento com sucesso.

6. Cancelar a Assinatura
Para finalizar, vamos testar o endpoint de cancelamento.

Invoke-RestMethod -Method Delete -Uri "http://localhost:8080/subscriptions/$subscriptionId" -Headers $headers

Resultado Esperado: Nenhuma saída no terminal (resposta 204 No Content).

Verifique o estado final:

Invoke-RestMethod -Method Get -Uri "http://localhost:8080/subscriptions/$subscriptionId" -Headers $headers

Resultado Esperado: O status da assinatura agora é CANCELLED, e o campo cancelledAt está preenchido.
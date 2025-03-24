#####################################################################################################################################
COMO ESTA IMPLEMENTAÇÃO DE RATE LIMITER FUNCIONA?
#####################################################################################################################################

Este rate limiter funciona baseados em 2 informações-alvos:
- Token
- IP

É feita a validação da requisição primeiramente pelo Token, caso essa informação exista. Caso não exista Token na requisição, a validação é feita baseada no IP. Com isso, se há Token na requisição, é feita sua validação e a validação por IP é descartada.

É possível ligar e desligar as validações de Token e IP pelo arquivo de configuração .env, na pasta raiz do projeto.
As chaves para esse chaveamento são os abaixo:

RATELIMITER_IP_ON = true | false
RATELIMITER_TOKEN_ON = true | false

#####################################################################################################################################
VALIDAÇÃO POR IP
#####################################################################################################################################

É possível configurar a quantidade de req/s que o rate limiter vai utilizar para a validação por IP.

A chave RATELIMITER_IP_MAX_REQUESTS indica a quantidade de req/s permitidas por segundo.
A chave RATELIMITER_IP_BLOCKED_SECONDS indica quantos segundos o IP deve ficar bloqueado para chamadas após ser atingida a quantidade máxima de requisições, conforme o parâmetro RATELIMITER_IP_MAX_REQUESTS.

#####################################################################################################################################
VALIDAÇÃO POR TOKEN
#####################################################################################################################################

Cada Token gerado, tem seus parâmetros de rate limiter em suas claims. Para isso, o token deve possuir 2 chaves:

rl-max-requests: indica a quantidade de req/s permitidas por segundo.
rl-seconds-blocked: indica quantos segundos que o token deve ficar bloqueado para chamadas após ser atingida a quantidade máxima de requisições, conforme a chave rl-max-requests.

Para obter um token, utilizar o seguinte request:

POST http://localhost:8080/login
Content-Type: application/json
Accept: application/json

{
    "email": "rene@test.com",
    "maxRequests" : 2,
    "blockedSeconds" : 60
}

Há um arquivo test.http na raiz do projeto com esta requisição.
 
#####################################################################################################################################
PARÂMETROS DEFAULT
#####################################################################################################################################

Caso algum parâmetro seja inválido, será assumido os valores default, encontrados no arquivo de .env.
 
Ex:
RATELIMITER_TOKEN_DEFAULT_MAX_REQUESTS=3
RATELIMITER_TOKEN_DEFAULT_BLOCKED_SECONDS=10

#####################################################################################################################################
REQUISIÇÕES DA API
#####################################################################################################################################

Desenvolvi uma rota que devolve uma lista de músicas, que tem o rate limiter "plugado" em seu pipelinte:

GET http://localhost:8080/songs
API_KEY: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDI1MjUxOTksImtleSI6ImFkbWluIiwicmwtbWF4LXJlcXVlc3RzIjoyLCJybC1zZWNvbmRzLWJsb2NrZWQiOjYwfQ.t8TRma0IHZVEdt8XaoXH2bqkK3oAQ8Ab7imFygfe5dE

Há um arquivo test.http na raiz do projeto com esta requisição.

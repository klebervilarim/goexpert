Aqui está um README.md completo, pronto para copiar e colar no seu projeto:

# Desafio: Cotação Dólar (Go)

## Descrição
Este desafio consiste em dois programas em Go:

- **server.go**: expõe o endpoint `/cotacao` na porta `8080` (ou `8081` se 8080 estiver ocupada localmente).  
  Ele consulta a API `https://economia.awesomeapi.com.br/json/last/USD-BRL` com timeout de 200ms, persiste a cotação no SQLite com timeout de 10ms e retorna JSON contendo o campo `bid`.

- **client.go**: conecta ao servidor (timeout de 300ms), recebe o valor do `bid` e salva em `cotacao.txt` no formato:



Dólar: {valor}


Logs informam timeouts ou outros erros caso ocorram.

---

## Arquivos
- `server.go`
- `client.go`
- `go.mod`
- `Dockerfile`
- `docker-compose.yml` (opcional)

---

## Requisitos
- Go 1.21+
- Docker (opcional)
- SQLite3 (opcional, para inspeção local)

---

## Execução Local (sem Docker)

1. Abra o terminal e entre na pasta do projeto:

```bash
cd ~/projetos/MBA-Golang/goexpert/desafio/client-server-api


Inicialize o módulo Go (se ainda não tiver):

go mod init cotacao


Baixe as dependências:

go get modernc.org/sqlite
go mod tidy


Compile o servidor:

go build -o cotacao-server server.go


Execute o servidor:

export RUN_LOCAL=1   # informa que é execução local
./cotacao-server     # Linux/macOS/Git Bash
cotacao-server.exe   # Windows


O servidor subirá em localhost:8080 ou localhost:8081 (fallback automático se 8080 estiver ocupada).

Teste o endpoint:

curl http://localhost:8080/cotacao   # ou :8081 se estiver usando fallback


Deve retornar JSON:

{"bid":"5.3461"}


Execute o cliente:

go run client.go


Criará o arquivo cotacao.txt com o valor do dólar:

Dólar: 5.3461


Verifique o conteúdo do arquivo:

cat cotacao.txt     # Linux/macOS/Git Bash
type cotacao.txt    # Windows

Execução via Docker

Build e execução dos containers:

docker-compose up --build


Servidor disponível em localhost:8080.

SQLite persistido no volume cotacoes-data.

Execute o cliente localmente:

go run client.go


Cria cotacao.txt com o valor do dólar.

Para parar os containers:

docker-compose down

Observações

Para mudar a porta local sem conflito, use a variável RUN_LOCAL=1.

Logs do servidor mostram timeouts de API ou de gravação no SQLite.

O cliente possui timeout de 300ms; se o servidor estiver lento, ele reportará erro.

O arquivo cotacao.txt é sobrescrito a cada execução do cliente.




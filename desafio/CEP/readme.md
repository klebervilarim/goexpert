# Desafio: CEP - API mais rápida (Go)

## Descrição
Este desafio consiste em criar um programa em Go que consulta simultaneamente duas APIs de CEP:

1. `https://brasilapi.com.br/api/cep/v1/{cep}`  
2. `http://viacep.com.br/ws/{cep}/json/`  

O objetivo é:

- Executar as duas requisições ao mesmo tempo.  
- Exibir no terminal a **primeira resposta que chegar**, informando qual API enviou os dados.  
- Descartar a resposta mais lenta.  
- Limitar o tempo total de execução a **1 segundo**.  
- Exibir erro caso nenhuma API responda dentro do timeout.

> Observação: A estrutura do JSON retornada pela BrasilAPI mudou. O programa já mapeia corretamente os campos `street → logradouro`, `neighborhood → bairro`, `city → cidade` e `state → uf`.

---

## Arquivos
- `main.go` : código principal do desafio  
- `go.mod`  : módulo Go para gerenciamento de dependências  

---

## Requisitos
- Go 1.21+  
- Conexão com internet para acessar as APIs  

---

## Execução Local

1. Abra o terminal e entre na pasta do projeto:

```bash
cd /MBA-Golang/goexpert/desafio/CEP

go mod init cep-fastest

go mod tidy

go build -o cep-main main.go

./cep-main -cep=01153000   # Linux/macOS/Git Bash
cep-main.exe -cep=01153000 # Windows

$  ./cep-main -cep=01153000

Resposta recebida da API: BrasilAPI
Logradouro: Rua Vitorino Carmilo
Bairro: Barra Funda
Cidade: São Paulo
UF: SP


# Etapa de build
FROM golang:tip-20250209-alpine3.22
WORKDIR /apphttps://hub.docker.com/layers/library/golang/1.24-alpine/images/sha256-3e347e3f7694d3a1db144f437bcd2674ad41726be54a3526538d476c7f0e1ce3

# Copia os arquivos de módulo e baixa as dependências para aproveitar o cache do Docker
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do código-fonte
COPY . .

# Instala ferramentas, constrói e otimiza o binário
RUN apk add --no-cache ca-certificates upx && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o cloudrun . && \
    upx --best --lzma cloudrun

# Etapa final
FROM gcr.io/distroless/static-debian12
WORKDIR /app

# Copia os certificados de autoridade da etapa de build para permitir chamadas HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copia o binário da aplicação
COPY --from=builder /app/cloudrun .
USER nonroot:nonroot
ENTRYPOINT ["/app/cloudrun"]
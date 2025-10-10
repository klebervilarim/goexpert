FROM golang:1.24 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/api ./cmd/api

FROM gcr.io/distroless/base-debian12
ENV PORT=8083
EXPOSE 8083
COPY --from=build /app/bin/api /api
USER nonroot:nonroot
ENTRYPOINT ["/api"]

# build
FROM golang:1.23-alpine AS build
WORKDIR /src

# Copia manifests, baixa deps (cache)
COPY go.mod ./
# COPY go.sum ./
RUN go mod download

# Copia código e garante go.sum
COPY . .
RUN go mod tidy

# Build estático
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -o /out/app main.go

# run (distroless)
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /out/app /app/app
EXPOSE 3000
USER nonroot:nonroot
ENTRYPOINT ["/app/app"]

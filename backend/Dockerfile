FROM golang:1.22-alpine AS builder
# go mod download использует GOMODCACHE, чтобы кэшировать модули
WORKDIR /build
# Сначала — только зависимости: слой кэшируется при изменении кода
COPY go.mod go.sum ./
RUN go mod download
# Теперь — исходники и сборка статического бинарника
COPY . .
# -ldflags="-s -w" убирает debug-информацию из бинарника (образ меньше)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/app ./main.go
# ---------- Этап 2: исполнение ----------
# distroless/static — без shell, без пакетного менеджера:
# атакующий surface минимален. Недостаток: нельзя docker exec bash.
FROM gcr.io/distroless/static-debian12:nonroot
# nonroot-вариант distroless запускается под пользователем nobody/nogroup
WORKDIR /app
COPY --from=builder /out/app /app/app

EXPOSE 8080
# ENTRYPOINT — бинарник, CMD — аргументы по умолчанию
ENTRYPOINT ["/app/app"]
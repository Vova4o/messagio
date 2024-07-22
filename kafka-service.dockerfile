# FROM golang:1.22

# WORKDIR /app

# COPY kafkaApp /app/
# COPY ./data/sql /app/data/sql
# COPY ./data/templates /app/data/templates

# RUN chmod +x /app/kafkaApp

# CMD ["/app/kafkaApp"]

# Использование multi-stage build
FROM golang:1.22 as builder
# FROM ubuntu:20.04 as builder

# # Установка зависимостей
# RUN apk --no-cache update && \
# 	apk --no-cache add git gcc libc-dev && \
# 	rm -rf /var/cache/apk/*

# Установка переменных окружения
# ENV CGO_ENABLED=1 \
	# GOFLAGS=-mod=vendor \
	# GOOS=linux \
	# GOARCH=amd64 \
	# GO111MODULE=on

# Копирование исходного кода и сборка приложения
WORKDIR /app
COPY . .
# RUN go build -tags musl -ldflags "-s -w" -o kafApp ./cmd/
RUN go build -o kafkaApp ./cmd/

# Финальный образ
FROM golang:1.22
COPY --from=builder /app/kafkaApp ./
# Копирование директорий sql и templates из контекста сборки в финальный образ
COPY ./data/sql /data/sql
COPY ./data/templates /data/templates
EXPOSE 8080
CMD ["./kafkaApp"]

# # Финальный образ
# FROM alpine:latest
# COPY --from=builder /app/kafApp /app/
# # Копирование директорий sql и templates из контекста сборки в финальный образ
# COPY ./data/sql /app/data/sql
# COPY ./data/templates /app/data/templates
# CMD ["/app/kafApp"]
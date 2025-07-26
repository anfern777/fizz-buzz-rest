## Dev
FROM golang:1.24.3-alpine as development

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/air-verse/air@latest

COPY . .

EXPOSE 8000

CMD air

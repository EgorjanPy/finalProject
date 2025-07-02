FROM golang:1.23

LABEL authors="parko"

WORKDIR /app
COPY . .

RUN CGO_ENABLED=1

RUN go build -o server ./cmd/orchestrator/main.go

RUN go build -o agent ./cmd/agent/main.go

EXPOSE 8080

CMD ./server & ./agent

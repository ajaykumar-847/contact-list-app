FROM golang:1.22.4 AS builder

COPY go.mod /app/
COPY go.sum /app/
COPY main.go /app/
COPY README.md /app/
COPY .env /app/
COPY static /app/static/
COPY templates /app/templates/

WORKDIR /app
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
RUN go mod download

EXPOSE 8000

CMD ["go", "run", "main.go"]

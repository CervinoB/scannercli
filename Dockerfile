FROM golang:1.20-alpine

WORKDIR /app

COPY main.go ./
RUN go run main.go

# COPY . .

# RUN go build -o exporter .

# CMD ["./exporter"]

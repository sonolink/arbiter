FROM golang:1.27.0-alpine

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["sleep", "infinity"]

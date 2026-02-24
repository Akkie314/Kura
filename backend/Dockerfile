FROM golang:1.25-bookworm

WORKDIR /app

RUN apt-get update && apt-get install -y ffmpeg

COPY go.mod go.sum ./
RUN go mod download

CMD ["go","run","main.go"]
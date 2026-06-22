FROM golang:1.26.4 AS builder
LABEL maintainer="tai@taigrr.com"
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o ssh-wars main.go

FROM scratch
WORKDIR /app
COPY --from=builder /src/ssh-wars .
EXPOSE 2222
CMD ["/app/ssh-wars"]

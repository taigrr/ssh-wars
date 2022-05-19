FROM golang:1.18 AS builder
MAINTAINER tai@taigrr.com
RUN mkdir -p /src
COPY starwars.ascii /src
COPY *.go /src/
COPY go.mod /src
COPY go.sum /src
WORKDIR /src
RUN go mod tidy
RUN CGO_ENABLED=0 go build ./...
RUN mv ssh-wars main

FROM scratch
WORKDIR /app
COPY --from=builder /src/main .
CMD ["/app/main"]

FROM golang:1.18 AS builder
MAINTAINER tai@taigrr.com
RUN mkdir -p /src
ADD . /src/
WORKDIR /src
RUN go mod tidy
RUN CGO_ENABLED=0 go build main.go

FROM scratch
WORKDIR /app
COPY --from=builder /src/main .
CMD ["/app/main"]

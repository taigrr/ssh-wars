FROM golang:1.26 AS builder
LABEL maintainer="tai@taigrr.com"
RUN mkdir -p /src
ADD . /src/
WORKDIR /src
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o ssh-wars main.go

FROM scratch
WORKDIR /app
COPY --from=builder /src/ssh-wars .
EXPOSE 2222
CMD ["/app/ssh-wars"]

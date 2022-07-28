FROM golang:1.17-alpine as builder
RUN apk --no-cache add ca-certificates
WORKDIR /docker-try-fiber
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /main ./cmd/main.go
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /main /main
COPY emails.json .
EXPOSE 8000
ENTRYPOINT ["./main"]
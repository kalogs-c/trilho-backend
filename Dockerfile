FROM golang:1.19.4-alpine as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o trilho .

RUN apk add -U --no-cache ca-certificates

FROM scratch

COPY --from=builder /app/trilho .

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ 

EXPOSE 80

ENTRYPOINT [ "./trilho" ]

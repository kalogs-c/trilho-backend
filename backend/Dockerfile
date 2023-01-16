FROM golang:1.19.4-alpine as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o ego .

FROM scratch

COPY --from=builder /app/ego .

ENTRYPOINT [ "./ego" ]
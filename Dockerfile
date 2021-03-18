FROM golang:1.16-alpine AS builder

WORKDIR /api

COPY . .

RUN CGO_ENABLED=0 go build -o server

FROM scratch

COPY --from=builder /api/server /api/server

ENTRYPOINT [ "/api/server" ]
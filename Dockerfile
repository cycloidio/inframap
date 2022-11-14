FROM golang:1.16.5-alpine3.15 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN apk -q --no-progress add git make; \
	make build

FROM alpine:3.15
COPY --from=builder /app/inframap /app/
ENTRYPOINT ["/app/inframap"]

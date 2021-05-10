FROM golang:1.16.4-alpine3.12 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN apk -q --no-progress add git make; \
	make build

FROM alpine
COPY --from=builder /app/inframap /app/
ENTRYPOINT ["/app/inframap"]

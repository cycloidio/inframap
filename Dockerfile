FROM golang:1.16.5-alpine3.12 as builder
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN apk -q --no-progress add git make \
	&& make build

FROM alpine

RUN apk -q --no-progress add graphviz ttf-dejavu \
        && addgroup -g 1000 inframap \
        && adduser -u 1000 -G inframap -s /bin/ash -D inframap

USER 1000

WORKDIR /home/inframap

COPY --from=builder /app/inframap /home/inframap

ENTRYPOINT ["./inframap"]

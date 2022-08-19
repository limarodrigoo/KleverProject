FROM golang:1.18-alpine AS stage

WORKDIR /app

ENV GO111MODULE=on

COPY . .

RUN go mod download

RUN go build -x server/main.go

FROM alpine

COPY --from=stage /app/main /app/

WORKDIR /app

CMD [ "./main" ]

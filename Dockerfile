FROM golang:1.22.1-alpine

WORKDIR /

ARG override_tables=false
ARG port=8888
ARG default_admin=false

ENV OVERRIDE=$override_tables
ENV PORT=$port
ENV DEF_ADMIN=$default_admin

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o out ./cmd/api

EXPOSE ${port}

CMD [ "sh", "-c", "./out -override_tables=$OVERRIDE -addr=:$PORT -default_admin=$DEF_ADMIN" ]
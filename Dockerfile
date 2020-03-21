FROM golang:1.13 as builder
WORKDIR /app
COPY ./src/go.mod ./src/go.sum ./
RUN go mod download
COPY ./src .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o keybase-bot .

FROM keybaseio/client
ARG KEYBASE_SERVICE=1
ARG KEYBASE_USERNAME
ARG KEYBASE_PAPERKEY
ENV KEYBASE_SERVICE=1
ENV KEYBASE_USERNAME=${KEYBASE_USERNAME}
ENV KEYBASE_PAPERKEY=${KEYBASE_PAPERKEY}
COPY --from=builder /app/keybase-bot /app/.env ./
CMD ["./keybase-bot"]
FROM golang:1.21 AS build

WORKDIR /password_manager

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/pm ./cmd/server/

FROM alpine:3.18 AS release

COPY --from=build /password_manager/bin/pm /bin

ENTRYPOINT ["/bin/pm"]

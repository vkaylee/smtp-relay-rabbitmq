# Build app stage
FROM ghcr.io/vleedev/cont-images:golang-1.15.6-alpine3.12 AS build
WORKDIR /workdir
COPY ./ ./
RUN go get -d -v ./...
RUN go install -v ./...
RUN env GOOS=linux GOARCH=amd64 go build -o /myapp main.go
# "app" stage
FROM ghcr.io/vleedev/cont-images:alpine-3.12.3 AS app

WORKDIR /srv/app

COPY --from=build /myapp ./

RUN chmod +x myapp

CMD ["./myapp"]
# syntax=docker/dockerfile:1
ARG GO_VERSION=1.18.3

## BUILD
FROM golang:${GO_VERSION}-alpine AS build
WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .
RUN echo "nobody:x:65534:65534:Nobody:/:" >> /etc/passwd

## DEPLOY
FROM scratch AS final
WORKDIR /bin/
COPY --from=build /src/app .
COPY --from=build /src/config.yaml .
COPY --from=0 /etc/passwd /etc/passwd
USER nobody
EXPOSE 8080
ENV GIN_RELEASE_MODE=release
ENTRYPOINT ["./app"]
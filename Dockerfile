FROM golang:1.17.1-buster as build
ENV GOPROXY "https://goproxy.io,direct"
WORKDIR /go/src/project/
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /bin/gohttpserver

FROM alpine
COPY --from=build /bin/gohttpserver /bin/gohttpserver
EXPOSE 8080
ENTRYPOINT [ "/bin/gohttpserver" ]
CMD [ "--help" ]
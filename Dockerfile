FROM golang:1.17.1-buster as build
COPY go.mod go.sum /go/src/project/
WORKDIR /go/src/project/
ENV GOPROXY "https://goproxy.io,direct"
RUN go mod download
COPY main.go /go/src/project/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/gohttpserver

FROM alpine
COPY --from=build /bin/gohttpserver /bin/gohttpserver
EXPOSE 8080
ENTRYPOINT [ "/bin/gohttpserver" ]
CMD [ "--help" ]
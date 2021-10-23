FROM golang:1.17.1-buster as build
COPY go.mod /go/src/project/
WORKDIR /go/src/project/
RUN go mod tidy
COPY main.go /go/src/project/
RUN  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/gohttpserver

FROM scratch
COPY --from=build /bin/gohttpserver /bin/gohttpserver
EXPOSE 8000
ENTRYPOINT [ "/bin/gohttpserver" ]
CMD [ "--help" ]
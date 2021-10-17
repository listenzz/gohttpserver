FROM golang:1.17.1-buster as build
COPY go.mod /go/src/project/
WORKDIR /go/src/project/
COPY main.go /go/src/project/
RUN go build -race -o /bin/gohttpserver

FROM busybox:1.32.1-glibc
COPY --from=build /bin/gohttpserver /bin/gohttpserver
RUN chmod +x /bin/gohttpserver
EXPOSE 8000
ENTRYPOINT [ "/bin/gohttpserver" ]
CMD [ "--help" ]
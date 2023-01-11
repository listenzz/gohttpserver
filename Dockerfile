FROM --platform=$BUILDPLATFORM golang:1.17.1-alpine as build
ENV GOPROXY "https://goproxy.io,direct"
ARG TARGETARCH
WORKDIR /go/src/project/
COPY go.* .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -o /bin/gohttpserver .

FROM --platform=$BUILDPLATFORM alpine
COPY --from=build /bin/gohttpserver /bin/gohttpserver
EXPOSE 8080
ENTRYPOINT [ "/bin/gohttpserver" ]
CMD [ "--help" ]
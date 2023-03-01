FROM golang:1.18 AS build
WORKDIR /go/src
COPY . .
RUN ls -al

RUN go mod tidy
RUN go build -o bin/immudblog-server cmd/immudblog-server/main.go
RUN go build -o bin/immudblog-cli cmd/immudblog-cli/main.go

FROM golang:1.18 AS runtime
COPY --from=build /go/src/bin/immudblog-server /usr/local/bin/immudblog-server
COPY --from=build /go/src/bin/immudblog-cli /usr/local/bin/immudblog-cli
RUN ls -al /
RUN ls -al /usr/local/bin/

EXPOSE 8080/tcp
#ENTRYPOINT ["/bin/immudblog-server"]
CMD ["immudblog-server"]

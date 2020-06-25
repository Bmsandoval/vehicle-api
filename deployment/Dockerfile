FROM golang:1.13
WORKDIR /app
RUN go get github.com/githubnemo/CompileDaemon
ENTRYPOINT CompileDaemon -log-prefix=false -build="go build -race -v -mod vendor -o .elf ./" -command="./.elf"

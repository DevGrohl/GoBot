# Dockerfile for GoBot - a Discord bot built with Golang
FROM golang:latest

# Set the Current Working Directory inside the container - GO111 is required by CompileDaemon
ENV PROJECT_DIR=/app \
    GO111MODULE=on \ 
    CGO_ENABLED=0

# Set the Current Working Directory inside the container
WORKDIR /app
RUN mkdir "/build"
COPY . .

RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon --build="go build -o main ." --command="./main"

#RUN go mod download

#RUN  CGO_ENABLED=0 GOOS=linux go build -o main .

#EXPOSE 8080

#CMD ["/app/main"]

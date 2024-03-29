FROM golang:1.19.2

WORKDIR "/app"

COPY go.mod .
COPY go.sum .

RUN go mod download && go mod verify

COPY . .

RUN go build -o trackspace ./cmd/web

EXPOSE 8080

ENTRYPOINT ["./trackspace"]

CMD ["nginx", "-g", "daemon off;"]
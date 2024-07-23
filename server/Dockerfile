# Build stage
FROM golang:alpine as server-build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

# Final stage
FROM golang:alpine

WORKDIR /app

COPY --from=server-build /app/main .

EXPOSE 8080

CMD ["./main"]
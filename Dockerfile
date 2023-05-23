#Build Stage

FROM golang:latest as builder
WORKDIR /ars-2023
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

#Run Stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
EXPOSE 8000
COPY --from=builder /ars-2023/main .
COPY ./swagger.yaml .
CMD ["./main"]
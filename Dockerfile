FROM golang:1.12 as builder
RUN mkdir /build
ADD . /build/
RUN mkdir /data
WORKDIR /build
RUN go get -d -v github.com/logrusorgru/aurora
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .
FROM scratch
COPY --from=builder /build/main /app/
WORKDIR /app
ENTRYPOINT ["./main"]
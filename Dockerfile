FROM golang:1.12 AS build
RUN mkdir /build
ADD . /build/
WORKDIR /build 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o helloworld-go .
FROM alpine:latest AS runtime
COPY --from=build /build/helloworld-go /app/
WORKDIR /app
CMD ["./helloworld-go"]
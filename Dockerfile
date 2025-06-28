FROM alpine:latest

COPY /bin /bin

CMD ["./bin/server"]

FROM alpine:3.16.0
ENTRYPOINT ["/forward"]
COPY forward /

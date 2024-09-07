FROM alpine:3.20.3
COPY forward /
ENTRYPOINT ["/forward"]

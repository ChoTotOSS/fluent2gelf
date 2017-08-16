FROM alpine
WORKDIR /fluent2gelf
COPY ./dist/fluent2gelf /app/fluent2gelf
EXPOSE 24224
CMD ["/app/fluent2gelf", "-c", "/etc/fluent2gelf.yml"]

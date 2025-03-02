FROM alpine

RUN mkdir -p /app/manager/build
COPY ./gorun /app/gorun
COPY ./manager/build /app/manager/build
RUN chmod +x /app/gorun

WORKDIR /app
CMD ["./gorun", "serve"]
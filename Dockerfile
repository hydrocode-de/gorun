FROM golang:1.23-alpine AS builder

# Install Node.js, npm, and make for the build process
RUN apk add --no-cache nodejs npm make

WORKDIR /app
COPY . .

# Run the build process using Makefile
RUN make build

FROM alpine

RUN mkdir -p /app
COPY --from=builder /app/gorun /app/gorun
RUN chmod +x /app/gorun

# Create directory for persistent data
RUN mkdir -p /data/gorun

WORKDIR /app

# Set environment variables for data storage
ENV GORUN_PATH=/data/gorun
ENV GORUN_DB=/data/gorun/gorun.db

# Expose the volume for persistent data
VOLUME /data/gorun

CMD ["./gorun", "serve"]
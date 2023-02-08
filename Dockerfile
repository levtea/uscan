FROM node:16 as frontend
WORKDIR /frontend
COPY ./ ./
RUN make pull-submodule
RUN make frontend-build


FROM golang:1.19 as builder

WORKDIR /app
COPY --from=frontend /frontend /app

RUN make statik
RUN make build

EXPOSE 4322
ENTRYPOINT ["/app/bin/uscan"]

# FROM alpine:3.16
# WORKDIR /app
# COPY --from=builder /app/bin/uscan /app

# EXPOSE 4322
# ENTRYPOINT ["/app/uscan"]

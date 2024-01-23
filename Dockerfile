FROM golang:1.21.6 as build_auth-service
ENV CGO_ENABLED 0
ARG BUILD_REF

COPY . /auth-service

WORKDIR /auth-service

RUN go build -ldflags "-X main.build=${BUILD_REF}"

FROM alpine:3.19
ARG BUILD_DATE
ARG BUILD_REF
RUN addgroup -g 1000 -S web && \
    adduser -u 1000 -h /auth-service -G web -S web \
COPY --from=build_auth-service --chown=web:web /auth-service/auth-service /service/auth-service
WORKDIR /service
USER web
CMD ["./auth-service"]
EXPOSE 8080

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="auth-service" \
      org.opencontainers.image.authors="Dther <dtherhtun.cw@gmail.com>" \
      org.opencontainers.image.source="https://github.com/lovehotel24/auth-service" \
      org.opencontainers.image.revision="${BUILD_REF}" \
      org.opencontainers.image.vendor="Love hotel"
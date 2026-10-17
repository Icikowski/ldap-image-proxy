ARG BUILD_VERSION=unknown
ARG BUILD_GIT_REF=unknown
ARG BUILD_TIMESTAMP=unknown
ARG BUILD_SHORT_REF=unknown

FROM golang:1.27 AS builder
ARG BUILD_VERSION
ARG BUILD_GIT_REF
ARG BUILD_TIMESTAMP
ARG TARGETARCH
ENV BUILD_VERSION=${BUILD_VERSION}
ENV BUILD_GIT_REF=${BUILD_GIT_REF}
ENV BUILD_TIMESTAMP=${BUILD_TIMESTAMP}
WORKDIR /app
COPY . .
RUN wget -nv -O- https://golangci-lint.run/install.sh | sh -s -- -b /usr/bin
RUN wget -nv -O- https://taskfile.dev/install.sh | sh -s -- -b /usr/bin
RUN go mod download -x
RUN task generate fmt build-static build-healthcheck

FROM gcr.io/distroless/static:nonroot
ARG BUILD_VERSION
ARG BUILD_GIT_REF
ARG BUILD_TIMESTAMP
ARG BUILD_SHORT_REF
LABEL org.opencontainers.image.created="${BUILD_TIMESTAMP}" \
    org.opencontainers.image.authors="Piotr Icikowski" \
    org.opencontainers.image.url="https://sr.ht/~icikowski/ldap-image-proxy" \
    org.opencontainers.image.documentation="https://git.sr.ht/~icikowski/ldap-image-proxy/blob/${BUILD_SHORT_REF}/README_CONTAINER.md" \
    org.opencontainers.image.source="https://git.sr.ht/~icikowski/ldap-image-proxy" \
    org.opencontainers.image.version="${BUILD_VERSION}" \
    org.opencontainers.image.revision="${BUILD_GIT_REF}" \
    org.opencontainers.image.vendor="Piotr Icikowski" \
    org.opencontainers.image.licenses="AGPL-3.0" \
    org.opencontainers.image.title="LDAP Image Proxy" \
    org.opencontainers.image.description="A proxy service for serving LDAP user avatars as images." \
    org.opencontainers.image.base.name="gcr.io/distroless/static:nonroot"
WORKDIR /app
COPY --from=builder /app/bin/ldap-image-proxy .
COPY --from=builder /app/bin/healthcheck .
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 CMD ["/app/healthcheck"]
ENTRYPOINT ["/app/ldap-image-proxy"]

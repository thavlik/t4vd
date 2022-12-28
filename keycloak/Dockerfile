ARG BASE_IMAGE=quay.io/keycloak/keycloak:latest
FROM ${BASE_IMAGE} as builder
ENV KC_HEALTH_ENABLED=true
ENV KC_METRICS_ENABLED=true
ENV KC_DB=postgres
WORKDIR /opt/keycloak
RUN /opt/keycloak/bin/kc.sh build
FROM ${BASE_IMAGE}
COPY --from=builder /opt/keycloak/ /opt/keycloak/
ENV KC_LOG_LEVEL=INFO
ENV KC_PROXY=edge
ENV KC_HTTP_ENABLED=true
ENV KC_HOSTNAME_STRICT=false
COPY entrypoint.sh /
USER root
RUN chmod +x /entrypoint.sh \
    && mkdir /opt/postgres \
    && chown -R keycloak:keycloak /opt/postgres
USER keycloak
ENTRYPOINT ["/entrypoint.sh"]

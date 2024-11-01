FROM debian:trixie-slim
ARG TARGETOS
ARG TARGETARCH
WORKDIR /
USER 65532:65532
COPY function-${TARGETOS:-linux}-${TARGETARCH} /function

ENTRYPOINT ["/function"]

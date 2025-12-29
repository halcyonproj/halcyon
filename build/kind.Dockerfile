FROM alpine

ARG VERSION

RUN apk add -u --no-cache curl docker-cli && \
    curl -Lo ./kind https://github.com/kubernetes-sigs/kind/releases/download/$VERSION/kind-$(uname)-amd64 && \
    chmod +x ./kind && mv ./kind /bin

ENTRYPOINT [ "kind" ]

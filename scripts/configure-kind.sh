#!/usr/bin/env bash

KIND_VERSION=$1
SCRIPT_DIR=$(realpath $0)
BIN_DIR=$(dirname ${SCRIPT_DIR})/../bin
echo $BIN_DIR

if [ -z "${KIND_VERSION}" ]; then
    echo "Failed to build kind: must supply kind version as an arg to $0" >&2
    exit 1
fi

docker build --build-arg VERSION=${KIND_VERSION} -t kind:${KIND_VERSION} - < ./build/kind.Dockerfile

cat << EOF > ${BIN_DIR}/kind
KUBECONFIG_DIR="\${KUBECONFIG_DIR:-\${HOME}/.kube/}"
KIND_VERSION="$1"

docker run -i -v /var/run/docker.sock:/var/run/docker.sock -v \$KUBECONFIG_DIR:/root/.kube/ kind:${KIND_VERSION} \$@
EOF

chmod +x ${BIN_DIR}/kind

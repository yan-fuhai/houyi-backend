#!/bin/bash

OS=linux
ARCH=amd64

COMPONENT=backend
BUILD_OUT_DIR=~/houyi/${COMPONENT}
WORK_DIR=../

mkdir -p ${BUILD_OUT_DIR}
CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -tags netgo -o ${BUILD_OUT_DIR}/${COMPONENT} -v ${WORK_DIR}/main.go

RUN_SHELL=run.sh
chmod u+x ${RUN_SHELL}
cp ${RUN_SHELL} ${BUILD_OUT_DIR}/

cat <<EOF > Dockerfile
FROM alpine:3.7
COPY ${COMPONENT} /opt/ms/
COPY ${RUN_SHELL} /opt/ms/
EXPOSE 80
WORKDIR /opt/ms/
ENTRYPOINT ["/opt/ms/${RUN_SHELL}"]
EOF
mv Dockerfile ${BUILD_OUT_DIR}/

BUILD_RUN_SHELL_DOCKER=build-docker.sh
chmod u+x ${BUILD_RUN_SHELL_DOCKER}
cp ${BUILD_RUN_SHELL_DOCKER} ${BUILD_OUT_DIR}/

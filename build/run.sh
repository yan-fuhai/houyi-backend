#!/bin/sh

ROOT_CMD='./backend'

if [[ ${LOG_LEVEL} ]]; then
  ROOT_CMD="${ROOT_CMD} --log.level=${LOG_LEVEL}"
fi

echo $ROOT_CMD
eval $ROOT_CMD

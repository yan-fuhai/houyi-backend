#!/bin/sh

ROOT_CMD='./backend'

if [[ ${LOG_LEVEL} ]]; then
  ROOT_CMD="${ROOT_CMD} --log.level=${LOG_LEVEL}"
fi

if [[ ${HTTP_LISTEN_PORT} ]]; then
  ROOT_CMD="${ROOT_CMD} --http.listen.port=${HTTP_LISTEN_PORT}"
fi

# Strategy Manager
if [[ ${STRATEGY_MANAGER_ADDR} ]]; then
  ROOT_CMD="${ROOT_CMD} --strategy.manager.addr=${STRATEGY_MANAGER_ADDR}"
fi
if [[ ${STRATEGY_MANAGER_PORT} ]]; then
  ROOT_CMD="${ROOT_CMD} --strategy.manager.port=${STRATEGY_MANAGER_PORT}"
fi

echo $ROOT_CMD
eval $ROOT_CMD

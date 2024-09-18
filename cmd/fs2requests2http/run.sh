#!/bin/sh

export ENV_TRUSTED_DIR_NAME=${ENV_TRUSTED_DIR_NAME:-./sample.d/requests.d}
export ENV_TARGET_URL=${ENV_TARGET_URL:-http://localhost:80}
export ENV_TARGET_TYP=${ENV_TARGET_TYP:-application/octet-stream}
export ENV_MAX_BODY_SIZE=${ENV_MAX_BODY_SIZE:-65536}
export ENV_SAVE_REQ=${ENV_SAVE_REQ:-false}
export ENV_SAVE_NAME=${ENV_SAVE_NAME:-./err.dat}
export ENV_LOG_LEVEL=${ENV_LOG_LEVEL:-info}

./fs2requests2http

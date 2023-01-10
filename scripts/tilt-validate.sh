#!/usr/bin/env bash

source .envrc

RET=$(tilt alpha tiltfile-result)
if [ $? -ne 0 ]; then
  echo "$RET"
else
  echo "Tiltfile validated successfully!"
fi

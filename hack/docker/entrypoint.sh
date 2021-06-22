#!/usr/bin/env bash

if [ "$1"  == "bash" ]; then
    exec /bin/bash
else
    exec /usr/bin/drone-sonar "$@"
fi
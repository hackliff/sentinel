#! /bin/bash
set -euo pipefail
IFS=$'\n\t'

# note : serf parses script string to find key=value syntax for dispatching
# information. Therefor we use a space instead for handler flags.
serf agent -rpc-addr=0.0.0.0:7373 \
  -event-handler="./sentinel -api-key ${PUSHBULLET_API_KEY}"

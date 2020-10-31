#!/bin/bash

# Inspired by https://github.com/alexedwards/go-reload

function monitor() {
  fswatch -r $directories |
  while read line; do
    restart
  done
}

# Terminate and rerun the main Go program
function restart {
  kill $CMD_PID
  echo ">> Reloading..."
  $cmd &
  CMD_PID=$!
}

if [ -z $(which fswatch) ]; then
  echo "Please install fswatch"
  echo "https://github.com/emcrisostomo/fswatch"
  exit 1
fi

cmd=$1
shift 1
directories=$@
echo "== reload.sh"

# Start the main Go program
echo ">> Watching directories $directories"
echo ">> Running command '$cmd', CTRL+C to stop"
$cmd &
CMD_PID=$!

monitor &

wait

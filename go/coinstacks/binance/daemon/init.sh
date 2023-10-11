#!/bin/bash

set -e

[ "$DEBUG" = "true" ] && set -x

HOME_DIR=/root/.bnbchaind
CONFIG_DIR=$HOME_DIR/config

# shapshots provided by: https://github.com/bnb-chain/bc-snapshots
if [ -n "$SNAPSHOT" ] && [ ! -f "$HOME_DIR/data/priv_validator_state.json" ]; then
  rm -rf $HOME_DIR/data;
  mkdir -p $HOME_DIR/data;
  wget -c $SNAPSHOT -O - | tar xvf - -C $DATA_DIR
fi

if [ ! -d "$CONFIG_DIR" ]; then
  mkdir -p $CONFIG_DIR
  cp app.toml config.toml genesis.json $CONFIG_DIR
fi

start() {
  bnbchaind start \
    --moniker unchained \
    --rpc.laddr tcp://0.0.0.0:26657 &
  DAEMON_PID="$!"

  bnbcli api-server \
    --chain-id Binance-Chain-Tigris \
    --laddr tcp://0.0.0.0:1317 &
  API_SERVER_PID="$!"
}

stop() {
  echo "Catching signal and sending to PID: $API_SERVER_PID" && kill $API_SERVER_PID
  while $(kill -0 $PID 2>/dev/null); do sleep 1; done
  echo "Catching signal and sending to PID: $DAEMON_PID" && kill $DAEMON_PID
  while $(kill -0 $PID 2>/dev/null); do sleep 1; done
}

trap 'stop' TERM INT
start
wait $PID
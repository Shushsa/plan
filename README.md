# Plan Node

Plan is a blockchain built on the Cosmos SDK that allows you to build your own passive income. 

## Requirements

 - At least 20 GB disk space
 - Golang >= 1.12 with the $GOPATH and $GOBIN variables set
 - *nix is not required, but preferred
 
## Building

Just run `./scripts/build.sh` or `make install`, then you'll be able to use it from your console by calling "pland" or "plancli".
 
## Running local network for dev purpose

Run `./scripts/init-genesis.sh`, then just start the node by `pland start` 

## Running the node

After building the node, you need to fetch the latest blocks from here https://plan-crypto.s3-eu-west-1.amazonaws.com/blockchain.tar.gz

To unpack them, please do

```shell script
tar -xvf blockchain.tar.gz -C ~/.pland/
```

After that, you'll be able to start your node via  ```pland start```

## What's next?

All the other things described here https://docs.plan-crypto.com/en/ including configuring your local console interface and running the API server.

Here are https://github.com/plan-crypto/nodejs_examples some examples of working with the blockchain including account generation, txs sending, etc.

If you still have any questions, please feel free to contact me via Telegram: @plan_developer
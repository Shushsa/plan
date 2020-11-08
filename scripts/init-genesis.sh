#!/usr/bin/env bash
# Sets up a network from scratch

PASSWORD="abcdef123"

rm -rf ~/.pland/
rm -rf ~/.plancli/

# Initializes blockchain
pland init anonymous --chain-id plan

# Creates a few testing accounts
echo ${PASSWORD} | plancli keys add jack
echo ${PASSWORD} | plancli keys add alice

# Adds genesis account
pland add-genesis-account $(plancli keys show jack -a) 10000000000000pln,10000000000000stake

# Конфигурируем cli
plancli config chain-id plan
plancli config output json
plancli config indent true
plancli config trust-node true

# Creates genesis tx
echo ${PASSWORD} | pland gentx --name jack --amount 100000000stake

pland collect-gentxs
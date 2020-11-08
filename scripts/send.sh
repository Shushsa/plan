#!/usr/bin/env bash
# Sends coins from one test account to another
plancli tx send  --gas auto --gas-adjustment 1.5 $(plancli keys show jack -a) $(plancli keys show alice -a) 2000000000pln
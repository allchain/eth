#!/bin/bash

n=$1

geth=./build/bin/geth
datadir=~/Library/Ethereum
verbosity=4
port=30302
rpcport=8544

#case ${n} in
#
#    ;;
#    2) ${geth} --datadir ${datadir}/2 --verbosity ${verbosity} --port `expr ${port} + ${n}` --rpc --rpcport `expr ${port} + ${n}` console
#    ;;
#    3) ${geth} --datadir ${datadir}/3 --verbosity ${verbosity} --port 30305 --rpc --rpcport 8547 console
#    ;;
#    4) ${geth} --datadir ${datadir}/3 --verbosity ${verbosity} --port 30305 --rpc --rpcport 8547 console
#    ;;
#    5) ${geth} --datadir ${datadir}/3 --verbosity ${verbosity} --port 30305 --rpc --rpcport 8547 console
#    ;;
#esac

${geth} --datadir ${datadir}/${n} --verbosity ${verbosity} --port `expr ${port} + ${n}` --rpc --rpcport `expr ${rpcport} + ${n}` --rpcapi "admin,miner" --syncmode "full" console
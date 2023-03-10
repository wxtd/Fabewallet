#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error
set -e

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1
starttime=$(date +%s)
CC_SRC_LANGUAGE="go"
CC_SRC_PATH="../asset-transfer-basic/experiment/"

# clean out any old identites in the wallet
rm -rf go/wallet/*

# launch network; create channel and join peer to channel
pushd ../test-network
./network.sh down
./network.sh up createChannel -ca -s couchdb
./network.sh deployCC -ccn fabtask -ccv 1.0 -ccl ${CC_SRC_LANGUAGE} -ccp ${CC_SRC_PATH}

#./network.sh deployCC -ccn fabtask -ccv 1.0 -ccl go -ccp ../asset-transfer-basic/experiment/ -c mychannel2

popd

cat <<EOF

Chaincode start !

EOF

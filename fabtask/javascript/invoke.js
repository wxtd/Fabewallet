/*
 * Copyright IBM Corp. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';


var format = require('string-format')
const { Gateway, Wallets } = require('fabric-network');
const fs = require('fs');
const path = require('path');

async function main() {
    try {
        // load the network configuration
        const ccpPath = path.resolve(__dirname, '..', '..', 'test-network', 'organizations', 'peerOrganizations', 'org1.example.com', 'connection-org1.json');
        let ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const identity = await wallet.get('appUser1');
        if (!identity) {
            console.log('An identity for the user "appUser1" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: 'appUser1', discovery: { enabled: true, asLocalhost: true } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mychannel');

        // Get the contract from the network.
        const contract = network.getContract('fabtask7');

        
        for(let j=0; j < 100; j++) {
            var curr_servers = []
            const result = await contract.evaluateTransaction('queryAllServer');
        // console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
            var re = JSON.parse(result)
            for(let i = 0; i < re.length; i++) {
                var record = {}
                record = {'mem': re[i]['Record']['memory'],
                          'name': re[i]['Key'],
                          'cpu': re[i]['Record']['CPU'],
                          'reputation': re[i]['Record']['reputation']}
                curr_servers.push(record)
                // console.log(re[i]['Record']['memory'])
            }
            var u = 0.8
            curr_servers = curr_servers.filter(record => parseInt(record['mem']) > 150 * u );
            curr_servers.sort((a, b) => a.reputation - b.reputation)
            var task = {'task': 'task'+ j, // taskid
                        'publisher': 'Tecent',
                        'executor': curr_servers[0]['name'],
                        'starttime': new Date().getTime(),
                        'ddl': '1',
                        'endtime': new Date().getTime() + Math.ceil(Math.random()*2),
                        'mem': '2',
                        'cpu': '2',}
            // console.log(task)
            await contract.submitTransaction('createTask', task['task'], task['publisher'], task['executor'], task['mem'], task['cpu'], task['starttime'], task['ddl'], task['endtime']);
            console.log('Transaction has been submitted');
            
            // for(let cc = 0; cc < curr_servers.length; cc++)
            //     console.log(curr_servers[cc])
        }


        // Submit the specified transaction.
        
        // Disconnect from the gateway.
        await gateway.disconnect();

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        process.exit(1);
    }
}

main();

/*
 * Copyright IBM Corp. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { Contract } = require('fabric-contract-api');

// class fabtask extends Contract {

//     async initLedger(ctx) {
//         console.info('============= START : Initialize Ledger ===========');
//         const cars = [
//             {
//                 color: 'blue',
//                 make: 'Toyota',
//                 model: 'Prius',
//                 owner: 'Tomoko',
//             },
//             {
//                 color: 'red',
//                 make: 'Ford',
//                 model: 'Mustang',
//                 owner: 'Brad',
//             },
//             {
//                 color: 'green',
//                 make: 'Hyundai',
//                 model: 'Tucson',
//                 owner: 'Jin Soo',
//             },
//             {
//                 color: 'yellow',
//                 make: 'Volkswagen',
//                 model: 'Passat',
//                 owner: 'Max',
//             },
//             {
//                 color: 'black',
//                 make: 'Tesla',
//                 model: 'S',
//                 owner: 'Adriana',
//             },
//             {
//                 color: 'purple',
//                 make: 'Peugeot',
//                 model: '205',
//                 owner: 'Michel',
//             },
//             {
//                 color: 'white',
//                 make: 'Chery',
//                 model: 'S22L',
//                 owner: 'Aarav',
//             },
//             {
//                 color: 'violet',
//                 make: 'Fiat',
//                 model: 'Punto',
//                 owner: 'Pari',
//             },
//             {
//                 color: 'indigo',
//                 make: 'Tata',
//                 model: 'Nano',
//                 owner: 'Valeria',
//             },
//             {
//                 color: 'brown',
//                 make: 'Holden',
//                 model: 'Barina',
//                 owner: 'Shotaro',
//             },
//         ];

//         for (let i = 0; i < cars.length; i++) {
//             cars[i].docType = 'car';
//             await ctx.stub.putState('CAR' + i, Buffer.from(JSON.stringify(cars[i])));
//             console.info('Added <--> ', cars[i]);
//         }
//         console.info('============= END : Initialize Ledger ===========');
//     }

//     async queryCar(ctx, carNumber) {
//         const carAsBytes = await ctx.stub.getState(carNumber); // get the car from chaincode state
//         if (!carAsBytes || carAsBytes.length === 0) {
//             throw new Error(`${carNumber} does not exist`);
//         }
//         console.log(carAsBytes.toString());
//         return carAsBytes.toString();
//     }

//     async createCar(ctx, carNumber, make, model, color, owner) {
//         console.info('============= START : Create Car ===========');

//         const car = {
//             color,
//             docType: 'car',
//             make,
//             model,
//             owner,
//         };

//         await ctx.stub.putState(carNumber, Buffer.from(JSON.stringify(car)));
//         console.info('============= END : Create Car ===========');
//     }

//     async queryAllCars(ctx) {
//         const startKey = '';
//         const endKey = '';
//         const allResults = [];
//         for await (const {key, value} of ctx.stub.getStateByRange(startKey, endKey)) {
//             const strValue = Buffer.from(value).toString('utf8');
//             let record;
//             try {
//                 record = JSON.parse(strValue);
//             } catch (err) {
//                 console.log(err);
//                 record = strValue;
//             }
//             allResults.push({ Key: key, Record: record });
//         }
//         console.info(allResults);
//         return JSON.stringify(allResults);
//     }

//     async changeCarOwner(ctx, carNumber, newOwner) {
//         console.info('============= START : changeCarOwner ===========');

//         const carAsBytes = await ctx.stub.getState(carNumber); // get the car from chaincode state
//         if (!carAsBytes || carAsBytes.length === 0) {
//             throw new Error(`${carNumber} does not exist`);
//         }
//         const car = JSON.parse(carAsBytes.toString());
//         car.owner = newOwner;

//         await ctx.stub.putState(carNumber, Buffer.from(JSON.stringify(car)));
//         console.info('============= END : changeCarOwner ===========');
//     }

// }

// module.exports = fabtask;

class fabtask extends Contract {
    getRandomInt(min, max) {
        return Math.floor(Math.random() * (max - min + 1) + min);
    }

    sleep(milliSeconds){ 
        var StartTime =new Date().getTime(); 
        let i = 0;
        while (new Date().getTime() <StartTime+milliSeconds);
    
    }

    async initLedger(ctx) {
        console.info('============= START : Initialize Ledger ===========');

        const servers = [
            {
                name: 'Tecent',
                memory: ''+150, // this.getRandomInt(150, 200)
                CPU: '0',
                reputation: '80',
            },
            {
                name: 'Alibaba',
                memory: ''+150,
                CPU: '0',
                reputation: '80',
            },
            {
                name: 'ByteDance',
                memory: ''+150,
                CPU: '0',
                reputation: '80',
            },
            {
                name: 'Baidu',
                memory: ''+150,
                CPU: '0',
                reputation: '80',
            },
            {
                name: 'Meituan',
                memory: ''+150,
                CPU: '0',
                reputation: '80',
            },
            {
                name: 'Apple',
                memory: ''+150,
                CPU: '0',
                reputation: '80',
            },
            {
                name: 'Microsoft',
                memory: ''+150,
                CPU: '0',
                reputation: '80',
            },
            {
                name: 'Google',
                memory: ''+150,
                CPU: '0',
                reputation: '80',
            },
            {
                name: 'Wangyi',
                memory: ''+150,
                CPU: '0',
                reputation: '80',
            },
            {
                name: 'Wechat',
                memory: ''+150,
                CPU: '0',
                reputation: '80',
            },
        ];

        for (let i = 0; i < servers.length; i++) {
            servers[i].docType = 'server';
            await ctx.stub.putState(servers[i].name, Buffer.from(JSON.stringify(servers[i])));
            console.info('Added <--> ', servers[i]);
        }
        console.info('============= END : Initialize Ledger ===========');
    }

    async queryServer(ctx, serverName) {
        const serverAsBytes = await ctx.stub.getState(serverName); // get the server from chaincode state
        if (!serverAsBytes || serverAsBytes.length === 0) {
            throw new Error(`${serverName} does not exist`);
        }
        console.log(serverAsBytes.toString());
        return serverAsBytes.toString();
    }

    async queryTask(ctx, taskNumber) {
        const taskAsBytes = await ctx.stub.getState(taskNumber); // get the server from chaincode state
        if (!taskAsBytes || taskAsBytes.length === 0) {
            throw new Error(`${taskNumber} does not exist`);
        }
        console.log(taskAsBytes.toString());
        return taskAsBytes.toString();
    }

    async createTask(ctx, taskNumber, publisher, executor, mem_consuming, cpu_consuming, release_time,
        time_consuming, deadline,) {
        console.info('============= START : Create Task ===========');

        const task = {
            docType: 'task',
            publisher,
            executor,
            mem_consuming,
            cpu_consuming,
            release_time,
            time_consuming,
            deadline,
        };

        const serverAsBytes = await ctx.stub.getState(executor); // get the server from chaincode state
        const server = JSON.parse(serverAsBytes.toString());
        if((parseInt(release_time)+parseInt(time_consuming)) > parseInt(deadline)) { // 锟斤拷时
            server.reputation = '' + (parseInt(server.reputation) - 1);
        }
        else {
            server.reputation = '' + (parseInt(server.reputation) + 1);
        }
        server.memory = '' + (parseInt(server.memory) - parseInt(mem_consuming))
        await ctx.stub.putState(server.name, Buffer.from(JSON.stringify(server)));
        this.sleep(parseInt(time_consuming) * 1000)
        await ctx.stub.putState(taskNumber, Buffer.from(JSON.stringify(task)));
        console.info('============= END : Create Task ===========');
    }

    async queryAll(ctx) {
        const startKey = '';
        const endKey = '';
        const allResults = [];
        for await (const {key, value} of ctx.stub.getStateByRange(startKey, endKey)) {
            const strValue = Buffer.from(value).toString('utf8');
            let record;
            try {
                record = JSON.parse(strValue);
            } catch (err) {
                console.log(err);
                record = strValue;
            }
            allResults.push({ Key: key, Record: record });
        }
        console.info(allResults);
        return JSON.stringify(allResults);
    }

    async queryAllServer(ctx) {
        const startKey = '';
        const endKey = '';
        const allResults = [];
        for await (const {key, value} of ctx.stub.getStateByRange(startKey, endKey)) {
            const strValue = Buffer.from(value).toString('utf8');
            let record;
            try {
                record = JSON.parse(strValue);
            } catch (err) {
                console.log(err);
                record = strValue;
            }
            if(record.docType == 'server') {
                allResults.push({ Key: key, Record: record });
            }
        }
        console.info(allResults);
        return JSON.stringify(allResults);
    }

    // async changeCarOwner(ctx, carNumber, newOwner) {
    //     console.info('============= START : changeCarOwner ===========');

    //     const carAsBytes = await ctx.stub.getState(carNumber); // get the car from chaincode state
    //     if (!carAsBytes || carAsBytes.length === 0) {
    //         throw new Error(`${carNumber} does not exist`);
    //     }
    //     const car = JSON.parse(carAsBytes.toString());
    //     car.owner = newOwner;

    //     await ctx.stub.putState(carNumber, Buffer.from(JSON.stringify(car)));
    //     console.info('============= END : changeCarOwner ===========');
    // }

}

module.exports = fabtask;

/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"

	"github.com/go-redis/redis"
	"strings"
)

func main() {
	// 为了使用 worker 线程池并且收集他们的结果，我们需要 2 个通道。
	jobs := make(chan int, 100)
	results := make(chan int, 100)

	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	for j := 1; j <= 10; j++ {
		jobs <- j
	}
	close(jobs)

	for a := 1; a <= 10; a++ {
		<-results
	}
}

func populateWallet(wallet *gateway.Wallet) error {
	credPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return errors.New("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	err = wallet.Put("appUser", identity)
	if err != nil {
		return err
	}
	return nil
}

var client *redis.Client

func initRedis()(err error){
	client = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",  // 指定
		Password: "",
		DB:0,		// redis一共16个库，指定其中一个库即可
	})
    _,err = client.Ping().Result()
	return
}

func worker(id int, jobs <-chan int, results chan<- int) {

	for j := range jobs {
		
		fmt.Println("worker", id, "processing job", j)


		err_redis := initRedis()
		if err_redis != nil {
			fmt.Printf("connect redis failed! err : %v\n",err_redis)
			return
		}
		// fmt.Println("redis connection success!")

		os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
		wallet, err := gateway.NewFileSystemWallet("wallet")
		if err != nil {
			fmt.Printf("Failed to create wallet: %s\n", err)
			os.Exit(1)
		}

		if !wallet.Exists("appUser") {
			err = populateWallet(wallet)
			if err != nil {
				fmt.Printf("Failed to populate wallet contents: %s\n", err)
				os.Exit(1)
			}
		}

		ccpPath := filepath.Join(
			"..",
			"..",
			"test-network",
			"organizations",
			"peerOrganizations",
			"org1.example.com",
			"connection-org1.yaml",
		)

		gw, err := gateway.Connect(
			gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
			gateway.WithIdentity(wallet, "appUser"),
		)
		if err != nil {
			fmt.Printf("Failed to connect to gateway: %s\n", err)
			os.Exit(1)
		}
		defer gw.Close()

		network, err := gw.GetNetwork("mychannel")
		if err != nil {
			fmt.Printf("Failed to get network: %s\n", err)
			os.Exit(1)
		}
		
		v_transaction, e_transaction := client.RPop("transaction_queue").Result()
		if e_transaction != nil {
			fmt.Printf("get transaction_queue failed! err : %v\n",e_transaction)
			return
		}

		arr := strings.Split(v_transaction, " ")

		_, payer, remittee, money := arr[0], arr[1], arr[2], arr[3]
		
		fmt.Println(payer, " ", remittee, " ", money)

		contract := network.GetContract("ewallet")

		result, err := contract.SubmitTransaction("transferAccount", payer, remittee, money)
		if err != nil {
			fmt.Printf("Failed to submit transaction: %s\n", err)
			os.Exit(1)
		}
		fmt.Println(string(result))
		
		
		results <- j*2
	}
}
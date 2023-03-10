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
	
	"strconv"
)

func main() {

	jobs := make(chan int, 100000)
	results := make(chan int, 100000)

	// worker num
	for w := 1; w <= 100; w++ {
		go worker(w, jobs, results)
	}

	cnt := 10000

	for i := 0; i <= cnt; i++ {
		jobs <- i 
	}

	close(jobs)

	for a := 0; a <= cnt; a++ {
		<-results
	}

	// for i := 0; i <= 10000; i++ {
	// 	result, err := contract.SubmitTransaction("saveMoney", "owner"+strconv.Itoa(i) , "1000000")
	// 	if err != nil {
	// 		fmt.Printf("Failed to submit transaction: %s\n", err)
	// 		os.Exit(1)
	// 	}
	// 	fmt.Println(i, string(result))
	// }

	// result, err := contract.SubmitTransaction("saveMoney", "owner2" , "1000000")
	// if err != nil {
	// 	fmt.Printf("Failed to submit transaction: %s\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println( string(result))

	// result, err = contract.SubmitTransaction("saveMoney", "owner1" , "1000000")
	// if err != nil {
	// 	fmt.Printf("Failed to submit transaction: %s\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println( string(result))
	
	// result, err = contract.SubmitTransaction("saveMoney", "b", "100")
	// if err != nil {
	// 	fmt.Printf("Failed to submit transaction: %s\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(string(result))

	// result, err := contract.EvaluateTransaction("queryAllAccounts")
	// if err != nil {
	// 	fmt.Printf("Failed to evaluate transaction: %s\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(string(result))


	// result, err = contract.EvaluateTransaction("queryCar", "CAR10")
	// if err != nil {
	// 	fmt.Printf("Failed to evaluate transaction: %s\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(string(result))

	// _, err = contract.SubmitTransaction("changeCarOwner", "CAR10", "Archie")
	// if err != nil {
	// 	fmt.Printf("Failed to submit transaction: %s\n", err)
	// 	os.Exit(1)
	// }

	// result, err = contract.EvaluateTransaction("queryCar", "CAR10")
	// if err != nil {
	// 	fmt.Printf("Failed to evaluate transaction: %s\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(string(result))
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


func worker(id int, jobs <-chan int, results chan<- int) {

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

	contract := network.GetContract("ewallet")

	for i := range jobs {

		result, err := contract.SubmitTransaction("saveMoney", "owner"+strconv.Itoa(i) , "1000000")
		if err != nil {
			fmt.Printf("Failed to submit transaction: %s\n", err)
			os.Exit(1)
		}
		fmt.Println(i, string(result))

	}

	results <- 0	
}
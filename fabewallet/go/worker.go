/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
	//"math/rand"

	"sync"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"

	"strings"

	"github.com/go-redis/redis"
)


// 线程数量
var thread_num = 100

// 轮询时间
var t_i = time.Second * 2

// 令牌
var token sync.Mutex

func main() {
	// 为了使用 worker 线程池并且收集他们的结果，我们需要 2 个通道。
	jobs := make(chan string, 100000)
	results := make(chan int, 100000)

	// worker num
	for w := 1; w <= thread_num; w++ {
		go worker(w, jobs, results) 
	}

	var cnt = 0

	err := initRedis()
	if err != nil {
		fmt.Printf("connect redis failed! err : %v\n",err)
		return
	}

	// get from transaction_queue
	for true {
		v_transaction, e_transaction := client.RPop("transaction_queue").Result()
		if e_transaction != nil {
			fmt.Printf("get transaction_queue failed! err : %v\n", e_transaction)
			break
		}

		jobs <- v_transaction
		
		cnt++
	}

	// input := make(chan interface{})
    // //producer - produce the messages
    // go func() {
    //     // for i := 0; i < 5; i++ {
    //     //     input <- i
    //     // }
    //     // input <- "hello, world"
    // }()
    t1 := time.NewTimer(t_i)  
    for {
        select {
            //consumer - consume the messages
            // case msg := <-input:
            //     fmt.Println(msg)
            case <-t1.C:
                	// get from merge_queue
				// val1, err1 := client.LRange("merge_queue",0,-1).Result()
				// if err1 != nil {
				// 	fmt.Println(err1)
				// }
				for true {
					v_transaction, e_transaction := client.RPop("merge_queue").Result()
					if e_transaction != nil {
						fmt.Printf("get merge_queue failed! err : %v\n", e_transaction)
						break
					}

					jobs <- v_transaction
					cnt++

				}
                t1.Reset(t_i)
        }
    }
	close(jobs)

	for a := 1; a <= cnt; a++ {
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

func initRedis() (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379", // 指定
		Password: "",
		DB:       0, // redis一共16个库，指定其中一个库即可
	})
	_, err = client.Ping().Result()
	return
}

func worker(id int, jobs <-chan string, results chan<- int) {
	
	err_redis := initRedis()
	if err_redis != nil {
		fmt.Printf("connect redis failed! err : %v\n", err_redis)
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

		contract := network.GetContract("ewallet")


	for j := range jobs {

		fmt.Println("worker", id, "processing job")

		arr := strings.Split(j, " ")

		funcName := arr[0]

		if funcName == "transferAccount_multiParty" {

			arg_map := make(map[string]string)
			_ = json.Unmarshal([]byte(arr[1]), &arg_map)
			isUnlocked := true

			token.Lock()
			for key, _ := range arg_map {
				v_payer, err_payer := client.HGet("test", key).Result()
				if err_payer != nil {
					fmt.Printf("err_player")
					// unlocked
					continue
				} else {
					// fmt.Println("v_payer:", v_payer)
					if v_payer != "0" && v_payer != "" {
						//locked
						isUnlocked = false
						break
					} else {
						// unlocked
						continue
					}
				}
			}
			
			if isUnlocked {

				var mynote [][]string
				_ = json.Unmarshal([]byte(arr[2]), &mynote)
				alltxnum := len(mynote)

				for key, _ := range arg_map {
					client.HIncrBy("test", key, 1)
				}
				token.Unlock()

				result, err := contract.SubmitTransaction(funcName, arr[1], arr[2])
				if err != nil {
					fmt.Printf("Failed to submit transaction_mutiply: %s\n", err)
					for i := 0; i < alltxnum; i++ {
						RecordTimetoFile("end_improve.txt", true)
					}
					// os.Exit(1)
				} else {
					for i := 0; i < alltxnum; i++ {
						RecordTimetoFile("end_improve.txt", false)
					}
				}
				fmt.Println(string(result))

				for key, _ := range arg_map {
					client.HIncrBy("test", key, -1)
				}

			} else {
				token.Unlock()
				client.LPush("merge_queue", j)
			}


			// arg_map := make(map[string]string)
			// _ = json.Unmarshal([]byte(arr[1]), &arg_map)
			// isUnlocked := true 			

			// for key, _ := range arg_map {
			// 	v_payer, err_payer := client.HGet("test", key).Result()
			// 	if err_payer != nil {
			// 		fmt.Printf("err_player")
			// 		// unlocked
			// 		continue
			// 	} else {
			// 		// fmt.Println("v_payer:", v_payer)
			// 		if v_payer != "0" && v_payer != "" {
			// 			//locked
			// 			isUnlocked = false
			// 			break
			// 		} else {
			// 			// unlocked
			// 			continue
			// 		}
			// 	}
			// }
			
			// if isUnlocked {

			// 	var mynote [][]string
			// 	_ = json.Unmarshal([]byte(arr[2]), &mynote)
			// 	alltxnum := len(mynote)

			// 	for key, _ := range arg_map {
			// 		client.HIncrBy("test", key, 1)
			// 	}

			// 	isUnlocked_again := true 			

			// 	for key, _ := range arg_map {
			// 		v_payer, err_payer := client.HGet("test", key).Result()
			// 		if err_payer != nil {
			// 			fmt.Printf("err_player")
			// 			// unlocked
			// 			continue
			// 		} else {
			// 			// fmt.Println("v_payer:", v_payer)
			// 			if v_payer != "1" {
			// 				//locked
			// 				isUnlocked_again = false
			// 				break
			// 			} else {
			// 				// unlocked
			// 				continue
			// 			}
			// 		}
			// 	}

			// 	if isUnlocked_again {
			// 		result, err := contract.SubmitTransaction(funcName, arr[1], arr[2])
			// 		if err != nil {
			// 			fmt.Printf("Failed to submit transaction_mutiply: %s\n", err)
			// 			for i := 0; i < alltxnum; i++ {
			// 				RecordTimetoFile("end_improve.txt", true)
			// 			}
			// 			// os.Exit(1)
			// 		} else {
			// 			for i := 0; i < alltxnum; i++ {
			// 				RecordTimetoFile("end_improve.txt", false)
			// 			}
			// 		}
			// 		fmt.Println(string(result))

			// 		for key, _ := range arg_map {
			// 			client.HIncrBy("test", key, -1)
			// 		}

			// 	} else {

			// 		for key, _ := range arg_map {
			// 			client.HIncrBy("test", key, -1)
			// 		}
					
			// 		time.Sleep(time.Millisecond * time.Duration(rand.Intn(3000)))

			// 		client.LPush("merge_queue", j)
					
			// 	}

			// } else {
			// 	client.LPush("merge_queue", j)
			// }

		} else {
			payer, remittee, money := arr[1], arr[2], arr[3]
			v_payer, err_payer := client.HGet("test", payer).Result()
			if err_payer != nil {
				fmt.Println(err_payer)
				return
			}
			// fmt.Println("v_payer: ", v_payer)
			if v_payer != "0" && v_payer != "" { // payer is locked
				// fmt.Println(payer, "is locked.")
				// lock this transaction
				RecordTimetoFile("start_improve.txt", false)
				client.LPush("block_transaction_queue", j)
				// break
			} else { // payer isn't locked, judge if remittee is locked
				token.Lock()
				v_remittee, err_remittee := client.HGet("test", remittee).Result()
				if err_remittee != nil {
					fmt.Println(err_remittee)
					token.Unlock()
					return
				} else {
					// fmt.Println("v_remittee: ", v_remittee)
					if v_remittee != "0" && v_remittee != "" { // remittee is locked
						// fmt.Println(remittee, " is locked.")
						RecordTimetoFile("start_improve.txt", false)

						client.HIncrBy("test", payer, 1)
						token.Unlock()
						result, err := contract.SubmitTransaction("transferAccount_new", payer, remittee, money)
						if err != nil {
							fmt.Printf("Failed to submit transaction_new: %s\n", err)
							RecordTimetoFile("end_improve.txt", true)
							// os.Exit(1)
						} else {
							RecordTimetoFile("end_improve.txt", false)
						}
						fmt.Println(string(result))

						client.HIncrBy("test", payer, -1)

					} else { // remittee isn't locked
						// fmt.Printf(remittee, " isn't locked.")
						RecordTimetoFile("start_improve.txt", false)

						client.HIncrBy("test", payer, 1)
						client.HIncrBy("test", remittee, 1)
						token.Unlock()
						
						result, err := contract.SubmitTransaction("transferAccount", payer, remittee, money)
						if err != nil {
							fmt.Printf("Failed to submit transaction: %s\n", err)
							RecordTimetoFile("end_improve.txt", true)
							// os.Exit(1)
						} else {
							RecordTimetoFile("end_improve.txt", false)
						}
						fmt.Println(string(result))

						client.HIncrBy("test", payer, -1)
						client.HIncrBy("test", remittee, -1)
				
						
					}
				}
			}
		}
		results <- 0
	}
}

// judge if all accounts are unlocked
// func isUnlockedForRoll(rollupTransactions_string string) bool {
// 	arr := strings.split(rollupTransactions_string, " ")
// 	arg_map := make(map[string]string)
// 	_ = json.Unmarshal([]byte(arr[1]), &arg_map)

// 	for key, _ := range arg_map {
// 		v_payer, err_payer := client.HGet("test", key).Result()
// 		if err_payer != nil {
// 			fmt.Printf("err_player")
// 			// unlocked
// 			continue
// 		} else {
// 			fmt.Println("v_payer:", v_payer)
// 			if v_payer != "0" && v_payer != "" {
// 				//locked
// 				return false
// 			} else {
// 				// unlocked
// 				continue
// 			}
// 		}
// 	}
// 	return true
// }

func RecordTimetoFile(filename string, isFail bool) {
	f, _ := os.OpenFile("./time/"+filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModeAppend|os.ModePerm)
	defer f.Close()
	if isFail {
		io.WriteString(f, time.Now().Format("2006-01-02 15:04:05")+" fail"+"\n")
	} else {
		io.WriteString(f, time.Now().Format("2006-01-02 15:04:05")+"\n")
	}
}

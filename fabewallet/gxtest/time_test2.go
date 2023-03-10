package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"log"
	"os"

	"github.com/go-redis/redis"
)

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

func main() {
	//这个设计目的只是为了防止批太大，是一个上限而已
	maxBatchSize := 10
	timeInterval := 2
	queue1 := "block_transaction_queue"
	queue2 := "merge_queue"

	logFile, err := os.OpenFile("./rollfunc.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        fmt.Println("open log file failed, err:", err)
        return
    }
    log.SetOutput(logFile)
    log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
    log.Println("这是一条很普通的日志。")

	err = initRedis()
	if err != nil {
		fmt.Printf("connect redis failed! err : %v\n", err)
		return
	}
	fmt.Println("redis连接成功！")
	
	t1 := time.NewTimer(time.Second * time.Duration(timeInterval)) 
	go func() {
		for {
			select {
			case <-t1.C:
				log.Println("timer info:",maxBatchSize," ",timeInterval," ",queue1," ",queue2)

				length, err := client.LLen(queue1).Result()
				if err != nil {
					panic(err)
				}
				s := strconv.FormatInt(length, 10)
				n, err := strconv.Atoi(s)
				if  err != nil {
					panic(err)
				}
				rollTimes := n/maxBatchSize + 1
				for i := 0; i < rollTimes; i++ {
					rollup(maxBatchSize, queue1, queue2)
				}
	
				// print 2 queues
				val1, err1 := client.LRange(queue1, 0, -1).Result()
				if err1 != nil {
					log.Printf("%s get failed! err : %v\n",queue1, err1)
					return
				}
				log.Println(queue1)
				log.Printf("data length=[%d]\n", len(val1))
	
				val1, err1 = client.LRange(queue2, 0, -1).Result()
				if err1 != nil {
					log.Printf("%s get failed! err : %v\n", queue2, err1)
					return
				}
				log.Println(queue2)
				log.Printf("data length=[%d]\n", len(val1))
	
				t1.Reset(time.Second * time.Duration(timeInterval))
			}
		}
	}()

	for {
		fmt.Println("reset info")
		fmt.Println("maxBatchSize timeInterval")
		fmt.Scanf("%d %d",&maxBatchSize,&timeInterval)
	}
	
}


func rollup(Rollup_batch_size int, queue1 string, queue2 string) {
	// !!! make sure two thread won't conflict when operating the same queue

	// cause i don't know how to use rpop, the pops code is pseudo-code
	// client.rpop("transaction_queue", function(e2,v2){
	// 	if(v2 != null) {
	// 		init(v2);
	// 	}
	// })

	// pseudo-code gx1
	// rpopn is a better choice , but i don't know if it's realized
	// make( , 0, Rollup_batch_size) is a better choice if necessary

	//优化性能，但是存疑，怕和想的不一样
	transactionFlow := make([][]string, 0, Rollup_batch_size)
	for i := 0; i < Rollup_batch_size; i++ {

		// length, err := client.LLen(queue1).Result()
		// if err != nil {
		// 	panic(err)
		// }
		// if length == 0 {
		// 	break
		// }

		v, e := client.RPop(queue1).Result()
		if e != nil {
			log.Println(queue1 + " rpop error ")
			// consider as empty
			break
		}
		if v == "" {
			log.Println(queue1 + " rpop empty")
			// consider as empty
			break
		}

		// none empty
		string_slice := strings.Split(v, " ")
		// string slice is expected as functionName payer remittee money
		slice_adjusted := []string{strconv.Itoa(i), string_slice[0], string_slice[1], "-" + string_slice[3], string_slice[2], string_slice[3]}
		transactionFlow = append(transactionFlow, slice_adjusted)
	}

	log.Println(len(transactionFlow), " transactions rolled in this call")

	//now we are talking
	//a tool go is needed
	var redundantKeySet []string
	for _, transactionI := range transactionFlow {
		redundantKeySet = append(redundantKeySet, transactionI[2], transactionI[4])
	}
	trimmedKeySet := RemoveRepByMap(redundantKeySet)
	lengthKey := len(trimmedKeySet)

	//need a function to translate keyString to index of trimmedKeySet
	var string2index map[string]int
	string2index = make(map[string]int)
	for i, str := range trimmedKeySet {
		string2index[str] = i
	}
	//string to index :  string2index
	//index to string : trimmedKeySet

	// second Step: get the union_find set
	ufs := NewUnionSet(lengthKey)
	for _, transactionI := range transactionFlow {
		ufs.Union(string2index[transactionI[2]], string2index[transactionI[4]])
	}

	// third Step:
	// { index1 : { "txIndexSet" : ["0","1","2"] , "A":"" ,"newTransaction": ["Transfer","A","","B","","C","","D",""](optional) } , index2: {} }
	// here index1 / index2 is the root of ufs, used to refer to the ufs.
	//{ index1 : { "txIndexSet" : "0,1,2" , "A":"-X", "B":"+X","Function":"Transfer"  } , index2: {} }
	//"newTransaction": ["Transfer","A","","B","","C","","D",""]
	var pkgMap map[int]map[string]string
	pkgMap = make(map[int]map[string]string)
	for index, transactionI := range transactionFlow {
		// transactionI[2] and transactionI[2] has the same root
		rootIndex := ufs.getRoot(string2index[transactionI[2]])

		if _, IsExisted := pkgMap[rootIndex]; !IsExisted {
			// init the infoMap
			var tempInfoMap map[string]string
			tempInfoMap = make(map[string]string)
			// this step is not needed, cause there is only Transfer
			tempInfoMap["Function"] = "transferAccount_multiParty"
			tempInfoMap["txIndexSet"] = ""

			pkgMap[rootIndex] = tempInfoMap
		}

		// ADD current tx to pkgMap, first record the tx index, second add or modify accounts
		if pkgMap[rootIndex]["txIndexSet"] == "" {
			pkgMap[rootIndex]["txIndexSet"] = strconv.Itoa(index)
		} else {
			pkgMap[rootIndex]["txIndexSet"] = pkgMap[rootIndex]["txIndexSet"] + "," + strconv.Itoa(index)
		}

		//a const is here
		account_parameter_place_default := []int{2, 4}
		for _, account_parameter_index := range account_parameter_place_default {
			if _, IsExisted := pkgMap[rootIndex][transactionI[account_parameter_index]]; !IsExisted {
				pkgMap[rootIndex][transactionI[account_parameter_index]] = transactionI[account_parameter_index+1]
			} else {
				num1, _ := strconv.Atoi(pkgMap[rootIndex][transactionI[account_parameter_index]])
				num2, _ := strconv.Atoi(transactionI[account_parameter_index+1])
				pkgMap[rootIndex][transactionI[account_parameter_index]] = strconv.Itoa(num1 + num2)
			}
		}
	}

	// last Step: commit
	// cause those packages are not relevant, whatever the order
	// commit the newTransaction and wait the response. If successed, continue to commit next one
	// if failed commit those orginal transactions according to txIndexSet
	//
	// think: what is the difference between normal submit and submitAsyc? should i use the latter?

	//fmt.Println(pkgMap)
	for _, ufsmap := range pkgMap {
		////string
		//str_tx_list := ufsmap["txIndexSet"]
		//delete(ufsmap, "txIndexSet")
		////cmd
		//arg_list := ""
		//arg_list = append(arg_list, ufsmap["Function"])

		arg_list_map := make(map[string]string)
		for mapkey, mapvalue := range ufsmap {
			if mapkey == "txIndexSet" || mapkey == "Function" {
				continue
			}
			//arg_list = arg_list + "," + mapkey + "," + mapvalue
			arg_list_map[mapkey] = mapvalue
		}

		//There should be only one Function here

		//must be Transfer！
		if ufsmap["Function"] != "transferAccount_multiParty" {
			log.Println("Not the right function to call the submit followed")
			return
		}

		txIndexSlice := strings.Split(ufsmap["txIndexSet"], ",")
		var note [][]string
		for _, txIndex := range txIndexSlice {
			intindex, _ := strconv.Atoi(txIndex)
			note = append(note, transactionFlow[intindex])
		}
		bytes, err := json.Marshal(note)
		if err != nil {
			panic("Serialization failed")
		}

		bytes_arglist, err := json.Marshal(arg_list_map)
		if err != nil {
			panic("Serialization failed")
		}

		log.Println(ufsmap["Function"]+":", arg_list_map)

		//resp, e := contract.SubmitTransaction(ufsmap["Function"], string(bytes_arglist), string(bytes))
		////important. error handling
		//if e != nil {
		//	fmt.Println(e, resp)
		//	//to do
		//	txIndexSlice := strings.Split(ufsmap["txIndexSet"], ",")
		//
		//	for _, txIndex := range txIndexSlice {
		//		txIndexInt, _ := strconv.Atoi(txIndex)
		//		txI := transactionFlow[txIndexInt]
		//		txOperation := txI[1]
		//		txArgs := "{\"" + txI[2] + "\":\"" + txI[3] + "\",\"" + txI[4] + "\":\"" + txI[5] + "\"}"
		//		fmt.Println(txOperation + "  " + txArgs)
		//		contract.SubmitTransaction(txOperation, txArgs, "null")
		//	}
		//}

		// //error handing is set for true
		// if enableThread {
		// 	go mySubmitHandle(contract, ufsmap["Function"], string(bytes_arglist), string(bytes), true, enableThread)
		// } else {
		// 	mySubmitHandle(contract, ufsmap["Function"], string(bytes_arglist), string(bytes), true, enableThread)
		// }

		// "function arglist note"
		rollupTransactions_string := ufsmap["Function"] + " " + string(bytes_arglist) + " " + string(bytes)

		// pseudo-code again
		// need to init the second block queue!!! gx2
		// client.LPush("block_rollUpTransaction_queue",rollupTransactions_string)
		_, err = client.LPush(queue2, rollupTransactions_string).Result()
		if err != nil {
			panic(err)
		}

	}

}

// judge if all accounts are unlocked
func isUnlockedForRoll(rollupTransactions_string string) bool {
	arg_map := make(map[string]string)
	_ = json.Unmarshal([]byte(rollupTransactions_string), &arg_map)

	for key, _ := range arg_map {
		v_payer, err_payer := client.HGet("test", key).Result()
		if err_payer != nil {
			fmt.Printf("err_player")
			// unlocked
			continue
		} else {
			fmt.Printf("v_payer:", v_payer)
			if v_payer != "0" && v_payer != "" {
				//locked
				return false
			} else {
				// unlocked
				continue
			}
		}
	}
	return true
}

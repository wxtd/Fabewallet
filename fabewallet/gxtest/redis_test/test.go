// test file

package main

import(
    "fmt"
	"github.com/go-redis/redis"
	// "strings"
)

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

func main() {
	err := initRedis()
	if err != nil {
		fmt.Printf("connect redis failed! err : %v\n",err)
		return
	}
	fmt.Println("redis连接成功！")

	v_payer, err_payer := client.HGet("test", "owner0").Result()
	fmt.Println("err: ", err_payer)
	fmt.Println(v_payer)

	// client.HIncrBy("test", "amy" , 1)

	// v, _ := client.HGet("test", "amy").Result()
	// fmt.Printf("%v\n", v);

	// val1, err1 := client.LRange("block_transaction_queue",0,-1).Result()
	// // val1, err1 := result1.val(), result1.err()
	// if err1 != nil {
	// 	fmt.Printf("transaction_queue get failed! err : %v\n",err1)
	// 	return
	// }
	// fmt.Printf("block_transaction_queue:");
	// fmt.Printf("read data from DB success. data length=[%d]\n",len(val1));
	// fmt.Printf("%v\n",val1)
	// client.LPush("transaction_queue", "ewallet transferAccount amy bob 10 ")
	// v_transaction, e_transaction := client.RPop("transaction_queue").Result()
	// if e_transaction != nil {
	// 	fmt.Printf("get transaction_queue failed! err : %v\n", e_transaction)
	// 	return
	// }
	// arr := strings.Split(v_transaction, " ")
	// fmt.Println(v_transaction)
	// // fmt.Println(arr)
	// // for k, v := range arr {
	// // 	fmt.Println(string(k) + "   " + v)
	// // }
	// chaincodeName, functionName, payer, remittee, money := arr[0], arr[1], arr[2], arr[3], arr[4]
	// fmt.Println(chaincodeName)
	// fmt.Println(functionName)
	// fmt.Println(payer)
	// fmt.Println(remittee)
	// fmt.Println(money)
	// fmt.Printf("%T\n", money)
	// client.LPush("transaction_queue", "ewallet transferAccount amy bob 10 ")


	// client.Quit()
}


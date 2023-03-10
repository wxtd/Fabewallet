// input data to redis-db

package main

import(
    "fmt"
	"github.com/go-redis/redis"
	
	"strconv"
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
	// fmt.Println("redis连接成功！")

	for i := 0; i < 20; i +=2 {
		// client.LPush("transaction_queue", "ewallet transferAccount user" + " b 1 ")
		client.LPush("transaction_queue", "ewallet transferAccount owner" + strconv.Itoa(i) + " owner" + strconv.Itoa(i+1) + " 1 ")
	}

	// val1, err1 := client.LRange("block_transaction_queue",0,-1).Result()
	// // val1, err1 := result1.val(), result1.err()
	// if err1 != nil {
	// 	fmt.Printf("transaction_queue get failed! err : %v\n",err1)
	// 	return
	// }
	// fmt.Printf("block_transaction_queue:");
	// fmt.Printf("read data from DB success. data length=[%d]\n",len(val1));
	// fmt.Printf("%v\n",val1)


	// client.Quit()
}
// get data from redis-db

package main

import(
    "fmt"
	"github.com/go-redis/redis"
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

	// transaction_queue
	val1, err1 := client.LRange("transaction_queue",0,-1).Result()
	if err1 != nil {
		fmt.Printf("transaction_queue get failed! err : %v\n",err1)
		return
	}
	fmt.Println("transaction_queue:");
	fmt.Printf("data length=[%d]\n",len(val1));
	for _, v := range val1 {
		fmt.Println(v)
	}
	fmt.Println()

	// block_transaction_queue
	val1, err1 = client.LRange("block_transaction_queue",0,-1).Result()
	if err1 != nil {
		fmt.Printf("block_transaction_queue get failed! err : %v\n",err1)
		return
	}
	fmt.Println("block_transaction_queue:");
	fmt.Printf("data length=[%d]\n",len(val1));
	for _, v := range val1 {
		fmt.Println(v)
	}
	fmt.Println()

	// hashmap
	v_hash, e_hash := client.HGetAll("test").Result()
	if e_hash != nil {
		fmt.Printf("hashmap get failed! err : %v\n",e_hash)
		return
	}
	fmt.Println("hash:");
	fmt.Printf("data length=[%d]\n",len(v_hash));
	fmt.Println(v_hash)
	fmt.Println()

	// client.Quit()
}
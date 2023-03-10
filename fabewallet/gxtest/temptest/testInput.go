package inputTest

import (
	"fmt"
)

func main() {
	maxBatchSize := 10
	timeInterval := 2

	for {
		fmt.Println("reset info")
		fmt.Println("maxBatchSize timeInterval")
		fmt.Scanf("%d %d",&maxBatchSize,&timeInterval)

		fmt.Println(maxBatchSize)
		fmt.Println(timeInterval)
	}
}
package content

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println(int64(math.Ceil((float64(42)/float64(41)))))
}


//tick := time.Tick(time.Second)
//for count :=10 ; count > 0; count-- {
//fmt.Println(count)
//fmt.Println(tick)
//}


//select {
//case <-ch1:
//fmt.Println("ch1 ch1")
//case x := <-ch2:
//fmt.Printf("default %s\n ",x)
//case ch3 <- y:
//fmt.Println("ch3 ch3")
//default:
//fmt.Printf("default %s\n ","1234")
//}
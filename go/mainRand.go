package main

import (
	"math/rand"
	"time"
	"fmt"
)

func main()  {
	seed	:=	time.Now().UTC().UnixNano()
	fmt.Print(seed)
	rng	:=	rand.New(rand.NewSource(seed))
	a := randomPalindrome(rng)
	fmt.Print(a)
}

func randomPalindrome(rng *rand.Rand) string {
	n := rng.Intn(25)
	runes := make([]rune, n)
	for i := 0; i < (n+1)/2; i++ {
		r := rune(rng.Intn(0x1000))
		runes[i] = r
		runes[n-1-i] = r
	}
	return string(runes)
}
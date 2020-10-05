package main

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	// "time"
)

func main() {
	var (
		userChan    = make(chan User, 1)
		group       = new(errgroup.Group)
		studentChan = make(chan Student, 1)
	)
	group.Go(func() error {
		userChan <- User{
			Name: "赵贷贷",
		}
		close(userChan)
		return nil
	})

	group.Go(func() error {
		studentChan <- Student{
			Name: "赵桥桥",
		}
		close(studentChan)
		return nil
	})
	if err := group.Wait(); err != nil {
		fmt.Println("Get errors: ", err)
	} else {
		fmt.Println("Get all num successfully!")
	}

	person := Person{
		User:    <-userChan,
		Student: <-studentChan,
	}
	fmt.Printf(" 人类 1 =>  %#v\n ", person)

	// for v := range userChan {
	// 	fmt.Printf(" 测试=>  %#v\n", v)
	// }

	// for {
	// 	select {
	// 	case user, ok := <-userChan:
	// 		fmt.Printf(" 测试  %#v\n", ok)
	// 		if !ok {
	// 			return
	// 		}
	// 		fmt.Printf(" 赵赵  %#v\n", user)
	// 	}
	// }

	// ch := make(chan int, 2)
	// go func() {
	// 	fmt.Println("Hello inline")
	// 	// send a value on channel
	// 	ch <- 1
	// }()
	// // call a function as goroutine
	// go printHello(ch)
	// fmt.Println("Hello from main")
	// i := <-ch
	// fmt.Println("Recieved ", i)
	// // time.Sleep(2*time.Second)
	// close(ch)
	// b := <-ch
	// fmt.Println("Recievedb", b)
}

type User struct {
	Name string
}

type Student struct {
	Name string
}

type Person struct {
	User
	Student
}

// prints to stdout and puts an int on channel
func printHello(ch chan int) {
	fmt.Println("Hello from printHello")
	// send a value on channel
	ch <- 2
}

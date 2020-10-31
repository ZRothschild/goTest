package main

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"sync"
	"sync/atomic"
)

// Lock try lock
type Lock struct {
	c chan struct{}
}

// NewLock generate a try lock
func NewLock() Lock {
	var l Lock
	l.c = make(chan struct{}, 1)
	l.c <- struct{}{}
	return l
}

// Lock try lock, return lock result
func (l Lock) Lock() bool {
	lockResult := false
	select {
	case <-l.c:
		lockResult = true
	default:
	}
	return lockResult
}

// Unlock , Unlock the try lock
func (l Lock) Unlock() {
	l.c <- struct{}{}
}

func main() {
	var counter int64
	var l = NewLock()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if !l.Lock() {
				// log error
				println("lock failed")
				return
			}
			counter++
			println("current counter", counter)
			l.Unlock()
		}()
	}
	wg.Wait()

	return

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// l.RLock()
			// counter++
			atomic.AddInt64(&counter, 1)
			// l.RUnlock()
		}()
	}

	wg.Wait()
	println(counter)

	return
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

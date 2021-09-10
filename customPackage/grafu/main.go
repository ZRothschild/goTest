package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var wg = &sync.WaitGroup{}

func main() {
	wg.Add(2)
	fmt.Println("pid: ", os.Getpid())

	process ,err := os.FindProcess(os.Getpid())
	if err != nil {
		return
	}

	err = process.Release()
	if err != nil {
		return
	}

	go func() {

		fmt.Println("pid: ", os.Getpid())
		c1 := make(chan os.Signal, 1)
		signal.Notify(c1, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
		fmt.Printf("goroutine 1 receive a signal : %v\n\n", <-c1)
		//path := os.Args[0]
		//var args []string
		//if len(os.Args) > 1 {
		//	args = os.Args[1:]
		//}
		//fmt.Println(path, "path", args, "args")
		//// 程序启动程序
		//cmd := exec.Command(path, args...)
		//cmd.Stdout = os.Stdout
		//cmd.Stderr = os.Stderr
		//cmd.ExtraFiles = nil
		//cmd.Env = nil
		//
		//if err := cmd.Start(); err != nil {
		//	log.Fatalf("Restart: Failed to launch, error: %v", err)
		//}
		//fmt.Println("pid: ", os.Getpid())
		wg.Done()
	}()

	wg.Wait()
	fmt.Printf("all groutine done!\n")
}

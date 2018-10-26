package test

import (
"net"
"os"
"fmt"
)

func main()  {
	conn ,err := net.Dial("tcp",":8080")
	defer conn.Close()
	if err != nil {
		fmt.Printf("net.Dial %v\n",err)
		os.Exit(1)
	}
	_,err = conn.Write([]byte("hello my firend"))
	if err != nil {
		fmt.Printf("conn.Read %v\n",err)
		os.Exit(2)
	}
}
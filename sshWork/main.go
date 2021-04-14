package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"runtime"
	"sync"
	"time"
)

type (
	HostInfo struct {
		host   string
		port   string
		user   string
		pass   string
		isWeak bool
	}

	List struct {
		ip   chan []string
		user chan []string
		pwd  chan []string
		port chan []string
		err
	}
)

// read lime from file and Scan
func Prepare(ipDict, userDict, pwdDict, portDict string) (ip, user, pwd, port []string, err error) {
	ipDictFile, err := os.Open(ipDict)
	if err != nil {
		return ip, user, pwd, port, err
	}
	defer ipDictFile.Close()
	scanner := bufio.NewScanner(ipDictFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		ip = append(ip, scanner.Text())
	}
	userDictFile, err := os.Open(userDict)
	if err != nil {
		return ip, user, pwd, port, err
	}
	defer userDictFile.Close()
	scanner = bufio.NewScanner(userDictFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		user = append(user, scanner.Text())
	}
	pwdDictFile, err := os.Open(pwdDict)
	if err != nil {
		return ip, user, pwd, port, err
	}
	defer pwdDictFile.Close()
	scanner = bufio.NewScanner(pwdDictFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		pwd = append(pwd, scanner.Text())
	}

	portDictFile, err := os.Open(portDict)
	if err != nil {
		return ip, user, pwd, port, err
	}
	defer portDictFile.Close()
	scanner = bufio.NewScanner(portDictFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		port = append(port, scanner.Text())
	}
	return ip, user, pwd, port, err
}

// Scan function
func Scan(sliceIpList, sliceUser, slicePass, slicePort []string) (err error) {
	var (
		hostInfo HostInfo
		client   *ssh.Client
		// total       = len(sliceIpList) * len(sliceUser) * len(slicePass)
		sucChanHost = make(chan HostInfo, 999)
	)
	go func() {
		var wg sync.WaitGroup
		for _, host := range sliceIpList {
			for _, user := range sliceUser {
				for _, pwd := range slicePass {
					for _, port := range slicePort {
						hostInfo = HostInfo{
							host:   host,
							port:   port,
							user:   user,
							pass:   pwd,
							isWeak: false,
						}
						wg.Add(1)
						go func(hostInfo HostInfo) {
							url := hostInfo.host + ":" + hostInfo.port + " " + hostInfo.user + " " + hostInfo.pass
							if err = crack(hostInfo, client); err != nil {
								fmt.Printf("\033[1;31;40m %v %s \033[0m\n", hostInfo.isWeak, url)
								wg.Done()
								return
							}
							hostInfo.isWeak = true
							fmt.Printf("\033[1;30;47m %v %s\033[0m\n", hostInfo.isWeak, url)
							sucChanHost <- hostInfo
							wg.Done()
						}(hostInfo)
						for runtime.NumGoroutine() > runtime.NumCPU()*10000 {
							time.Sleep(2 * time.Microsecond)
						}
					}
				}
			}
		}
		wg.Wait()
		close(sucChanHost)
	}()
	fileSuc, err := os.Create("./suc.txt")
	if err != nil {
		return err
	}
	defer fileSuc.Close()
	for v := range sucChanHost {
		url := v.host + ":" + v.port + " " + v.user + " " + v.pass + "\n"
		if _, err = fileSuc.WriteString(url); err != nil {
			return err
		}
	}
	return nil
}

func crack(hostInfo HostInfo, client *ssh.Client) (err error) {
	var (
		pass = []ssh.AuthMethod{ssh.Password(hostInfo.pass)}
		conf = ssh.ClientConfig{
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			User:            hostInfo.user,
			Auth:            pass,
			Timeout:         400 * time.Millisecond,
		}
	)
	if client, err = ssh.Dial("tcp", hostInfo.host+":"+hostInfo.port, &conf); err != nil {
		return err
	}
	defer client.Close()
	return nil
}

// main function
func main() {
	var (
		now                                 = time.Now()
		ipDict, userDict, pwdDict, portDict = "./sshWork/IPCp.txt", "./sshWork/userCP.txt", "./sshWork/passwdCp.txt", "./sshWork/portCp.txt"
	)
	runtime.GOMAXPROCS(runtime.NumCPU())
	ip, user, pwd, port, err := Prepare(ipDict, userDict, pwdDict, portDict)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := Scan(ip, user, pwd, port); err != nil {
		fmt.Println(err)
		return
	}
	gap := time.Now().Sub(now).Seconds()
	fmt.Printf("\033[1;33;46m 结束了总耗时 %f 秒\033[0m\n", gap)
	// hostInfo := HostInfo{
	// 	host: "123.150.189.157",
	// 	port: "55667",
	// 	user: "root",
	// 	pass: "QG2geJzd8MBX7uWq",
	// }
	//
	// // sucChanHost := make(chan HostInfo, 10)
	// if err := crack(hostInfo); err != nil {
	// 	fmt.Println(err)
	// }
	// runtime.GOMAXPROCS(runtime.NumCPU())
	// Scan(Prepare(ipList, userDict, passDict))
}

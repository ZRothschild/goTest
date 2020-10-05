package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type safeMap struct {
	m map[string]string
	l *sync.RWMutex
}

func (s *safeMap) Set(key string, value string) {
	s.l.Lock()
	defer s.l.Unlock()
	s.m[key] = value
}

func (s *safeMap) Get(key string) string {
	s.l.RLock()
	defer s.l.RUnlock()
	return s.m[key]
}

func newSafeMap() *safeMap {
	return &safeMap{l: new(sync.RWMutex), m: make(map[string]string)}
}

func main() {
	var syn sync.WaitGroup
	confPath := "./"
	ViperConf := viper.New()
	ViperConf.SetConfigName("config") // 配置文件名称
	ViperConf.AddConfigPath(confPath) // 设置配置文件的搜索目录
	ViperConf.SetConfigType("toml")
	if err := ViperConf.ReadInConfig(); err != nil { // 加载配置文件内容
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	go ReloadConfig(ViperConf)

	syn.Add(1)

	syn.Wait()

	sy := newSafeMap()

	f, err := os.OpenFile("./text.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	logger := log.New(f, "prefix ", log.LstdFlags)

	for i := 0; i < 999999; i++ {
		syn.Add(1)
		go func(i int) {
			ii := strconv.FormatInt(int64(i), 10)
			sy.Set(ii, ii)
			syn.Done()
		}(i)
	}
	syn.Wait()

	time.Sleep(20000000000)

	// for i := 0;  i < 999999 ; i++  {
	// 	ii := strconv.FormatInt(int64(i),10)
	// 	ii = sy.Get(ii)
	// 	logger.Println(i, "  ii =>  ",ii)
	// }
	ii := sy.Get("999998")
	logger.Println(999998, "  ii =>  ", ii)

	// time.Sleep(10000000000)
}

// 热加载
func ReloadConfig(ViperCon *viper.Viper) {
	ViperCon.WatchConfig()
	ViperCon.OnConfigChange(func(in fsnotify.Event) {
		fmt.Printf("Detect config change: %s \n", in.String())
	})
}

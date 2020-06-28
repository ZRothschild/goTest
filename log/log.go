package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type Data struct {
	Id    int64
	Value int64
}

var Mock = []Data{
	{Id: 1, Value: 1}, {Id: 1, Value: 2}, {Id: 1, Value: 3}, {Id: 1, Value: 4}, {Id: 1, Value: 5}, {Id: 1, Value: 6}, {Id: 1, Value: 7}, {Id: 1, Value: 8},
	{Id: 2, Value: 1}, {Id: 2, Value: 2}, {Id: 2, Value: 3}, {Id: 2, Value: 4}, {Id: 2, Value: 5}, {Id: 2, Value: 6}, {Id: 2, Value: 7}, {Id: 2, Value: 8},
	{Id: 3, Value: 1}, {Id: 3, Value: 2}, {Id: 3, Value: 3}, {Id: 3, Value: 4}, {Id: 3, Value: 5}, {Id: 3, Value: 6}, {Id: 3, Value: 7}, {Id: 3, Value: 8},
	{Id: 4, Value: 1}, {Id: 4, Value: 2}, {Id: 4, Value: 3}, {Id: 4, Value: 4}, {Id: 4, Value: 5}, {Id: 4, Value: 6}, {Id: 4, Value: 7}, {Id: 4, Value: 8},
	{Id: 4, Value: 2}, {Id: 5, Value: 2}, {Id: 1, Value: 3}, {Id: 4, Value: 12}, {Id: 4, Value: 5}, {Id: 1, Value: 6}, {Id: 4, Value: 7}, {Id: 4, Value: 8},
	{Id: 4, Value: 3}, {Id: 6, Value: 2}, {Id: 2, Value: 3}, {Id: 4, Value: 13}, {Id: 1, Value: 5}, {Id: 3, Value: 6}, {Id: 1, Value: 7}, {Id: 7, Value: 8},
	{Id: 4, Value: 4}, {Id: 7, Value: 2}, {Id: 3, Value: 3}, {Id: 4, Value: 14}, {Id: 3, Value: 5}, {Id: 2, Value: 6}, {Id: 2, Value: 7}, {Id: 5, Value: 8},
	{Id: 4, Value: 5}, {Id: 8, Value: 2}, {Id: 4, Value: 3}, {Id: 4, Value: 15}, {Id: 2, Value: 5}, {Id: 4, Value: 6}, {Id: 3, Value: 7}, {Id: 6, Value: 8},
}

// 自制的简易并发安全字典
type ConcurrentMap struct {
	m  map[int64][]int64
	mu sync.RWMutex
}

func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{
		m: make(map[int64][]int64),
	}
}

func (cMap *ConcurrentMap) Delete(key int64) {
	cMap.mu.Lock()
	defer cMap.mu.Unlock()
	delete(cMap.m, key)
}

func (cMap *ConcurrentMap) Get() map[int64][]int64 {
	cMap.mu.RLock()
	defer cMap.mu.RUnlock()
	return cMap.m
}

func (cMap *ConcurrentMap) Load(key int64) (value []int64, ok bool) {
	cMap.mu.RLock()
	defer cMap.mu.RUnlock()
	value, ok = cMap.m[key]
	return
}

func (cMap *ConcurrentMap) SetOrStore(key int64, value int64) map[int64][]int64 {
	cMap.mu.Lock()
	defer cMap.mu.Unlock()
	actual, loaded := cMap.m[key]
	if loaded {
		actual = append(actual, value)
		cMap.m[key] = actual
		return cMap.m
	}
	var tmp []int64
	tmp = append(tmp, value)
	cMap.m[key] = tmp
	return cMap.m
}

func (cMap *ConcurrentMap) LoadOrStore(key int64, value []int64) (actual []int64, loaded bool) {
	cMap.mu.Lock()
	defer cMap.mu.Unlock()
	actual, loaded = cMap.m[key]
	if loaded {
		return
	}
	cMap.m[key] = value
	actual = value
	return
}

func main() {
	//O_APPEND 添加写  O_CREATE 不存在则生成   O_WRONLY 只写模式
	f, err := os.OpenFile("./text.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	//时间展示格式  LstdFlags
	logger := log.New(f, "prefix ", log.LstdFlags)
	logger.Println("text to append")
	logger.Println("more text to append")

	var cM = NewConcurrentMap()
	for _, v := range Mock {
		go func(v Data) {
			cM.SetOrStore(v.Id, v.Value)
		}(v)
	}

	va := cM.Get()
	fmt.Printf("==== %#v\n", va)

}

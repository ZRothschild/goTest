package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"sync"
)

const (
	NULL = ""
)

type BodyMap map[string]interface{}

type xmlMapMarshal struct {
	XMLName xml.Name
	Value   interface{} `xml:",cdata"`
}

type xmlMapUnmarshal struct {
	XMLName xml.Name
	Value   string `xml:",cdata"`
}

var mu = new(sync.RWMutex)

// 设置参数
func (bm BodyMap) Set(key string, value interface{}) BodyMap {
	mu.Lock()
	bm[key] = value
	mu.Unlock()
	return bm
}

func (bm BodyMap) SetBodyMap(key string, value func(bm BodyMap)) BodyMap {
	_bm := make(BodyMap)
	value(_bm)

	mu.Lock()
	bm[key] = _bm
	mu.Unlock()
	return bm
}

// 设置 FormFile
func (bm BodyMap) SetFormFile(fieldName string, filePath string) (err error) {
	_FileBm := make(BodyMap)
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("bm.SetFormFile(%s, %s),err:%w", fieldName, filePath, err)
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("bm.SetFormFile(%s, %s),err:%w", fieldName, filePath, err)
	}
	fileContent := make([]byte, stat.Size())
	_, err = file.Read(fileContent)
	if err != nil {
		return fmt.Errorf("bm.SetFormFile(%s, %s),err:%w", fieldName, filePath, err)
	}
	_FileBm[stat.Name()] = fileContent

	mu.Lock()
	bm[fieldName] = _FileBm
	mu.Unlock()
	return nil
}

// 获取参数，同 GetString()
func (bm BodyMap) Get(key string) string {
	return bm.GetString(key)
}

// 获取参数转换string
func (bm BodyMap) GetString(key string) string {
	if bm == nil {
		return NULL
	}
	mu.RLock()
	defer mu.RUnlock()
	value, ok := bm[key]
	if !ok {
		return NULL
	}
	v, ok := value.(string)
	if !ok {
		return convertToString(value)
	}
	return v
}

// 获取原始参数
func (bm BodyMap) GetInterface(key string) interface{} {
	if bm == nil {
		return nil
	}
	mu.RLock()
	defer mu.RUnlock()
	return bm[key]
}

// 删除参数
func (bm BodyMap) Remove(key string) {
	mu.Lock()
	delete(bm, key)
	mu.Unlock()
}

// 置空BodyMap
func (bm BodyMap) Reset() {
	mu.Lock()
	for k := range bm {
		delete(bm, k)
	}
	mu.Unlock()
}

func (bm BodyMap) JsonBody() (jb string) {
	mu.Lock()
	defer mu.Unlock()
	bs, err := json.Marshal(bm)
	if err != nil {
		return ""
	}
	jb = string(bs)
	return jb
}

func (bm BodyMap) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	if len(bm) == 0 {
		return nil
	}
	start.Name = xml.Name{NULL, "xml"}
	if err = e.EncodeToken(start); err != nil {
		return
	}
	for k := range bm {
		if v := bm.GetString(k); v != NULL {
			e.Encode(xmlMapMarshal{XMLName: xml.Name{Local: "xxxx"}, Value: v})
		}
	}
	return e.EncodeToken(start.End())
}

func (bm *BodyMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	for {
		var e xmlMapUnmarshal
		err = d.Decode(&e)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		bm.Set(e.XMLName.Local, e.Value)
	}
}

func convertToString(v interface{}) (str string) {
	if v == nil {
		return NULL
	}
	var (
		bs  []byte
		err error
	)
	if bs, err = json.Marshal(v); err != nil {
		return NULL
	}
	str = string(bs)
	return
}

type Name struct {
	Name string
	F    CC
}
type CC struct {
	Age int64
}

func main() {
	n := Name{
		Name: "xxx",
		F: CC{
			Age: 1,
		},
	}
	b, err := xml.Marshal(n)
	fmt.Println(string(b), err)
}

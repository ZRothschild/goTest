package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/tidwall/gjson"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type Handle func(gKey, gValue *gjson.Result)

type Handler struct {
	FieldNames []string   // 需要处理的字段集合
	FieldKind  gjson.Type // 字段的类型
	HandleList []Handle   // 处理方法集合
}

// TrimSpaceHandle 去除空格符号
func TrimSpaceHandle(gKey, gValue *gjson.Result) {
	if gValue.Type == gjson.String {
		gValue.Str = strings.TrimSpace(gValue.Str)
	}
}

// DeepModifyWrap 深度修改参数函数
func DeepModifyWrap(bt []byte, handlers []Handler) interface{} {
	var (
		btGJson   = gjson.ParseBytes(bt)
		handleMap = make(map[string][]Handle)
	)

	if btGJson.Type == gjson.Null {
		return nil
	}
	// 组装过滤函数,循环后可以直接使用key获取到对应的处理函数
	for _, handler := range handlers {
		for _, fieldName := range handler.FieldNames {
			flag := fieldName + "_#_" + strconv.Itoa(int(handler.FieldKind))
			handleMap[flag] = handler.HandleList
		}
	}
	// 递归处理过滤函数
	return DeepModify(btGJson, handleMap)
}

// DeepModify 递归处理 filters 里面的每一个设置字段 且使用对应的自定义函数
func DeepModify(gResult gjson.Result, handleMap map[string][]Handle) interface{} {

	if gResult.IsArray() {
		jsonSlice := make([]interface{}, 0)
		for _, data := range gResult.Array() {
			jsonSlice = append(jsonSlice, DeepModify(data, handleMap))
		}
		return jsonSlice
	}

	mapJson := make(map[string]interface{})
	gResult.ForEach(func(gKey, gValue gjson.Result) bool {
		var (
			k    = gKey.String()
			flag = k + "_#_" + strconv.Itoa(int(gValue.Type))

			handles, ok = handleMap[flag]
		)

		if ok {
			for _, handle := range handles {
				// 链式执行 注意 k,v 变化
				handle(&gKey, &gValue)
			}
		}
		
		if gValue.IsArray() {
			jsonSlice := make([]interface{}, 0)
			for _, data := range gValue.Array() {
				jsonSlice = append(jsonSlice, DeepModify(data, handleMap))
			}
			mapJson[k] = jsonSlice
		} else if gValue.IsObject() {
			mapJson[k] = DeepModify(gValue, handleMap)
		} else {
			mapJson[k] = gValue.Value()
		}
		return true
	})

	return mapJson
}

func main() {
	jsonStr := `{
"data": {
    "businessId": 17831,
    "pageIndex": 1,
    "pageSize": 2000,
    "relation": 1,
    "accountId"  :     " us-be1ba459df3b4 986a43092dbdde242ea ",
    "serverId": "{\"errorCode\":\"000001005\",\"accountId\":\"resource not found: 'rsp-5a152356fd1d46f18f8e8b9860e1e775'\",\"target\":\"\",\"details\":null}",
    "regionId": "1",
    "approveBusinessRoleInfoList": [
        {
            "accountId": " us-be1ba459df3b4 986a43092dbdde242ea ",
            "userId": "us-09b344b175a340ec81e1a69855c64fa5",
            "requestName": "INTLProject - us-09b344b175a340ec81e1a69855c64fa5 - INTLProject-测试专用",
            "requestId": "49842",
            "approveFlag": "Y",
            "note": " "
        }
    ]
}
}`

	//	jsonStr = `[
	//{
	//"name": "zhangsan",
	//"age": "10",
	//"phone": "11111",
	//"email": "11111@11.com"
	//},
	//{
	//"name": "lisi",
	//"age": "20",
	//"phone": "22222",
	//"email": "22222@22.com"
	//},
	//]`
	//jsonStr = ``
	//jsonStr = `dfsfsf`
	filters := []Handler{
		{
			FieldNames: []string{"accountId"},
			FieldKind:  gjson.String,
			HandleList: []Handle{TrimSpaceHandle},
		},
	}
	data := DeepModifyWrap([]byte(jsonStr), filters)

	bt, err := json.Marshal(data)

	fmt.Println(string(bt), err)

	//age.Raw
	//fmt.Println(age,age.Raw)
	//DeepModify(age)

	return
	var t interface{}
	j := 1
	t = j

	fmt.Println(t)

	accaa := []int{1, 2, 5, 3, 4}
	fmt.Println(accaa) // [1 2 5 3 4]
	sort.Sort(sort.IntSlice(accaa))
	fmt.Println(accaa)
	sort.Sort(sort.Reverse(sort.IntSlice(accaa)))
	fmt.Println(accaa) // [5 4 3 2 1]

	return

	dataByte := []byte("")
	dataTrr := []byte(`"accountId"`)
	aa := bytes.Index(dataByte, dataTrr)

	a := len(dataTrr)

	aStr := string(dataByte[aa : aa+a])

	affer := dataByte[aa+a:]
	cc := []byte(":")

	acccr := string(affer)
	fmt.Println(acccr)
	aindex := bytes.Index(affer, cc)

	add := dataByte[aa : aa+a+aindex+1]

	adddStr := string(add)

	maoHao := dataByte[aa+a+aindex+1:]
	dd := []byte(`"`)
	aindexa := bytes.Index(maoHao, dd)
	maoHaocc := dataByte[aa+a+aindex+1+aindexa+1:]

	fmt.Println(aa, aStr, adddStr, maoHaocc)

	return

	var s interface{}
	//err := json.Unmarshal([]byte(data), &s)
	//fmt.Printf("%+#v\n", s)

	sType := reflect.ValueOf(s)
	if sType.Kind() == reflect.Map {
		mapIter := sType.MapRange()
		for mapIter.Next() {
			if mapIter.Key().Kind() == reflect.String && mapIter.Key().String() == "accountId" {
				if mapIter.Value().Kind() == reflect.Interface {
					i := mapIter.Value().Interface()
					iValueOf := reflect.ValueOf(i)
					if iValueOf.Kind() == reflect.String {
						aa := iValueOf.String()
						testAA := strings.TrimSpace(aa)
						aaValue := reflect.ValueOf(testAA)
						sType.SetMapIndex(mapIter.Key(), aaValue)
					}
				}
				//fmt.Println(mapIter.Value().CanConvert())
			} else if mapIter.Key().Kind() == reflect.String && mapIter.Key().String() == "relation" {
				if mapIter.Value().Kind() == reflect.Interface {
					i := mapIter.Value().Interface()
					iValueOf := reflect.ValueOf(i)
					fmt.Println(iValueOf.Kind(), "relaiton")
				}
			}

		}

	} else {
		fmt.Println("test aaaaaaaaaaa")
	}

	fmt.Printf("end %+#v\n", s)
	return
	app := iris.Default()

	app.Get("/ping", func(ctx iris.Context) {
		ctx.GetBody()
		_, _ = ctx.JSON(iris.Map{
			"message": "pong",
		})
	})
	// listen and serve on http://0.0.0.0:8080.
	_ = app.Run(iris.Addr(":8080"))

	root := Root{Id: 33, Name: "name", Age: 10}
	bXml, err := xml.Marshal(root)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%#v\n", string(bXml))

	d := xml.NewDecoder(strings.NewReader(testInput))

	tk, err := d.Token()
	fmt.Printf("%#v\n", tk)

	start := xml.StartElement{Name: xml.Name{Local: "test"}}
	//a, err := xml.NewEncoder(nil)
	//
	//fmt.Printf("%#v\n", string(a))
	fmt.Printf("%#v\n", start)
}

const testInput = `<?xml version="1.0" encoding="UTF-8"?>`

type Root struct {
	XMLName xml.Name `xml:"person"`
	Id      int      `xml:"id,attr"`
	Name    string   `xml:",Value"`
	Age     int64    `xml:",cdata"`
}

func ValueTypeKind(int) {

}

func TestData(aa interface{}) (err error) {

	switch aa.(type) {
	case string:
		fmt.Println("cccccc")
	case interface{}:
		fmt.Println("interface")
	default:
		fmt.Println("stir")
	}

	return err
	defer func() {
		if err != nil {
			fmt.Println("xxxxxxxxx", err)
		}
	}()

	{
		err, aa = errors.New("cuowu"), "xxx"
		fmt.Println(err, aa)
	}

	fmt.Println("======", err)
	return err
}

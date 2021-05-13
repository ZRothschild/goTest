package main

import (
	"encoding/xml"
	"fmt"
	"strings"
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.Default()


	app.Get("/ping", func(ctx iris.Context) {
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
	a, err := xml.NewEncoder()

	fmt.Printf("%#v\n", string(a))
	fmt.Printf("%#v\n", start)
}

const testInput = `<?xml version="1.0" encoding="UTF-8"?>`

type Root struct {
	XMLName xml.Name `xml:"person"`
	Id      int      `xml:"id,attr"`
	Name    string   `xml:",Value"`
	Age     int64    `xml:",cdata"`
}

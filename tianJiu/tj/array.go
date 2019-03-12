/*
	数据结构设置
*/
package tj

const (
	wen string = "文子"
	wu  string = "武子"
	g   string = "宫"
	m   string = "么"
	d   string = "点"
)

var Child = [2]string{wu, wen}
var Attributes = [3]string{d, m, g}

type Tj struct {
	Name       string //牌名称
	Number     int8   //牌点数
	Child      int8   //0:武子 1:文子
	Attributes int8   // 0:点  1:么 2:宫
	Level      int8   //对大小
	Flag       int8   //单只大小
}

var Tjs = [32]Tj{
	{Name: "天", Number: 12, Child: 1, Attributes: 2, Level: 15, Flag: 16},
	{Name: "天", Number: 12, Child: 1, Attributes: 2, Level: 15, Flag: 16},
	{Name: "地", Number: 2, Child: 1, Attributes: 2, Level: 14, Flag: 15},
	{Name: "地", Number: 2, Child: 1, Attributes: 2, Level: 14, Flag: 15},
	{Name: "人", Number: 8, Child: 1, Attributes: 2, Level: 13, Flag: 14},
	{Name: "人", Number: 8, Child: 1, Attributes: 2, Level: 13, Flag: 14},
	{Name: "和", Number: 4, Child: 1, Attributes: 2, Level: 12, Flag: 13},
	{Name: "和", Number: 4, Child: 1, Attributes: 2, Level: 12, Flag: 13},
	{Name: "梅", Number: 10, Child: 1, Attributes: 2, Level: 11, Flag: 12},
	{Name: "梅", Number: 10, Child: 1, Attributes: 2, Level: 11, Flag: 12},
	{Name: "长", Number: 6, Child: 1, Attributes: 2, Level: 10, Flag: 11},
	{Name: "长", Number: 6, Child: 1, Attributes: 2, Level: 10, Flag: 11},
	{Name: "板", Number: 4, Child: 1, Attributes: 2, Level: 9, Flag: 10},
	{Name: "板", Number: 4, Child: 1, Attributes: 2, Level: 9, Flag: 10},
	{Name: "斧头", Number: 11, Child: 1, Attributes: 1, Level: 8, Flag: 9},
	{Name: "斧头", Number: 11, Child: 1, Attributes: 1, Level: 8, Flag: 9},
	{Name: "四六", Number: 10, Child: 1, Attributes: 1, Level: 7, Flag: 8},
	{Name: "四六", Number: 10, Child: 1, Attributes: 1, Level: 7, Flag: 8},
	{Name: "么六", Number: 7, Child: 1, Attributes: 1, Level: 6, Flag: 7},
	{Name: "么六", Number: 7, Child: 1, Attributes: 1, Level: 6, Flag: 7},
	{Name: "么五", Number: 6, Child: 1, Attributes: 1, Level: 5, Flag: 6},
	{Name: "么五", Number: 6, Child: 1, Attributes: 1, Level: 5, Flag: 6},
	{Name: "红九", Number: 9, Child: 0, Attributes: 0, Level: 4, Flag: 5},
	{Name: "黑九", Number: 9, Child: 0, Attributes: 0, Level: 4, Flag: 5},
	{Name: "弯八", Number: 8, Child: 0, Attributes: 0, Level: 3, Flag: 4},
	{Name: "平八", Number: 8, Child: 0, Attributes: 0, Level: 3, Flag: 4},
	{Name: "红七", Number: 7, Child: 0, Attributes: 0, Level: 2, Flag: 3},
	{Name: "黑七", Number: 7, Child: 0, Attributes: 0, Level: 2, Flag: 3},
	{Name: "红六点", Number: 6, Child: 0, Attributes: 0, Level: 16, Flag: 2},
	{Name: "红五", Number: 5, Child: 0, Attributes: 0, Level: 1, Flag: 1},
	{Name: "黑五", Number: 5, Child: 0, Attributes: 0, Level: 1, Flag: 1},
	{Name: "红三点", Number: 3, Child: 0, Attributes: 0, Level: 16, Flag: 0},
}

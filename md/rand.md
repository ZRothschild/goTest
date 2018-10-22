## go语言生成自定义组合随机字符串

> go语言生成自定义组合随机字符串包括（纯数字，小写字母，大写字母，数字小大写字母特殊字符,
> 数字加小写字母,数字加大写字母，大小写字母，数字大小写字母）

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)
func main() {
	//获取可能包含数字大小写字母长度为4的字符串
	//kind 小于零返回可能包含数字小大写字母特殊字符
	//kind 大于7返回可能包含数字大小写字母字符
	by := kindRand(4,7)
	fmt.Printf("rand %v\n",string(by))
}

/**
 * len 随机码的位数
 * kind 0    // 纯数字
 *      1    // 小写字母
 *      2    // 大写字母
 *		3    // 数字小大写字母特殊字符
 *      4    // 数字、小写字母
 *      5    // 数字、大写字母
 *		6    // 大小写字母
 *		7    // 数字大小写字母
*/
func kindRand(len int, kind int) string {
	//数字0ASCII值为48， 字符a ASCII值为97  字符A ASCII值为65
	kindCp,kinds, result := kind,[][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65},[]int{93,33}}, make([]byte, len)
	//满足两种要求所有要随机找其中一种随机字符
	randMap := map[int][2]int{4:[2]int{0,1}, 5:[2]int{0,2}, 6:[2]int{1,2}}
	for i := 0; i < len; i++ {
		//随机动态种子
		rand.Seed(time.Now().UnixNano())
		if kind < 0  {
			kindCp = 3
		} else if kind > 3 && kind < 7  {
			//rand.Intn 随机值是左闭右开区间
			kindGap := rand.Intn(2)
			kindCp = randMap[kind][kindGap]
		}else if kind >= 7 {
			kindCp = rand.Intn(3)
		}
		scope, base := kinds[kindCp][0], kinds[kindCp][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return string(result)
}
```


package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"golang.org/x/crypto/ripemd160"
)

// s := "string1"
//产生一个散列值得方式是 sha1.New()，sha1.Write(bytes)，然后 sha1.Sum([]byte{})。这里我们从一个新的散列开始。
// h := sha1.New() // md5加密类似md5.New()
//写入要处理的字节。如果是一个字符串，需要使用[]byte(s) 来强制转换成字节数组。
// h.Write([]byte(s))
//这个用来得到最终的散列值的字符切片。Sum 的参数可以用来对现有的字符切片追加额外的字节切片：一般不需要要。也就是前面加固定值比如比特币是btc 开头
// bs := h.Sum(nil)
//SHA1 值经常以 16 进制输出，使用%x 来将散列结果格式化为 16 进制字符串。
// fmt.Printf("%x\n", bs)
//如果需要对另一个字符串加密，要么重新生成一个新的散列，要么一定要调用h.Reset()方法，不然生成的加密字符串会是拼接第一个字符串之后进行加密
// h.Reset()//重要！！！
// h.Write([]byte("string2"))
// fmt.Printf("%x\n", h.Sum(nil))
// hex.DecodeString 字符串转十六进制 字符串与十六进制加密产生的结果不一样，数字货币用十六进制
/**
package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func main() {
	sha1New := md5.New()
	src := "123456";
	srcB,_ := hex.DecodeString(src)
	fmt.Printf("123456 r => %x\n", srcB)
	b := []byte("123456")
	sha1New.Write(b)
	nSum := sha1New.Sum(nil)
	fmt.Printf("123456 r => %x\n", nSum)
	sha1New.Reset()
	sha1New.Write(srcB)
	nSum = sha1New.Sum(nil)
	fmt.Printf("123 r => %x\n", nSum)
}

*/

func Sha1Bytes(str string, pre []byte) string {
	sha1New := sha1.New()
	b := []byte(str)
	sha1New.Write(b)
	b = sha1New.Sum(pre)
	bStr := hex.EncodeToString(b)
	return bStr
}

func Md5Bytes(str string, pre []byte) string {
	mdNew := md5.New()
	b := []byte(str)
	mdNew.Write(b)
	b = mdNew.Sum(pre)
	//转换成字符串
	bStr := hex.EncodeToString(b)
	return bStr
}

func Sha512Bytes(str string, pre []byte) string {
	sha512New := sha512.New()
	b := []byte(str)
	sha512New.Write(b)
	b = sha512New.Sum(pre)
	bStr := hex.EncodeToString(b)
	return bStr
}

func Sha256Bytes(str string, pre []byte) string {
	sha256New := sha256.New()
	b := []byte(str)
	sha256New.Write(b)
	b = sha256New.Sum(pre)
	bStr := hex.EncodeToString(b)
	return bStr
}

func Ripemd160Bytes(str string, pre []byte) string {
	sha256New := ripemd160.New()
	b := []byte(str)
	sha256New.Write(b)
	b = sha256New.Sum(pre)
	bStr := hex.EncodeToString(b)
	return bStr
}

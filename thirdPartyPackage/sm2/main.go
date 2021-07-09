package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	gmSm2 "github.com/ZZMarquis/gm/sm2"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/fs"
	"io/ioutil"
	"math/big"
)

var (
	pubKeyStr = "AqH8MHA4c7yoZH+X+e8lqZlibkL4Ti9jMary93p8f69f" // 压缩后的公钥 使用 Compress函数
	// AqH8MHA4c7yoZH+X+e8lqZlibkL4Ti9jMary
	privKeyStr = "oBpT5FgdQXhIRJgBqY6jWcFZ1Ptd35sSOrwieHLdIg8="
	msgStr     = []byte("123456")
	// 阿里使用的是  r与s 拼接成一个byte 再base64 所以我们要先解码然后再以三十二位的位置切割成两份，r与s
	signStr = "asy5gEjCBFm9wE2fBH52FZRGiGz96NwI+pzDVoYwzPOtcT4r1g4Lokoi/E+scWZFKx2xvVworQ501krdMfXtxA=="
	// asy5gEjCBFm9wE2fBH52FZRGiGz96NwI+pzDVoYwzPOtcT4r1g4Lokoi/E+scWZFKx2xvVworQ501krdMfXtxA==
	// 21qmC0cXSP2MqKYRHVppyb5F9aX8G9MGNx86c7p5YVK8BUcnuFG2N2zLNynd6hocINoeDSWX8YtMNFLrPixgkA
	// MEUCIQC9R5QRlcVT+qtEEdU832Enh/KnvEByTyTvyz9ixUDwIQIgCFNNvTQA8jPstE/bOib6j1CJ37T5KBITjRjphCyG3K4=
	//ofwwcDhzvKhkf5f57yWpmWJuQvhOL2MxqvL3enx/r18
	//JX1mA421012QTzENrOtFAUmT3HXs8rXWj8j7XUeMTR4=

	//kZ7l1IZcprHfonuD7rxi4pBRl6Gk78zpLDI8MpGEIHc=
	//u9ClNyPYWtp8GQmBRQaB35QHlXwQ0+JBFdHEOINlobQ=

	// kZ7l1IZcprHfonuD7rxi4pBRl6Gk78zpLDI8MpGEIHc=
	// u9ClNyPYWtp8GQmBRQaB35QHlXwQ0+JBFdHEOINlobQ=
	//privKeyStr = "9WlS9KsxCZMbCeztUKgELgO3d4V/jqJGZvV94gDBKFc="

)

// 	bjt := []byte{213,212,199,197,205,250}
//	aacc, errjj := GbkToUtf8(bjt)
//	fmt.Println(aacc,errjj,string(aacc))
// GbkToUtf8 transform GBK bytes to UTF-8 bytes
func GbkToUtf8(bt []byte) ([]byte, error) {
	r := transform.NewReader(bytes.NewReader(bt), simplifiedchinese.GBK.NewDecoder())
	return ioutil.ReadAll(r)
}

// Utf8ToGbk transform UTF-8 bytes to GBK bytes
func Utf8ToGbk(bt []byte) (b []byte, err error) {
	r := transform.NewReader(bytes.NewReader(bt), simplifiedchinese.GBK.NewEncoder())
	return ioutil.ReadAll(r)
}

func main() {
	id := []byte("1234567812345678")
	p, pub, err := gmSm2.GenerateKey(rand.Reader)
	fmt.Printf("p => %#v pub => %#v  err => %v \n", p, pub, err)

	bt := p.GetRawBytes()
	str := base64.StdEncoding.EncodeToString(bt)
	strRaw := base64.StdEncoding.EncodeToString(p.D.Bytes())
	fmt.Printf("p.D => %v   str => %v  strRaw => %v\n", p.D, str, strRaw)

	fmt.Printf("pub.Y => %v   pub.X => %v\n", pub.Y, pub.X)

	pubX := base64.StdEncoding.EncodeToString(pub.X.Bytes())
	pubY := base64.StdEncoding.EncodeToString(pub.Y.Bytes())
	fmt.Printf("pubX => %v pubY => %v \n", pubX, pubY)

	btc := pub.GetUnCompressBytes()
	str = base64.StdEncoding.EncodeToString(btc)

	btc = pub.GetRawBytes()
	strRaw = base64.StdEncoding.EncodeToString(btc)
	fmt.Printf("pub GetUnCompressBytes => %v  GetRawBytes => %v\n", str, strRaw)

	fmt.Println("====================================================================")

	bt, err = base64.StdEncoding.DecodeString(privKeyStr)
	fmt.Println(bt, err)
	p, err = gmSm2.RawBytesToPrivateKey(bt)
	fmt.Println(bt, err)

	//  加签 start
	r, s, err := gmSm2.SignToRS(p, id, msgStr)

	str = base64.StdEncoding.EncodeToString(r.Bytes())

	strRaw = base64.StdEncoding.EncodeToString(s.Bytes())
	fmt.Printf("r. => %v s. => %v\n", str, strRaw)

	bt, err = gmSm2.Sign(p, id, msgStr)
	str = base64.StdEncoding.EncodeToString(bt)
	fmt.Println(str, err)

	// 加签 end

	// r. => eBuXV1F8oYwuYnuEQHotPgus4ZzWCNYLvfNbRNdNU+k= s. => BFsIyzizeAwYyHYxA7CLp3qe+IPoWkugAnQQ7vBgC/4=
	// MEYCIQDYNpbwCh3zI00SnOo99TYU0IasGxccO+oDwjHs8cdYPQIhAJqXNTZQhqKv/vP7oRdHmisxIRycuXHmVykkrwGTVE47

	// signStr = "MEQCIEoHKKnfevqzyOfzjo4CmccCjvDTC3+c2vN6OugjyLh7AiBW3kZkfBZlYDXDBwm0DwYsOFDj54nOIUyzHB+pB78WoQ=="

	// 支付宝加密后的字段需要base64解压，然后再通过三十二位一份切割，成两份r,s。
	// 而不是直接将加签得到的r,s直接 asn1.Marshal 编码
	bt, err = base64.StdEncoding.DecodeString(signStr)
	fmt.Println(bt, err)

	r = new(big.Int).SetBytes(bt[:32])
	s = new(big.Int).SetBytes(bt[32:])

	rBt := append(r.Bytes(), s.Bytes()...)
	rBtStr := base64.StdEncoding.EncodeToString(rBt)
	fmt.Println(rBtStr)

	//r, s, err = gmSm2.UnmarshalSign(bt)
	//fmt.Println(r, s, err)
	//rStr := base64.StdEncoding.EncodeToString(r.Bytes())
	//sStr := base64.StdEncoding.EncodeToString(s.Bytes())

	//fmt.Printf("rStr => %v sStr => %v rBtStr %v",rStr,sStr,rBtStr)
	// rStr => Sgcoqd96+rPI5/OOjgKZxwKO8NMLf5za83o66CPIuHs= sStr => Vt5GZHwWZWA1wwcJtA8GLDhQ4+eJziFMsxwfqQe/FqE=
	pub = gmSm2.CalculatePubKey(p)

	b := gmSm2.VerifyByRS(pub, id, msgStr, r, s)
	fmt.Println(b)

	bt = Compress(pub)
	str = base64.StdEncoding.EncodeToString(bt)

	pubX = base64.StdEncoding.EncodeToString(pub.X.Bytes())
	pubY = base64.StdEncoding.EncodeToString(pub.Y.Bytes())
	fmt.Printf("Compress pub => %v pubX => %v pubY => %v \n", str, pubX, pubY)

	bt, err = base64.StdEncoding.DecodeString(pubKeyStr)
	fmt.Println(bt, err)

	fmt.Println("***************************************************************")

	cc, err := ioutil.ReadFile("./pri1.pem")
	pk, err := x509.ReadPrivateKeyFromPem(cc, nil)
	fmt.Println(cc, err, pk)

	aa, err := sm2.GenerateKey(rand.Reader)
	fmt.Println(err)

	a, err := x509.MarshalSm2UnecryptedPrivateKey(aa)
	fmt.Println(a, err, base64.StdEncoding.EncodeToString(a))

	strB := base64.StdEncoding.EncodeToString(aa.PublicKey.X.Bytes())
	fmt.Println(strB)

	strB = base64.StdEncoding.EncodeToString(aa.PublicKey.Y.Bytes())
	fmt.Println(strB)

	strA := base64.StdEncoding.EncodeToString(aa.D.Bytes())
	fmt.Println(strA)

	bt = sm2.Compress(&aa.PublicKey)
	str = base64.StdEncoding.EncodeToString(bt)
	fmt.Println(str)

	bt, err = x509.WritePrivateKeyToPem(aa, nil)
	fmt.Println(err)

	err = ioutil.WriteFile("pri1.pem", bt, 0666)
	fmt.Println(err)

	//privKeyB, err := base64.StdEncoding.DecodeString(privKeyStr)
	//if err != nil {
	//	fmt.Println(err)
	//}

	//pri := &sm2.PrivateKey{
	//	PublicKey: sm2.PublicKey{
	//		Curve: sm2.P256Sm2(),
	//	},
	//	D: new(big.Int).SetBytes(privKeyB),
	//}
	//
	//pric, err := pkcs12.MarshalECPrivateKey(pri)
	//if err != nil {
	//	fmt.Println(err, pric)
	//}

	r, s, err = sm2.Sm2Sign(nil, msgStr, id, rand.Reader)
	if err != nil {
		fmt.Println(err)
	}
	data, err := sm2.SignDigitToSignData(r, s)
	if err != nil {
		fmt.Println(err)
	}
	str = base64.StdEncoding.EncodeToString(data)
	fmt.Println(str)
}

func Compress(a *gmSm2.PublicKey) []byte {
	buf := []byte{}
	yp := getLastBit(a.Y)
	buf = append(buf, a.X.Bytes()...)
	if n := len(a.X.Bytes()); n < 32 {
		buf = append(zeroByteSlice()[:(32-n)], buf...)
	}
	buf = append([]byte{byte(yp + 2)}, buf...)
	return buf
}

func getLastBit(a *big.Int) uint {
	return a.Bit(0)
}

// 32byte
func zeroByteSlice() []byte {
	return []byte{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}
}

func GeneratePrivateKey(privKeyStr string) (err error) {
	var (
		bt         []byte
		privateKey *sm2.PrivateKey
	)
	if bt, err = base64.StdEncoding.DecodeString(privKeyStr); err != nil {
		return err
	}
	curve := sm2.P256Sm2()
	privateKey = new(sm2.PrivateKey)
	privateKey.PublicKey.Curve = curve
	privateKey.D = new(big.Int).SetBytes(bt)
	privateKey.PublicKey.X, privateKey.PublicKey.Y = curve.ScalarBaseMult(bt)
	fmt.Println(privateKey.PublicKey.X, privateKey.PublicKey.Y)
	if bt, err = x509.WritePrivateKeyToPem(privateKey, nil); err != nil {
		return err
	}
	// 私钥写入文件
	if err = ioutil.WriteFile("./privateKey.pem", bt, fs.FileMode(0644)); err != nil {
		return err
	}
	buf := []byte{}
	yp := getLastBit(privateKey.PublicKey.Y)
	buf = append(buf, privateKey.PublicKey.X.Bytes()...)
	if n := len(privateKey.PublicKey.X.Bytes()); n < 32 {
		buf = append(zeroByteSlice()[:(32-n)], buf...)
	}
	buf = append([]byte{byte(yp + 2)}, buf...)
	str := base64.StdEncoding.EncodeToString(buf)
	fmt.Println(str)
	// 私钥写入文件
	if bt, err = x509.WritePublicKeyToPem(&privateKey.PublicKey); err != nil {
		return err
	}
	// 私钥写入文件
	if err = ioutil.WriteFile("./publicKey.pem", bt, fs.FileMode(0644)); err != nil {
		return err
	}
	return err
}

// Sm2Sign 私钥签名
func Sm2Sign(msg, uid []byte) (r, s *big.Int, err error) {
	var (
		bt  []byte
		pri *sm2.PrivateKey
	)
	if bt, err = ioutil.ReadFile("./privateKey.pem"); err != nil {
		return r, s, err
	}
	if pri, err = x509.ReadPrivateKeyFromPem(bt, nil); err != nil {
		return r, s, err
	}
	return sm2.Sm2Sign(pri, msg, uid, rand.Reader)
}

// Sm2Verify 公钥验签
func Sm2Verify(msg, uid []byte, r, s *big.Int) (b bool, err error) {
	var (
		bt  []byte
		pub *sm2.PublicKey
	)
	if bt, err = ioutil.ReadFile("./publicKey.pem"); err != nil {
		return b, err
	}
	if pub, err = x509.ReadPublicKeyFromPem(bt); err != nil {
		return b, err
	}
	b = sm2.Sm2Verify(pub, msg, uid, r, s)
	return b, err
}

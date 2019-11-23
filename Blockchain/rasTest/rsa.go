package rasTest

import "math"

/**
必备数学知识
		RSA加密算法中，只用到素数、互质数、指数运算、模运算等几个简单的数学知识。所以，我们也需要了解这几个概念即可
	1. 素数
		素数又称质数，指在一个大于1的自然数中，除了1和此整数自身外，不能被其他自然数整除的数。这个概念，我们在上初中，甚至小学的时候都学过了，这里就不再过多解释了。
	2. 互质数
		百度百科上的解释是：公因数只有1的两个数，叫做互质数。；维基百科上的解释是：互质，又称互素。若N个整数的最大公因子是1，则称这N个整数互质。

		常见的互质数判断方法主要有以下几种：
			1. 两个不同的质数一定是互质数。例如，2与7、13与19。
			2. 一个质数，另一个不为它的倍数，这两个数为互质数。例如，3与10、5与 26。
			3. 相邻的两个自然数是互质数。如 15与 16。
			4. 相邻的两个奇数是互质数。如 49与 51。
			5. 较大数是质数的两个数是互质数。如97与88。
			6. 小数是质数，大数不是小数的倍数的两个数是互质数。例如 7和 16。
			7. 2和任何奇数是互质数。例如2和87。
			8.  1不是质数也不是合数，它和任何一个自然数在一起都是互质数。如1和9908。
			9. 辗转相除法。
	3. 指数运算
		指数运算又称乘方计算，计算结果称为幂。nm指将n自乘m次。把nm看作乘方的结果，叫做”n的m次幂”或”n的m次方”。其中，n称为“底数”，m称为“指数”。
	4. 模运算
	　　模运算即求余运算。“模”是“Mod”的音译。和模运算紧密相关的一个概念是“同余”。数学上，当两个整数除以同一个正整数，
		若得相同余数，则二整数同余。

	　　两个整数a，b，若它们除以正整数m所得的余数相等，则称a，b对于模m同余，记作: a ≡ b (mod m)；读作：a同余于b模m，
		或者，a与b关于模m同余。例如：26 ≡ 14 (mod 12)。
RSA加密算法
	1. RSA加密算法简史
		RSA是1977年由罗纳德·李维斯特（Ron Rivest）、阿迪·萨莫尔（Adi Shamir）和伦纳德·阿德曼（Leonard Adleman）一起提出的。
		当时他们三人都在麻省理工学院工作。RSA就是他们三人姓氏开头字母拼在一起组成的。
	2. 公钥与密钥的产生
		假设Alice想要通过一个不可靠的媒体接收Bob的一条私人讯息。她可以用以下的方式来产生一个公钥和一个私钥：

		1. 随意选择两个大的质数p和q，p不等于q，计算N=pq。
		2. 根据欧拉函数，求得r = (p-1)(q-1)
		3. 选择一个小于 r 的整数 e，求得 e 关于模 r 的模反元素 (e*d)%r=1 也就是 e*d+1= r*y ，命名为 d。（模反元素存在，当且仅当e与r互质）
		4. 将 p 和 q 的记录销毁。
		5. (N,e)是公钥，(N,d)是私钥。Alice将她的公钥(N,e)传给Bob，而将她的私钥(N,d)藏起来。
	3. 加密消息
		假设Bob想给Alice送一个消息m，他知道Alice产生的N和e。他使用起先与Alice约好的格式将m转换为一个小于N的整数n，
		比如他可以将每一个字转换为这个字的Unicode码，然后将这些数字连在一起组成一个数字。假如他的信息非常长的话，
		他可以将这个信息分为几段，然后将每一段转换为n。用下面这个公式他可以将n加密为c：
			ne ≡ c (mod N)  [m^e%N= c] m的e次方
		计算c并不复杂。Bob算出c后就可以将它传递给Alice。
	4. 解密消息
		Alice得到Bob的消息c后就可以利用她的密钥d来解码。她可以用以下这个公式来将c转换为n：

		　　cd ≡ n (mod N)  [c^d%N = m] a的d次方

		得到n后，她可以将原来的信息m重新复原。

		解码的原理是：

		　　cd ≡ n e·d(mod N)

		以及ed ≡ 1 (mod p-1)和ed ≡ 1 (mod q-1)。由费马小定理可证明（因为p和q是质数）

		　　n e·d ≡ n (mod p) 　　和 　n e·d ≡ n (mod q)

		这说明（因为p和q是不同的质数，所以p和q互质）

		　　n e·d ≡ n (mod pq)
	5. 签名消息
		RSA也可以用来为一个消息署名。假如甲想给乙传递一个署名的消息的话，那么她可以为她的消息计算一个
		散列值(Message digest)，然后用她的密钥(private key)加密这个散列值并将这个“署名”加在消息的后面。这个消息只有用她的
		公钥才能被解密。乙获得这个消息后可以用甲的公钥解密这个散列值，然后将这个数据与他自己为这个消息计算的散列值相比较。
		假如两者相符的话，那么他就可以知道发信人持有甲的密钥，以及这个消息在传播路径上没有被篡改过。
	6. 编程实践
		1. 假设p = 3、q = 11（p，q都是素数即可。），则N = pq = 33；
		2. r = (p-1)(q-1) = (3-1)(11-1) = 20；
		3. 根据模反元素的计算公式，我们可以得出，(e为整数。且 1 < e < r ) e·d ≡ 1 (mod 20),即e·d = 20n+1 (n为正整数)；我们假设n=1，
		则e·d = 21。e、d为正整数，并且e与r互质，则e = 3，d = 7。（两个数交换一下也可以。）

　　	到这里，公钥和密钥已经确定。公钥为(N, e) = (33, 3)，密钥为(N, d) = (33, 7)。

		package main

		import (
			"fmt"
			"github.com/ZRothschild/goTest/Blockchain/rasTest"
		)

		func main() {
			N,d,e := rasTest.RsaTest()
			fmt.Printf("N => %d d => %d e => %d\n",N,d,e)
			var src int = 6
			c := rasTest.RsaEncryption(src,N,e)
			fmt.Println(c)
			src = rasTest.RsaDecrypt(c,N,d)
			fmt.Println(src)
		}
*/

func RsaTest() (N, d, e int) {
	var (
		p int = 3
		q int = 11
		r int = (p - 1) * (q - 1)
	)
	for i := 2; i < r; i++ {
		if (r+1)%i == 0 {
			e = i
			break
		}
	}
	N = p * q
	d = (r + 1) / e
	return
}

//加密 (N,e) 公钥 加密的值小于 N src <N
func RsaEncryption(src, N, e int) (c int) {
	pow := math.Pow(float64(src), float64(e))
	c = int(pow) % N
	return
}

//解密 (N,d) 密钥
func RsaDecrypt(c, N, d int) (src int) {
	pow := math.Pow(float64(c), float64(d))
	src = int(pow) % N
	return
}

/**

package main

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)
//rsa加密解密 签名验签
func main() {
	//生成私钥
	priv, e := rsa.GenerateKey(rand.Reader, 1024)
	if e != nil {
		fmt.Println(e)
	}

	//根据私钥产生公钥
	pub := &priv.PublicKey

	//明文
	plaintext := []byte("Hello world")

	//加密生成密文
	fmt.Printf("%q\n加密:\n", plaintext)
	ciphertext, e := rsa.EncryptOAEP(md5.New(), rand.Reader, pub, plaintext, nil)
	if e != nil {
		fmt.Println(e)
	}
	fmt.Printf("\t%x\n", ciphertext)

	//解密得到明文
	fmt.Printf("解密:\n")
	plaintext, e = rsa.DecryptOAEP(md5.New(), rand.Reader, priv, ciphertext, nil)
	if e != nil {
		fmt.Println(e)
	}
	fmt.Printf("\t%q\n", plaintext)

	//消息先进行Hash处理
	h := md5.New()
	h.Write(plaintext)
	hashed := h.Sum(nil)
	fmt.Printf("%q MD5 Hashed:\n\t%x\n", plaintext, hashed)

	//签名
	opts := &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: crypto.MD5}
	sig, e := rsa.SignPSS(rand.Reader, priv, crypto.MD5, hashed, opts)
	if e != nil {
		fmt.Println(e)
	}
	fmt.Printf("签名:\n\t%x\n", sig)

	//认证
	fmt.Printf("验证结果:")
	if e := rsa.VerifyPSS(pub, crypto.MD5, hashed, sig, opts); e != nil {
		fmt.Println("失败:", e)
	} else {
		fmt.Println("成功.")
	}
}
*/

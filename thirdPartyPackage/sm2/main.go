package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/guanzhi/GmSSL/go/gmssl"
	"github.com/tjfoc/gmsm/pkcs12"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
	"io/ioutil"
	//"github.com/tjfoc/gmsm/x509"
)

var (
	pubKeyStr  = "AqH8MHA4c7yoZH+X+e8lqZlibkL4Ti9jMary93p8f69f"
	privKeyStr = "oBpT5FgdQXhIRJgBqY6jWcFZ1Ptd35sSOrwieHLdIg8="
	msgStr     = "helloworld"
	// base64.StdEncoding.EncodeToString(sign)
	signStr = "DIk02bFoZJLbR5EgLX5ZULRF78aP0EYYAKbDjmV8PeNwG4EEkcbRdiw4BpIIyKGpID4lFv0u26+KfgDeEfN/MQ=="
)

func main() {
	/* Engines */
	fmt.Print("Engines:")
	engines := gmssl.GetEngineNames()
	for _, engine := range engines {
		fmt.Print(" " + engine)
	}
	fmt.Println("\n");

	/* private key */
	rsa_args := [][2]string{
		{"rsa_keygen_bits", "2048"},
		{"rsa_keygen_pubexp", "65537"},
	}

	rsa, err := gmssl.GeneratePrivateKey("RSA", rsa_args, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	rsa_pem, err := rsa.GetPublicKeyPEM()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rsa_pem)

	/* Engine */
	eng, _ := gmssl.NewEngineByName(engines[1])
	cmds, _ := eng.GetCommands()
	for _, cmd := range cmds {
		fmt.Print(" " + cmd)
	}
	fmt.Println()

	/* SM2 key pair operations */
	sm2keygenargs := [][2]string{
		{"ec_paramgen_curve", "sm2p256v1"},
		{"ec_param_enc", "named_curve"},
	}
	sm2sk, _ := gmssl.GeneratePrivateKey("EC", sm2keygenargs, nil)
	sm2sktxt, _ := sm2sk.GetText()
	sm2skpem, _ := sm2sk.GetPEM("SMS4", "password")
	sm2pkpem, _ := sm2sk.GetPublicKeyPEM()

	fmt.Println(sm2sktxt)
	fmt.Println(sm2skpem)
	fmt.Println(sm2pkpem)

	sm2pk, _ := gmssl.NewPublicKeyFromPEM(sm2pkpem)
	sm2pktxt, _ := sm2pk.GetText()
	sm2pkpem_, _ := sm2pk.GetPEM()

	fmt.Println(sm2pktxt)
	fmt.Println(sm2pkpem_)

	/* SM2 sign/verification */
	sm3ctx, _ := gmssl.NewDigestContext("SM3")
	sm2zid, _ := sm2pk.ComputeSM2IDDigest("1234567812345678")
	sm3ctx.Reset()
	sm3ctx.Update(sm2zid)
	sm3ctx.Update([]byte("message"))
	tbs, _ := sm3ctx.Final()

	sig, _ := sm2sk.Sign("sm2sign", tbs, nil)
	fmt.Printf("sm2sign(sm3(\"message\")) = %x\n", sig)

	if ret := sm2pk.Verify("sm2sign", tbs, sig, nil); ret != nil {
		fmt.Printf("sm2 verify failure\n")
	} else {
		fmt.Printf("sm2 verify success\n")
	}

	/* SM2 encrypt */
	sm2msg := "01234567891123456789212345678931234567894123456789512345678961234567897123"
	sm2encalg := "sm2encrypt-with-sm3"
	sm2ciphertext, _ := sm2pk.Encrypt(sm2encalg, []byte(sm2msg), nil)
	sm2plaintext, _ := sm2sk.Decrypt(sm2encalg, sm2ciphertext, nil)
	fmt.Printf("sm2enc(\"%s\") = %x\n", sm2plaintext, sm2ciphertext)
	if sm2msg != string(sm2plaintext) {
		fmt.Println("SM2 encryption/decryption failure")
	}


	cc ,err := ioutil.ReadFile("./pri.pem")
	fmt.Println(err,cc)

	sm2sk, err = gmssl.NewPrivateKeyFromPEM(string(cc), "")

	fmt.Println(sm2sk,err)

	p , _ := pem.Decode(cc)
	fmt.Println(len(p.Bytes))

	pk ,err := x509.ParseSm2PrivateKey(p.Bytes)
	fmt.Println(err,pk)

	id := []byte("1234567812345678")
	privKeyB, err := base64.StdEncoding.DecodeString(privKeyStr)
	if err != nil {
		fmt.Println(err)
	}

	aa, err := sm2.GenerateKey(rand.Reader)
	fmt.Println(err)

	bt, err := x509.WritePrivateKeyToPem(aa, nil)
	fmt.Println(err)

	err = ioutil.WriteFile("pri1.pem", bt, 0666)
	fmt.Println(err)

	pric, err := pkcs12.ParsePKCS8PrivateKey(privKeyB)
	if err != nil {
		fmt.Println(err, pric)
	}
	r, s, err := sm2.Sm2Sign(nil, []byte(msgStr), id, rand.Reader)
	if err != nil {
		fmt.Println(err)
	}
	data, err := sm2.SignDigitToSignData(r, s)
	if err != nil {
		fmt.Println(err)
	}
	str := base64.StdEncoding.EncodeToString(data)
	fmt.Println(str)
}

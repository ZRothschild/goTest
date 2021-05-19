// Package tlv contains code to read and write type-length-value messages.
// https://gw.alipayobjects.com/os/skylark-tools/public/files/4df27ca9443dcf9131c4a1b69376fb53.pdf
package main

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/tjfoc/gmsm/x509"
)

// MaxMessageSize defines how large a message can be before we reject it.
const MaxMessageSize = 1024 * 1024 * 1024 // 1GB

// ReadTLV reads a type-length-value record from r.
func ReadTLV(r io.Reader) (byte, []byte, error) {
	typ, err := ReadType(r)
	if err != nil {
		return 0, nil, err
	}

	buf, err := ReadLV(r)
	if err != nil {
		return 0, nil, err
	}
	return typ, buf, err
}

// ReadType reads the type from a TLV record.
func ReadType(r io.Reader) (byte, error) {
	var typ [1]byte
	if _, err := io.ReadFull(r, typ[:]); err != nil {
		if err == io.EOF {
			return 0, err
		} else {
			return 0, fmt.Errorf("read message type: %s", err)
		}
	}
	return typ[0], nil
}

// ReadLV reads the length-value from a TLV record.
func ReadLV(r io.Reader) ([]byte, error) {
	// Read the size of the message.
	var sz int64
	if err := binary.Read(r, binary.BigEndian, &sz); err != nil {
		return nil, fmt.Errorf("read message size: %s", err)
	}

	if sz < 0 {
		return nil, fmt.Errorf("negative message size is invalid: %d", sz)
	}

	if sz >= MaxMessageSize {
		return nil, fmt.Errorf("max message size of %d exceeded: %d", MaxMessageSize, sz)
	}

	// Read the value.
	buf := make([]byte, sz)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("read message value: %s", err)
	}

	return buf, nil
}

type Writer struct {
	bytes.Buffer
}

// WriteBigEndianInt32 writes the BigEndian int32 to w.
func (w *Writer) WriteBigEndianInt32(value int32) error {
	if err := binary.Write(w, binary.BigEndian, value); err != nil {
		return fmt.Errorf("write BigEndian int32: %s", err)
	}
	return nil
}

// WriteBigEndianInt64 writes the BigEndian int32 to w.
func (w *Writer) WriteBigEndianInt64(value int64) error {
	if err := binary.Write(w, binary.BigEndian, value); err != nil {
		return fmt.Errorf("write BigEndian int64: %s", err)
	}
	return nil
}

// WriteLV writes the Length(byte)+value to w.
func (w *Writer) WriteLV(buf []byte) error {
	// Write the size of the message.
	length := len(buf)
	if length > 255 {
		return errors.New("buf len greater than 255")
	}
	if err := w.WriteByte(byte(length)); err != nil {
		return err
	}

	// Write the value.
	if _, err := w.Write(buf); err != nil {
		return fmt.Errorf("write message value: %s", err)
	}
	return nil
}

// WriteBLV writes the Length(int32)+value to w.
func (w *Writer) WriteBLV(buf []byte) error {
	// Write the size of the message.
	if err := binary.Write(w, binary.BigEndian, int32(len(buf))); err != nil {
		return fmt.Errorf("write message size: %s", err)
	}

	// Write the value.
	if _, err := w.Write(buf); err != nil {
		return fmt.Errorf("write message value: %s", err)
	}
	return nil
}

// WriteTLV writes a Tag(byte)+Length(byte)+Value to w.
func (w *Writer) WriteTLV(tag, length byte, buf []byte) error {
	if err := w.WriteByte(tag); err != nil {
		return err
	}
	if err := w.WriteByte(length); err != nil {
		return err
	}
	// Write the value.
	if _, err := w.Write(buf); err != nil {
		return fmt.Errorf("write message value: %s", err)
	}
	return nil
}

func (w *Writer) WriteTLVByString(tag byte, value string) error {
	b := []byte(value)
	length := len(b)
	if length > 255 {
		return errors.New("value len greater than 255")
	}
	// Write the value.
	if err := w.WriteTLV(tag, byte(length), b); err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteTLVByByte(tag byte, value byte) error {
	// Write the value.
	if err := w.WriteTLV(tag, 1, []byte{value}); err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteTLVByUint64(tag byte, value uint64) error {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, value)
	// Write the value.
	if err := w.WriteTLV(tag, 8, b); err != nil {
		return err
	}
	return nil
}

type Order struct {
	OutOrderId    string `json:"out_order_id"`   // 离线订单号
	PID           uint64 `json:"pid"`            // PID, ⼤端序
	CompanyCode   string `json:"company_code"`   // 航空公司2位 代码
	FlightNum     string `json:"flight_num"`     // 航班编号
	SeatNum       string `json:"seat_num"`       // 座位号
	GoodsName     string `json:"goods_name"`     // 商品名称
	GoodsAmount   uint64 `json:"goods_amount"`   // ⾦额, ⼤端序
	IdentityType  byte   `json:"identity_type"`  // 旅客身份信息类型
	IdentityId    string `json:"identity_id"`    // 旅客身份ID 脱敏信息
	DeviceId      string `json:"device_id"`      // ⽣成订单码的设备编号
	GoodsCategory byte   `json:"goods_category"` // 商品类⽬
}

// WriteUserinfo 写入订单主体信息
func (w *Writer) WriteUserinfo(order Order) error {
	// Write data.
	var dataW Writer
	if err := dataW.WriteTLVByString(1, order.OutOrderId); err != nil {
		return err
	}
	if err := dataW.WriteTLVByUint64(2, order.PID); err != nil {
		return err
	}
	if err := dataW.WriteTLVByString(3, order.CompanyCode); err != nil {
		return err
	}
	if err := dataW.WriteTLVByString(4, order.FlightNum); err != nil {
		return err
	}
	if err := dataW.WriteTLVByString(5, order.SeatNum); err != nil {
		return err
	}
	if err := dataW.WriteTLVByString(6, order.GoodsName); err != nil {
		return err
	}
	if err := dataW.WriteTLVByUint64(7, order.GoodsAmount); err != nil {
		return err
	}
	if err := dataW.WriteTLVByByte(9, order.IdentityType); err != nil {
		return err
	}
	if err := dataW.WriteTLVByString(10, order.IdentityId); err != nil {
		return err
	}
	if err := dataW.WriteTLVByString(11, order.DeviceId); err != nil {
		return err
	}
	if err := dataW.WriteTLVByByte(12, order.GoodsCategory); err != nil {
		return err
	}

	// Write userinfo.
	// 加密算法
	if err := w.WriteByte(0); err != nil {
		return err
	}
	// 密钥ID
	if err := w.WriteByte(0); err != nil {
		return err
	}
	// 数据模版ID ⼤端序的int类型byte[]
	if err := w.WriteBigEndianInt32(1000); err != nil {
		return err
	}
	// 数据(BLV结构)
	if err := w.WriteBLV(dataW.Bytes()); err != nil {
		return fmt.Errorf("write message value: %s", err)
	}
	return nil
}

func (w *Writer) WriteOrder(order Order, pubKey, priKey string) ([]byte, error) {
	// 写死28.表示让支付宝受理
	if err := w.WriteBigEndianInt32(23); err != nil {
		return nil, err
	}
	// 数据格式版本6
	if err := w.WriteByte(6); err != nil {
		return nil, err
	}
	// 编码格式gb2312
	if err := w.WriteByte(1); err != nil {
		return nil, err
	}
	// 根密钥ID,可以在支付宝配置多组公钥，这个值指定公钥的key. 从1开始。0表示 不验证签名
	if err := w.WriteByte(2); err != nil {
		return nil, err
	}
	// 算法 - SM2
	if err := w.WriteByte(1); err != nil {
		return nil, err
	}
	// 航司设备的公钥
	if err := w.WriteLV([]byte(pubKey)); err != nil {
		return nil, err
	}

	// 航司根私钥签名
	privPem, err := ioutil.ReadFile("./certs/sm2PriKey.pem")
	if err != nil {
		fmt.Println("ReadFile err: ", err)
		return nil, err
	}
	privKey, err := x509.ReadPrivateKeyFromPem(privPem, nil) // 读取密钥
	if err != nil {
		fmt.Println("ReadPrivateKeyFromPem err: ", err)
		return nil, err
	}
	sign, err := privKey.Sign(nil, w.Bytes(), nil)
	if err != nil {
		fmt.Println("Sign err: ", err)
		return nil, err
	}
	if err := w.WriteLV(sign); err != nil {
		return nil, err
	}
	fmt.Println("order: ", order)
	fmt.Println("1 w len: ", w.Len())
	// 写入订单主体信息
	if err := w.WriteUserinfo(order); err != nil {
		return nil, err
	}
	fmt.Println("2 w len: ", w.Len())
	// 过期时间
	if err := w.WriteBigEndianInt64(time.Now().Unix() + 7200); err != nil {
		return nil, err
	}

	// 航司设备私钥签名
	sign2, err := privKey.Sign(nil, w.Bytes(), nil)
	if err != nil {
		return nil, err
	}
	if err := w.WriteLV(sign2); err != nil {
		return nil, err
	}
	fmt.Println("3 w len: ", w.Len())
	return w.Bytes(), nil
}

// WriteOrder
func (w *Writer) GetEncodeToString() (string, error) {
	return base64.StdEncoding.EncodeToString(w.Bytes()), nil
}

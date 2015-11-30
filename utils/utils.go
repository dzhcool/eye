package utils

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"text/template"
	// "time"
	"github.com/dzhcool/eye"
	"github.com/dzhcool/eye/qhydra"
)

func Trace(msg string) {
	log.Println(msg)
}

func IsLan(ipAddr string) bool {
	Addr := strings.Split(ipAddr, ":")
	rAddr := net.ParseIP(Addr[0])
	if rAddr == nil {
		return false
	}
	ip := Inet_aton(rAddr)
	return ((ip&0xff000000) == 0x0a000000 || (ip&0xfff00000) == 0xac100000 || (ip&0xffff0000) == 0xc0a80000)
}

func Inet_ntoa(ipnr int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)
	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

func Inet_aton(ipnr net.IP) int64 {
	bits := strings.Split(ipnr.String(), ".")
	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}

//过滤函数
func Input_val(val string) string {
	return template.HTMLEscapeString(val)
}
func Input(val string) string {
	return Input_val(val)
}

func Urlencode(val string) string {
	//todo
	/*val = strings.Replace(val, "+", "%2B", -1)
	val = strings.Replace(val, "/", "%2F", -1)
	val = strings.Replace(val, "=", "%3D", -1)*/
	val = url.QueryEscape(val)
	return val
}

func Urldecode(val string) string {
	//todo
	// val = strings.Replace(val, "%2B", "+", -1)
	// val = strings.Replace(val, "%2F", "/", -1)
	// val = strings.Replace(val, "%3D", "=", -1)
	val, _ = url.QueryUnescape(val)
	return val
}

func Md5(val string) string {
	if val == "" {
		return ""
	}
	h := md5.New()
	io.WriteString(h, val)
	secret := fmt.Sprintf("%x", h.Sum(nil))
	return secret
}

//处理轮数 *****  和项目相关
func Round(all int, size int) int {
	if size <= 0 {
		size = 1
	}
	mod := all % size
	if mod <= 0 {
		mod = size
	}
	round := 100 + mod
	return round
}

//随机数
func Rand(max int) int {
	// r := rand.New(rand.NewSource(time.Now().UnixNano()))
	c := rand.Intn(max)
	return c
}

//数据类型转换
var ErrNil = errors.New("nil or type error")

func Int(reply interface{}, err error) (int, error) {
	if err != nil {
		return 0, err
	}
	switch reply := reply.(type) {
	case int64:
		x := int(reply)
		if int64(x) != reply {
			return 0, strconv.ErrRange
		}
		return x, nil
	case []byte:
		n, err := strconv.ParseInt(string(reply), 10, 0)
		return int(n), err
	case string:
		n, err := strconv.Atoi(reply)
		return n, err
	case nil:
		return 0, ErrNil
	}
	return 0, ErrNil
}

func Int64(reply interface{}, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	switch reply := reply.(type) {
	case int64:
		return reply, nil
	case []byte:
		n, err := strconv.ParseInt(string(reply), 10, 64)
		return n, err
	case nil:
		return 0, ErrNil
	}
	return 0, ErrNil
}

func Float64(reply interface{}, err error) (float64, error) {
	if err != nil {
		return 0, err
	}
	switch reply := reply.(type) {
	case []byte:
		n, err := strconv.ParseFloat(string(reply), 64)
		return n, err
	case nil:
		return 0, ErrNil
	}
	return 0, ErrNil
}

func String(reply interface{}, err error) (string, error) {
	if err != nil {
		return "", err
	}
	switch reply := reply.(type) {
	case []byte:
		return string(reply), nil
	case string:
		return reply, nil
	case int:
		n := strconv.Itoa(reply)
		return n, nil
	case nil:
		return "", ErrNil
	}
	return "", ErrNil
}

//获取http请求
func GetUrl(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", err
	}
	return String(result, nil)
}

func SendQueueMsg(data map[string]string) error {
	if msg, num, err := qhydra.QMsg.Trigger(data, ""); err != nil {
		logkit.Logger.Error("qhydra.Trigger:"+fmt.Sprintf("msg:%s,num:%d", msg, num), "utils/sendQueueMsg")
		return err
	} else {
		return nil
	}
}

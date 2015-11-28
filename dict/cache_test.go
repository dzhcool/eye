/**
 * 广告数据存储定义
 * @date 2015-08-05
 */
package dict

import (
	// "reflect"
	"github.com/pquerna/ffjson/ffjson"
	"log"
	"testing"
	"time"
)

type Ad struct {
	Ad_id   int    `json:"ad_id"`
	Pos     string `json:"pos"`
	Type    string `json:"type"`
	Jumpurl string `json:"jumpurl"`
	Url     string `json:"url"`
}

type AdRest struct {
	Errno  int    `json:"errno"`
	Errmsg string `json:"errmsg"`
	Data   []Ad   `json:"data"`
}

func TestCache(t *testing.T) {
	var (
		k = "testkey"
		v = "testval"
	)
	cacheTest := Cache("test")
	//存储
	cacheTest.Set(k, v, 0*time.Second)

	//查询缓存值
	val, err := cacheTest.String(cacheTest.Get(k))

	if len(val) > 0 {
		t.Log("Cache Ok")
	} else {
		t.Fatal("Cache Error", err)
	}
}

func TestJson(t *testing.T) {

	js := `{"ad_id":1, "pos":"31057001", "ext":"xxx"}`

	var buf Ad
	err := ffjson.Unmarshal([]byte(js), &buf)
	if err != nil {
		t.Fatal("decode json error", err)
	}

	var rest AdRest
	d := []Ad{buf}
	rest.Data = d

	json, err := ffjson.Marshal(rest)
	if err != nil {
		t.Fatal("encode json error", err)
	}
	log.Println(string(json))
}

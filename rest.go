/**
 * rest返回值结构类型定义
 */
package eye

import (
	"github.com/pquerna/ffjson/ffjson"
)

type RestCtx struct {
	Errno  int    `json:"error"`
	Errmsg string `json:"errmsg"`
	Count  int    `json:"count"`

	callbackName string
}

type IRestCtx interface {
	Prepare()
	JsonByte()
	JsonString()
}

func (p *RestCtx) Prepare() {
	//扩展
}

//结构体转化json byte数据
func (p *RestCtx) JsonByte() ([]byte, error) {
	json, err := ffjson.Marshal(&p)
	return json, err
}

//结构体转化json string数据
func (p *RestCtx) JsonString() (string, error) {
	json, err := p.JsonByte()
	return string(json), err
}

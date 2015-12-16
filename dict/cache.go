/**
 * 数据存储定义
 * @author dzh
 * @date 2015-08-05
 */
package dict

import (
	"errors"
	"github.com/dzhcool/eye/utils"
	"sync"
	"time"
)

var (
	cache = make(map[string]*CacheTable)
	mutex sync.RWMutex
)
var ErrNil = errors.New("nil or type error")

type CacheTable struct {
	sync.RWMutex
	name    string
	items   map[interface{}]*CacheItem
	addTime time.Time
}

func Cache(table string) *CacheTable {
	mutex.RLock()
	t, ok := cache[table]
	mutex.RUnlock()

	if !ok {
		addtime := time.Now()
		t = &CacheTable{
			name:    table,
			items:   make(map[interface{}]*CacheItem),
			addTime: addtime,
		}
		mutex.Lock()
		cache[table] = t
		mutex.Unlock()
	}
	return t
}

func (p *CacheTable) Count(k interface{}, v interface{}) int {
	num := len(p.items)
	return num
}

func (p *CacheTable) Set(k interface{}, v interface{}, l ...int) *CacheItem {
	lifetime := 0 * time.Second
	if len(l) > 0 {
		lifetime = time.Duration(l[0]) * time.Second
	}
	item := createCacheItem(k, v, lifetime)

	p.Lock()
	p.items[k] = &item
	p.Unlock()
	return &item
}

func (p *CacheTable) Add(k interface{}, v interface{}, l ...int) *CacheItem {
	return p.Set(k, v, l...)
}

func (p *CacheTable) Get(k interface{}) (interface{}, error) {
	p.RLock()
	r, ok := p.items[k]
	if !ok {
		p.RUnlock()
		return nil, ErrNil
	}
	p.RUnlock()

	if r.Expired() {
		p.Lock()
		delete(p.items, k)
		p.Unlock()
		return nil, ErrNil
	}

	return r.data, nil
}

func (p *CacheTable) Item(k interface{}) *CacheItem {
	p.RLock()
	r, ok := p.items[k]
	p.RUnlock()

	r.Lock()
	defer r.Unlock()

	if !ok {
		r = nil
	}
	return r
}

func (p *CacheTable) Items() map[interface{}]*CacheItem {
	p.Lock()
	defer p.Unlock()

	return p.items
}

func (p *CacheTable) Exists(k interface{}) bool {
	p.RLock()
	defer p.RUnlock()

	_, ok := p.items[k]
	return ok
}

func (p *CacheTable) Delete(k interface{}) (*CacheItem, error) {
	p.RLock()
	r, ok := p.items[k]
	p.RUnlock()

	if ok {
		p.Lock()
		delete(p.items, k)
		p.Unlock()
		return r, nil
	}
	return nil, ErrNil
}

func (p *CacheTable) Int(reply interface{}, err error) (int, error) {
	return utils.Int(reply, err)
}

func (p *CacheTable) Int64(reply interface{}, err error) (int64, error) {
	return utils.Int64(reply, err)
}

func (p *CacheTable) Float64(reply interface{}, err error) (float64, error) {
	return utils.Float64(reply, err)
}

func (p *CacheTable) String(reply interface{}, err error) (string, error) {
	return utils.String(reply, err)
}

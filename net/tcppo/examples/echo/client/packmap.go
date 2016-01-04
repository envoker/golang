package main

import (
	"fmt"
	"sync"
)

type PackMap struct {
	m    sync.Mutex
	bmap map[int]*bool
}

func NewPackMap() *PackMap {
	return &PackMap{
		bmap: make(map[int]*bool),
	}
}

func (p *PackMap) AddSendId(id int) bool {

	p.m.Lock()
	defer p.m.Unlock()

	_, ok := p.bmap[id]
	if ok {
		return false
	}

	p.bmap[id] = new(bool)

	return true
}

func (p *PackMap) AddReceiveId(id int) bool {

	p.m.Lock()
	defer p.m.Unlock()

	b, ok := p.bmap[id]
	if !ok {
		return false
	}

	*b = true

	return true
}

func (p *PackMap) test() {

	p.m.Lock()
	defer p.m.Unlock()

	for key, b := range p.bmap {
		if !(*b) {
			fmt.Println("key:", key)
			return
		}
	}

	fmt.Println("ok")
}

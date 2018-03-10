package main

type Holder struct {
	clients map[string]*Client //记录
}

func newHolder() *Holder {
	holder := &Holder{
		clients: make(map[string]*Client),
	}
	return holder
}
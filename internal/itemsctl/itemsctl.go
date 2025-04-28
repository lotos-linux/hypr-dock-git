package itemsctl

import "hypr-dock/internal/item"

type List struct {
	list map[string]*item.Item
}

func New() *List {
	return &List{
		list: make(map[string]*item.Item),
	}
}

func (l *List) GetMap() map[string]*item.Item {
	return l.list
}

func (l *List) Get(className string) *item.Item {
	return l.list[className]
}

func (l *List) Add(className string, item *item.Item) {
	l.list[className] = item
}

func (l *List) Remove(className string) {
	delete(l.list, className)
}

func (l *List) Len() int {
	return len(l.list)
}

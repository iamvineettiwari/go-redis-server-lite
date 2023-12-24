package list

import (
	"sync"

	"github.com/iamvineettiwari/go-redis-server-lite/resp"
)

type ListNode struct {
	data resp.ArrayType
	prev *ListNode
	next *ListNode
}

type List struct {
	head         *ListNode
	tail         *ListNode
	totalElement int64
	lock         *sync.RWMutex
}

func NewList() *List {
	return &List{
		head:         nil,
		tail:         nil,
		totalElement: 0,
		lock:         &sync.RWMutex{},
	}
}

func (l *List) IsEmpty() bool {
	return l.totalElement == 0
}

func (l *List) InsertLast(data any, dataType string) {
	l.lock.Lock()
	defer l.lock.Unlock()

	newNode := ListNode{
		data: resp.ArrayType{
			Value: data,
			Type:  dataType,
		},
	}

	if l.IsEmpty() {
		l.head = &newNode
		l.tail = &newNode
	} else {
		newNode.prev = l.tail
		l.tail.next = &newNode
		l.tail = &newNode
	}

	l.totalElement++
}

func (l *List) InsertFirst(data any, dataType string) {
	l.lock.Lock()
	defer l.lock.Unlock()

	newNode := ListNode{
		data: resp.ArrayType{
			Value: data,
			Type:  dataType,
		},
	}

	if l.IsEmpty() {
		l.head = &newNode
		l.tail = &newNode
	} else {
		newNode.next = l.head
		l.head.prev = &newNode
		l.head = &newNode
	}

	l.totalElement++
}

func (l *List) GetValues() []resp.ArrayType {
	data := []resp.ArrayType{}

	if l.IsEmpty() {
		return data
	}

	l.lock.RLock()
	defer l.lock.RUnlock()

	headPtr := l.head

	for headPtr != nil {
		data = append(data, headPtr.data)
		headPtr = headPtr.next
	}

	return data
}

package lfu

import "sync"

type node struct {
	key   string
	value string
	freq  int
	prev  *node
	next  *node
}

type doublyLinkedList struct {
	head *node
	tail *node
}

func (dll *doublyLinkedList) addNode(n *node) {
	n.prev = dll.tail
	n.next = nil
	dll.tail.next = n
	dll.tail = n
}

func (dll *doublyLinkedList) removeNode(n *node) {
	prev := n.prev
	next := n.next
	prev.next = next
	if next != nil {
		next.prev = prev
	} else {
		dll.tail = prev
	}
}

type LFUCache struct {
	capacity int
	size     int
	minFreq  int
	freqMap  map[int]*doublyLinkedList
	keyMap   map[string]*node
	mu       sync.Mutex
}

func NewLFUCache(capacity int) *LFUCache {
	return &LFUCache{
		capacity: capacity,
		size:     0,
		minFreq:  0,
		freqMap:  make(map[int]*doublyLinkedList),
		keyMap:   make(map[string]*node),
	}
}

func (lfu *LFUCache) Get(key string) (string, int) {
	lfu.mu.Lock()
	defer lfu.mu.Unlock()

	if _, ok := lfu.keyMap[key]; !ok {
		return "", -1
	}

	node := lfu.keyMap[key]
	dll := lfu.freqMap[node.freq]
	dll.removeNode(node)

	if dll.head.next == dll.tail {
		delete(lfu.freqMap, node.freq)
		if lfu.minFreq == node.freq {
			lfu.minFreq++
		}
	}

	node.freq++
	lfu.freqMap[node.freq] = &doublyLinkedList{
		head: node,
		tail: node,
	}
	lfu.freqMap[node.freq].addNode(node)

	return node.value, node.freq
}

func (lfu *LFUCache) Put(key string, value string) {
	lfu.mu.Lock()
	defer lfu.mu.Unlock()

	if lfu.capacity == 0 {
		return
	}

	if node, ok := lfu.keyMap[key]; ok {
		node.value = value
		lfu.Get(key) // increase freq
		return
	}

	if lfu.size == lfu.capacity {
		lfuNode := lfu.freqMap[lfu.minFreq].head.next
		lfu.freqMap[lfu.minFreq].removeNode(lfuNode)
		delete(lfu.keyMap, lfuNode.key)
		lfu.size--
	}

	node := &node{
		key:   key,
		value: value,
		freq:  1,
	}
	lfu.keyMap[key] = node
	lfu.freqMap[1] = &doublyLinkedList{
		head: node,
		tail: node,
	}
	lfu.freqMap[1].addNode(node)
	lfu.minFreq = 1
	lfu.size++
}

func (lfu *LFUCache) Size() int {
	return lfu.size
}

func (lfu *LFUCache) Capacity() int {
	return lfu.capacity
}

func (lfu *LFUCache) MinFreq() int {
	return lfu.minFreq
}

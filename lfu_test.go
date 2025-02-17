package lfu

import (
	"testing"
)

func TestDoublyLinkedListAddNode(t *testing.T) {
	tests := []struct {
		name     string
		state    doublyLinkedList
		node     node
		expected doublyLinkedList
	}{
		{
			name:  "add node to empty list",
			state: doublyLinkedList{},
			node:  node{key: "key1", value: "value1", freq: 1},
			expected: doublyLinkedList{
				head: &node{key: "key1", value: "value1", freq: 1},
				tail: &node{key: "key1", value: "value1", freq: 1},
			},
		},
		{
			name: "add node to non-empty list",
			state: doublyLinkedList{
				head: &node{key: "key1", value: "value1", freq: 1},
				tail: &node{key: "key1", value: "value1", freq: 1},
			},
			node: node{
				key:   "key2",
				value: "value2",
				freq:  1,
			},
			expected: doublyLinkedList{
				head: &node{key: "key1", value: "value1", freq: 1, next: &node{key: "key2", value: "value2", freq: 1}},
				tail: &node{key: "key2", value: "value2", freq: 1, prev: &node{key: "key1", value: "value1", freq: 1}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.state.addNode(&test.node)
			if test.state.head.key != test.expected.head.key || test.state.head.value != test.expected.head.value {
				t.Errorf("expected head %v, got %v", test.expected.head, test.state.head)
			}
			if test.state.tail.key != test.expected.tail.key || test.state.tail.value != test.expected.tail.value {
				t.Errorf("expected tail %v, got %v", test.expected.tail, test.state.tail)
			}
		})
	}
}

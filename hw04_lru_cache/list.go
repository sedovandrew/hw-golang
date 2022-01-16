package hw04lrucache

// List is a doubly linked list interface.
type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(value interface{}) *ListItem
	PushBack(value interface{}) *ListItem
	Remove(item *ListItem)
	MoveToFront(item *ListItem)
}

// ListItem is an element of the doubly linked list.
type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

// list is a doubly linked list structure.
type list struct {
	length int
	front  *ListItem
	back   *ListItem
}

// Len returns the length of the doubly linked list.
func (l list) Len() int {
	return l.length
}

// Front returns the first element of the doubly linked list.
func (l list) Front() *ListItem {
	return l.front
}

// Back returns the last element of the doubly linked list.
func (l list) Back() *ListItem {
	return l.back
}

// PushFront adds an element to the front of a doubly linked list.
func (l *list) PushFront(value interface{}) *ListItem {
	frontListItem := l.front
	newFrontListItem := &ListItem{
		Value: value,
		Next:  frontListItem,
		Prev:  nil,
	}
	if frontListItem != nil {
		frontListItem.Prev = newFrontListItem
	}
	l.front = newFrontListItem
	if l.back == nil {
		l.back = newFrontListItem
	}
	l.length++
	return newFrontListItem
}

// PushBack adds an element to the end of a doubly linked list.
func (l *list) PushBack(value interface{}) *ListItem {
	backListItem := l.back
	newBackListItem := &ListItem{
		Value: value,
		Next:  nil,
		Prev:  backListItem,
	}
	if backListItem != nil {
		backListItem.Next = newBackListItem
	}
	l.back = newBackListItem
	if l.front == nil {
		l.front = newBackListItem
	}
	l.length++
	return newBackListItem
}

// Remove removes an element from a doubly linked list.
func (l *list) Remove(item *ListItem) {
	if item == l.front {
		l.front = item.Next
		if item.Next != nil {
			item.Next.Prev = nil
		}
	} else {
		item.Prev.Next = item.Next
	}
	if item == l.back {
		l.back = item.Prev
		if item.Prev != nil {
			item.Prev.Next = nil
		}
	} else {
		item.Next.Prev = item.Prev
	}
	l.length--
}

// MoveToFront puts an element at the beginning of a doubly linked list.
func (l *list) MoveToFront(item *ListItem) {
	// If the element is already the first one, then there is nothing to do.
	if item == l.front {
		return
	}

	// Connect neighbors with each other, if any.
	if item.Next != nil {
		item.Next.Prev = item.Prev
	}
	if item.Prev != nil {
		item.Prev.Next = item.Next
	}

	// Specify the last element of the list, if necessary.
	if l.back == item {
		l.back = item.Prev
	}

	// Put the element at the beginning of the list.
	item.Next = l.front
	item.Prev = nil
	l.front.Prev = item
	l.front = item
}

// NewList creates a doubly linked list.
func NewList() List {
	return new(list)
}

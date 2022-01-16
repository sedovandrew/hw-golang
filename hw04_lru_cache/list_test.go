package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListSimple(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("first item to the end", func(t *testing.T) {
		l := NewList()
		l.PushBack(10) // [10]
		require.Equal(t, 1, l.Len())
		require.Equal(t, 10, l.Front().Value)
		require.Equal(t, 10, l.Back().Value)
	})
}

func TestListRemove(t *testing.T) {
	t.Run("remove all elements", func(t *testing.T) {
		l := NewList()
		item10 := l.PushFront(10) // [10]
		l.Remove(item10)
		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("removing the first item from a two-item list", func(t *testing.T) {
		l := NewList()
		item10 := l.PushFront(10) // [10]
		l.PushBack(20)            // [10, 20]
		l.Remove(item10)          // [20]
		require.Equal(t, 1, l.Len())
		require.Equal(t, 20, l.Front().Value)
		require.Equal(t, 20, l.Back().Value)
		require.Nil(t, l.Front().Prev)
		require.Nil(t, l.Back().Next)
	})

	t.Run("removing the first item from a three-item list", func(t *testing.T) {
		l := NewList()
		item10 := l.PushFront(10) // [10]
		l.PushBack(20)            // [10, 20]
		l.PushBack(30)            // [10, 20, 30]
		l.Remove(item10)          // [20, 30]
		require.Equal(t, 2, l.Len())
		require.Equal(t, 20, l.Front().Value)
		require.Equal(t, 30, l.Back().Value)
		require.Nil(t, l.Front().Prev)
		require.Nil(t, l.Back().Next)
	})

	t.Run("removing the last item from a two-item list", func(t *testing.T) {
		l := NewList()
		l.PushFront(10)          // [10]
		item20 := l.PushBack(20) // [10, 20]
		l.Remove(item20)         // [10]
		require.Equal(t, 1, l.Len())
		require.Equal(t, 10, l.Front().Value)
		require.Equal(t, 10, l.Back().Value)
		require.Nil(t, l.Front().Prev)
		require.Nil(t, l.Back().Next)
	})

	t.Run("removing the last item from a three-item list", func(t *testing.T) {
		l := NewList()
		item30 := l.PushFront(30) // [30]
		l.PushFront(20)           // [20, 30]
		l.PushFront(10)           // [10, 20, 30]
		l.Remove(item30)          // [10, 20]
		require.Equal(t, 2, l.Len())
		require.Equal(t, 10, l.Front().Value)
		require.Equal(t, 20, l.Back().Value)
		require.Nil(t, l.Front().Prev)
		require.Nil(t, l.Back().Next)
	})

	// The test is disabled for the reason:
	// **Считаем, что методы Remove и MoveToFront вызываются только от существующих в списке элементов.**
	// t.Run("remove item twice", func(t *testing.T) {
	// 	l := NewList()
	// 	item10 := l.PushFront(10) // [10]
	// 	l.Remove(item10)          // []
	// 	l.Remove(item10)          // []
	// 	require.Equal(t, 0, l.Len())
	// })
}

func TestListMove(t *testing.T) {
	t.Run("move single item", func(t *testing.T) {
		l := NewList()
		item10 := l.PushFront(10) // [10]
		l.MoveToFront(item10)
		require.Equal(t, 1, l.Len())
		require.Equal(t, 10, l.Front().Value)
		require.Equal(t, 10, l.Back().Value)
	})
}

func TestListSequence(t *testing.T) {
	t.Run("sequence", func(t *testing.T) {
		// Creating List for testing
		l := NewList()
		listSize := 10
		itemLinks := make([]*ListItem, 0, listSize)
		for v := 0; v < listSize; v++ {
			itemLinks = append(itemLinks, l.PushBack(v))
		} // [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]

		// List manipulation
		l.MoveToFront(itemLinks[5]) // [5, 0, 1, 2, 3, 4, 6, 7, 8, 9]
		l.MoveToFront(itemLinks[8]) // [8, 5, 0, 1, 2, 3, 4, 6, 7, 9]
		l.MoveToFront(itemLinks[8]) // [8, 5, 0, 1, 2, 3, 4, 6, 7, 9]
		l.Remove(itemLinks[9])      // [8, 5, 0, 1, 2, 3, 4, 6, 7]
		l.Remove(itemLinks[8])      // [5, 0, 1, 2, 3, 4, 6, 7]
		l.MoveToFront(itemLinks[3]) // [3, 5, 0, 1, 2, 4, 6, 7]
		l.PushFront(10)             // [10, 3, 5, 0, 1, 2, 4, 6, 7]
		l.PushBack(11)              // [10, 3, 5, 0, 1, 2, 4, 6, 7, 11]
		l.PushFront(12)             // [12, 10, 3, 5, 0, 1, 2, 4, 6, 7, 11]

		// Checking the forward pass of the sequence
		sequence := [...]int{12, 10, 3, 5, 0, 1, 2, 4, 6, 7, 11}
		index := 0
		var prev *ListItem
		for item := l.Front(); item != nil; item = item.Next {
			// Checking values
			require.Equal(t, item.Value, sequence[index])
			index++

			// Checking the link to the previous list item
			require.Equal(t, item.Prev, prev)
			prev = item
		}

		// Checking the reverse pass of the sequence
		var next *ListItem
		for item := l.Back(); item != nil; item = item.Prev {
			// Checking the link to the next list item
			require.Equal(t, item.Next, next)
			next = item
		}
	})
}

func TestListTypes(t *testing.T) {
	t.Run("store strings", func(t *testing.T) {
		l := NewList()
		hello := l.PushFront("Hello") // ["Hello"]
		l.PushFront("World!")         // ["World!", "Hello"]
		require.Equal(t, 2, l.Len())
		require.Equal(t, "World!", l.Front().Value)
		require.Equal(t, "Hello", l.Back().Value)
		l.MoveToFront(hello) // ["Hello", "World!"]
		require.Equal(t, 2, l.Len())
		require.Equal(t, "Hello", l.Front().Value)
		require.Equal(t, "World!", l.Back().Value)
	})
}

func TestListComplex(t *testing.T) {
	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}

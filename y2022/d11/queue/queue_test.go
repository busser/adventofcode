package queue

import "fmt"

func ExampleQueue() {
	q := New[int](10)

	fmt.Println(q.Len())

	q.Enqueue(123)
	q.Enqueue(456)
	fmt.Println(q.Len())

	fmt.Println(q.Dequeue())
	fmt.Println(q.Len())

	q.Enqueue(789)
	fmt.Println(q.Dequeue())
	fmt.Println(q.Dequeue())

	// Output:
	// 0
	// 2
	// 123
	// 1
	// 456
	// 789
}

package main

import "fmt"

func main() {

	a := []int{1, 2, 3}
	b := []int{4, 5, 6}

	// APPEND
	a = append(a, b...)
	fmt.Println(a) // output: [1 2 3 4 5 6]

	// INSERT a value 99 at position i
	a = []int{1, 2, 3}
	i := 2
	a = append(a[:i], append([]int{99}, a[i:]...)...)
	fmt.Println(a) // output: [1 2 99 3]

	// DELETE 1st element i.e. i = 0
	a = []int{1, 2, 3, 4, 5, 6, 7, 8}
	i = 0
	a = append(a[:i], a[i+1:]...)
	fmt.Println(a) // output: [2 3 4 5 6 7 8]

	// DELETE last element, i = 7
	a = []int{1, 2, 3, 4, 5, 6, 7, 8}
	i = len(a) - 1
	a = append(a[:i], a[i+1:]...)
	fmt.Println(a) // output: [1 2 3 4 5 6 7]

	// DELETE 4th element, i = 3
	a = []int{1, 2, 3, 4, 5, 6, 7, 8}
	i = 3
	a = append(a[:i], a[i+1:]...)
	fmt.Println(a) // output: [1 2 3 5 6 7 8]

	// FILTER in place e.g. keep evens
	a = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	keep := func(x int) bool { return x%2 == 00 }
	n := 0
	for _, x := range a {
		if keep(x) {
			a[n] = x
			n++
		}
	}
	a = a[:n]
	fmt.Println(a) // output: [2 4 6 8 10]

	// COPY a to b
	a = []int{10, 20, 30, 40}
	b = make([]int, len(a))
	copy(b, a)
	fmt.Println(b) // output: [10 20 30 40]

	// COPY a to b, alternative
	b = append(a[:0:0], a...)
	fmt.Println(b) // output: [10 20 30 40]

	// EXPAND insert n elements at position i
	a = []int{2, 4, 6}
	n = 3
	i = 2
	fmt.Println(len(a)) // output: 3
	fmt.Println(a)      // output: [2 4 6]

	a = append(a[:i], append(make([]int, n), a[i:]...)...)
	fmt.Println(len(a)) // output: 6
	fmt.Println(a)      // output: [2 4 0 0 0 6]

	// APPEND n elements
	a = []int{2, 4, 6}
	a = append(a, make([]int, n)...)
	fmt.Println(a) // output: [2 4 6 0 0 0]

}

package fibonacci

func FibonacciActivity(n int) (int, error) {
	return fibonacciRecursive(n), nil
}

func fibonacciRecursive(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacciRecursive(n-1) + fibonacciRecursive(n-2)
}

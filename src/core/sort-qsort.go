// Package gosorts
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-01
// This copies code from Go's sort.go because we can't use something like
// sort.SortRange(data, a, b) to sort a range of data.  Wrapping incoming
// data in another sort.Interface is possible, but kills speed.
package gosorts

import (
	"sort"
)

// Insertion sort
func insertionSort(data sort.Interface, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && data.Less(j, j-1); j-- {
			data.Swap(j, j-1)
		}
	}
}

// siftDown implements the heap property on data[lo, hi).
// first is an offset into the array where the root of the heap lies.
func siftDown(data sort.Interface, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && data.Less(first+child, first+child+1) {
			child++
		}
		if !data.Less(first+root, first+child) {
			return
		}
		data.Swap(first+root, first+child)
		root = child
	}
}

func heapSort(data sort.Interface, a, b int) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown(data, i, hi, first)
	}

	// Pop elements, largest first, into end of data.
	for i := hi - 1; i >= 0; i-- {
		data.Swap(first, first+i)
		siftDown(data, lo, i, first)
	}
}

// medianOfThree moves the median of the three values data[m0], data[m1], data[m2] into data[m1].
func medianOfThree(data sort.Interface, m1, m0, m2 int) {
	// sort 3 elements
	if data.Less(m1, m0) {
		data.Swap(m1, m0)
	}
	// data[m0] <= data[m1]
	if data.Less(m2, m1) {
		data.Swap(m2, m1)
		// data[m0] <= data[m2] && data[m1] < data[m2]
		if data.Less(m1, m0) {
			data.Swap(m1, m0)
		}
	}
	// now data[m0] <= data[m1] <= data[m2]
}

func doPivot(data sort.Interface, lo, hi int) (midlo, midhi int) {
	m := lo + (hi-lo)/2 // Written like this to avoid integer overflow.
	if hi-lo > 40 {
		// Tukey's ``Ninther,'' median of three medians of three.
		s := (hi - lo) / 8
		medianOfThree(data, lo, lo+s, lo+2*s)
		medianOfThree(data, m, m-s, m+s)
		medianOfThree(data, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThree(data, lo, m, hi-1)

	// Invariants are:
	//	data[lo] = pivot (set up by ChoosePivot)
	//	data[lo < i < a] < pivot
	//	data[a <= i < b] <= pivot
	//	data[b <= i < c] unexamined
	//	data[c <= i < hi-1] > pivot
	//	data[hi-1] >= pivot
	pivot := lo
	a, c := lo+1, hi-1

	for ; a != c && data.Less(a, pivot); a++ {
	}
	b := a
	for {
		for ; b != c && !data.Less(pivot, b); b++ { // data[b] <= pivot
		}
		for ; b != c && data.Less(pivot, c-1); c-- { // data[c-1] > pivot
		}
		if b == c {
			break
		}
		// data[b] > pivot; data[c-1] <= pivot
		data.Swap(b, c-1)
		b++
		c--
	}
	// If hi-c<3 then there are duplicates (by property of median of nine).
	// Let be a bit more conservative, and set border to 5.
	protect := hi-c < 5
	if !protect && hi-c < (hi-lo)/4 {
		// Lets test some points for equality to pivot
		dups := 0
		if !data.Less(pivot, hi-1) { // data[hi-1] = pivot
			data.Swap(c, hi-1)
			c++
			dups++
		}
		if !data.Less(b-1, pivot) { // data[b-1] = pivot
			b--
			dups++
		}
		// m-lo = (hi-lo)/2 > 6
		// b-lo > (hi-lo)*3/4-1 > 8
		// ==> m < b ==> data[m] <= pivot
		if !data.Less(m, pivot) { // data[m] = pivot
			data.Swap(m, b-1)
			b--
			dups++
		}
		// if at least 2 points are equal to pivot, assume skewed distribution
		protect = dups > 1
	}
	if protect {
		// Protect against a lot of duplicates
		// Add invariant:
		//	data[a <= i < b] unexamined
		//	data[b <= i < c] = pivot
		for {
			for ; a != b && !data.Less(b-1, pivot); b-- { // data[b] == pivot
			}
			for ; a != b && data.Less(a, pivot); a++ { // data[a] < pivot
			}
			if a == b {
				break
			}
			// data[a] == pivot; data[b-1] < pivot
			data.Swap(a, b-1)
			a++
			b--
		}
	}
	// Swap pivot into middle
	data.Swap(pivot, b-1)
	return b - 1, c
}

func quickSort(data sort.Interface, a, b, maxDepth int) {
	for b-a > 12 {
		if maxDepth == 0 {
			heapSort(data, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivot(data, a, b)
		// Avoiding recursion on the larger subproblem guarantees
		// a stack depth of at most lg(b-a).
		if mlo-a < b-mhi {
			quickSort(data, a, mlo, maxDepth)
			a = mhi // i.e., quickSort(data, mhi, b)
		} else {
			quickSort(data, mhi, b, maxDepth)
			b = mlo // i.e., quickSort(data, a, mlo)
		}
	}
	if b-a > 1 {
		// Do ShellSort pass with gap 6
		// It could be written in this simplified form cause b-a <= 12
		for i := a + 6; i < b; i++ {
			if data.Less(i, i-6) {
				data.Swap(i, i-6)
			}
		}
		insertionSort(data, a, b)
	}
}

// qSort quicksorts data immediately.
// It performs O(n*log(n)) comparisons and swaps. The sort is not stable.
func qSort(data sort.Interface, a, b int) {
	// Switch to heapsort if depth of 2*ceil(lg(n+1)) is reached.
	n := b - a
	maxDepth := 0
	for i := n; i > 0; i >>= 1 {
		maxDepth++
	}
	maxDepth *= 2
	quickSort(data, a, b, maxDepth)
}

// Quicksort performs a parallel quicksort on data.
func Quicksort(data sort.Interface) {
	a, b := 0, data.Len()
	n := b - a
	maxDepth := 0
	for i := n; i > 0; i >>= 1 {
		maxDepth++
	}
	maxDepth *= 2
	parallelSort(data, quickSortWorker, task{-maxDepth - 1, a, b})
}

// qSortPar starts a parallel quicksort.
func qSortPar(data sort.Interface, t task, sortRange func(task)) {
	a, b := t.pos, t.end
	n := b - a
	maxDepth := 0
	for i := n; i > 0; i >>= 1 {
		maxDepth++
	}
	maxDepth *= 2
	quickSortWorker(data, task{-maxDepth - 1, a, b}, sortRange)
}

// quickSortWorker is a parallel analogue of quickSort: it performs a pivot
// and might asynchronously sort one of the halves if it's large enough.
func quickSortWorker(data sort.Interface, t task, sortRange func(task)) {
	maxDepth, a, b := 1-t.offs, t.pos, t.end
	for b-a > minOffload {
		if maxDepth == 0 {
			heapSort(data, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivot(data, a, b)
		// Avoiding recursion on the larger subproblem guarantees
		// a stack depth of at most lg(b-a).
		if mlo-a < b-mhi {
			sortRange(task{-maxDepth - 1, a, mlo})
			a = mhi // i.e., quickSortWorker(data, mhi, b)
		} else {
			sortRange(task{-maxDepth - 1, mhi, b})
			b = mlo // i.e., quickSortWorker(data, a, mlo)
		}
	}
	if b-a > 7 {
		quickSort(data, a, b, maxDepth)
	} else if b-a > 1 {
		insertionSort(data, a, b)
	}
}

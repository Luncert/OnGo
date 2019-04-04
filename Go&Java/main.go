package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

func main() {
	content, err := ioutil.ReadFile("data")
	if err != nil {
		panic(err)
	}

	data := ParseIntArray(strings.Split(string(content), ","))
	start := time.Now()
	QuickSort(data)
	ok := check(data)
	fmt.Printf("%vms %v\n", time.Now().Sub(start).Nanoseconds()/1e6, ok)
}

func ParseIntArray(raw []string) []int64 {
	v := make([]int64, len(raw))
	for i, s := range raw {
		v[i] = ParseInt([]byte(s))
	}
	return v
}

func ParseInt(bytes []byte) (v int64) {
	var b int64 = 0
	// all value is positive
	for _, i := range bytes {
		b = v*10 + int64(i-'0')
		if b < v {
			panic("overflow!")
		}
		v = b
	}
	return
}

func QuickSort(data []int64) {
	quickSort(data, 0, len(data)-1)
}

func quickSort(data []int64, start, end int) {
	if end-start == 1 {
		if data[start] > data[end] {
			data[start], data[end] = data[end], data[start]
		}
	} else if end > start {
		v := data[start]
		i, j := start+1, end
		// move i to find a value bigger than v
		for ; i < j; i++ {
			// found it!
			if data[i] > v {
				// move j to find a value smaller than v
				for ; i < j; j-- {
					// find it
					if data[j] < v {
						data[i], data[j] = data[j], data[i]
						break
					}
				}
				// i == j: no compatiable value could be found by moving j!
				// if i has never been moved, that means all the value after data[start], from i to end are bigger
				// if not, that means all the value before data[i] is smaller
				if i == j {
					if i == start+1 {
						// i has never been moved
						quickSort(data, i, end)
					} else {
						data[start] = data[i-1]
						data[i-1] = v
						quickSort(data, start, i-1)
						quickSort(data, i, end)
					}
					return
				}
			}
		}
		if data[i] >= v {
			i--
		}
		data[start], data[i] = data[i], data[start]
		quickSort(data, start, i)
		quickSort(data, i+1, end)
	}
}

func check(data []int64) bool {
	for i, limit := 0, len(data)-1; i < limit; i++ {
		if data[i] > data[i+1] {
			return false
		}
	}
	return true
}

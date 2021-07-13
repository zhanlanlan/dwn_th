package main

import (
	"dwn_th/route"

	"github.com/gin-gonic/gin"
)

func main() {

	{
		route.InitClients()
	}

	{
		server := gin.Default()
		route.Route(server)
		server.Run()
	}

}

// func bubblesort(s []int, m int) {
// 	if m == 1 {
// 		return
// 	}

// 	for i := 0; i < m-1; i++ {
// 		if s[i] > s[i+1] {
// 			swap(s, i, i+1)
// 		}
// 	}
// 	bubblesort(s, m-1)
// }

// func insertsort(s []int) {

// 	if len(s) < 2 {
// 		return
// 	}

// 	for i := 1; i < len(s); i++ {
// 		for j := i; j > 0 && s[j] < s[j-1]; j-- {
// 			swap(s, j, j-1)
// 		}
// 	}
// }

// // arr := []int{9, 4, 5, 11, 22, 8, 4, 18, 0, 2}
// // quickSort(arr, 0, 9)
// func partition(arr []int, low, high int) int {
// 	pivot := arr[low]
// 	for low < high {
// 		for low < high && pivot <= arr[high] {
// 			high--
// 		}
// 		arr[low] = arr[high]

// 		for low < high && pivot >= arr[low] {
// 			low++
// 		}
// 		arr[high] = arr[low]
// 	}
// 	arr[low] = pivot
// 	return low
// }
// func quickSort(arr []int, low, high int) {
// 	if low >= high {
// 		return
// 	}
// 	p := partition(arr, low, high)
// 	quickSort(arr, low, p-1)
// 	quickSort(arr, p+1, high)
// }

// func selectSort(arr []int, n int) {
// 	for i := 0; i < n; i++ {
// 		min := i
// 		for j := i + 1; j < n; j++ {
// 			if arr[j] < arr[min] {
// 				min = j
// 			}
// 		}
// 		swap(arr, i, min)
// 	}
// }

// func swap(s []int, i, j int) {
// 	s[i], s[j] = s[j], s[i]
// }

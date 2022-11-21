package main

import (
	"fmt"
	"math"
	"sync"
)

// takes an indefnite amount of funcs to apply a waitgroup for each of them so they all wait for each other
func Par(procs ...func()) {
	var wg sync.WaitGroup // implementation of WaitGroup
	wg.Add(len(procs))    //for every func() wg.add()
	defer wg.Wait()       // Wait for each func() to finish
	for _, proc := range procs {
		go func(proc func()) {
			defer wg.Done()
			proc()
		}(proc)
	}
}


// Returns the x number then, forever Takes a number and returns it
func Prefix(x int, in, out chan int) {
	buf := x
	out <- buf
	for {
		buff := <-in
		out <- buff
	}
}

//Compares 2 integers and returns the higher one in out_h and lower one in out_l
func Comp(in_x, in_y, out_h, out_l chan int) {
	for {

		bufhigh := <-in_x
		buflow := <-in_y

		// compare and return accordingly 
		if bufhigh < buflow { 	// other way
			Par(func() {
				out_l <- bufhigh
			}, func() {
				out_h <- buflow
			})
		} else { 				// right way
			Par(func() {
				out_l <- buflow
			}, func() {
				out_h <- bufhigh
			})
		}

	}
}

//The repetative component of the Pump Sort
func PSort(in, out chan int) {

	a := make(chan int)
	b := make(chan int)

	go Comp(a, in, b, out)         // in,a > comp() > b,out
	go Prefix(math.MinInt64, b, a) // b > prefix(0) > a

	/*
	* compare the entered (in channel) number and held (a channel) number
	* and push the larger out (out channel) while keeping the smaller number (b channel)
	 */
}

// Returns a sorted array of the entered array using Pump Sort principles
func PumpSort(array []int) []int {

	//allocates the right amount of channels
	in := make(chan int)
	ch := make(chan int, len(array)-2)
	out := make(chan int)

	//pre-connection,  in > PSort() > ch
	go PSort(in, ch)

	for i := 1; i < (len(array) - 1); i++ { // for loop for each compnent in sort,  ch > PSort() > ch
		go PSort(ch, ch)
	}

	go PSort(ch, out)
	//post-connection  ch > PSort() > out

	for i := 0; i < len(array); i++ {
		in <- array[i] //pump every number in the array
		<-out          //dump every useless number of Min values
	}

	sortedArray := make([]int, len(array)) //crate the new array

	for i := 0; i < len(array); i++ {
		in <- math.MaxInt64    // flush the system using Max values
		sortedArray[i] = <-out //save the number in the array
	}

	return sortedArray // return the new sorted array

}

func PumpSortRevised(size int, in , out chan int){
	ch := make(chan int, size-2)

	//pre-connection,  in > PSort() > ch
	go PSort(in, ch)

	for i := 1; i < (size - 1); i++ { // for loop for each compnent in sort,  ch > PSort() > ch
		go PSort(ch, ch)
	}

	go PSort(ch, out)
}

func main(){
 	/* pump Sort implementation with only int */

	unsortedArray := [20]int{10, 4, 8, 3, 9, 2, 7, 1, 8, 3, 7, 4, 0, 1, 6, 10, 9, 2, 5, 0}

	sortedArray := PumpSort(unsortedArray[:])

	for i := 0; i < len(sortedArray); i++ {
		fmt.Println(sortedArray[i])
	}

	//in := make(chan int)
	//out := make(chan int)



}
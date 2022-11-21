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
func Prefix[T Others](x T, in, out chan T) {
	buf := x
	out <- buf
	for {
		buff := <-in
		out <- buff
	}
}

//Compares 2 integers and returns the higher one in out_h and lower one in out_l
func Comp[T Others](in_x, in_y, out_h, out_l chan T) {
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
		in <- math.MaxInt64   // flush the system using Max values
		sortedArray[i] = <-out //save the number in the array
	}

	return sortedArray // return the new sorted array

}

//Implementing Generics

type Others interface{
	int | float32 | float64 | int64 
}

func PumpSortRevised[T Others](array []T, in , out chan T){
	ch := make(chan T, len(array)-2)
	
	min := FindMin(array[:])
	//pre-connection,  		in > PSort() > ch
	go PSortRevised(min,in, ch)

	for i := 1; i < (len(array) - 1); i++ { // for loop for each compnent in sort,  ch > PSort() > ch
		go PSortRevised(min,ch, ch)
	}

	go PSortRevised(min,ch, out) 		// ch > PSort() > out
}

func PSortRevised[T Others](min T,in , out chan T){
	a := make(chan T)
	b := make(chan T)

	go Comp(a, in, b, out)         // in,a > comp() > b,out
	go Prefix(min, b, a) // b > prefix(0) > a
}

func FindMin[T Others](array []T) T{
	min := array[0]
	for i := 1; i < len(array); i++ {
		if(min > array[i]){
			min = array[i]
		}
	}
	return min
}

func main(){
 	/* pump Sort implementation with only int */

	unsortedArray := [20]int{10, 4, 8, 3, 9, 2, 7, 1, 8, 3, 7, 4, 0, 1, 6, 10, 9, 2, 5, 0}
	fmt.Print("unsorted int\t> ")
	fmt.Println(unsortedArray)

	sortedArray := PumpSort(unsortedArray[:])
	fmt.Print("sorted int\t> ")
	fmt.Println(sortedArray)
	
	//pump sort Implementation such that floats are accepted as well to represet non-integers

	in := make(chan float64)
	out := make(chan float64)
	
	unsortedArray2 := [20]float64{10, 4, 8, 3, 9, 2, 7, 1, 8, 3, 7, 4, 0, 1, 6, 10, 9, 2, 5, 0}
	fmt.Print("unsorted float64> ")
	fmt.Println(unsortedArray2)

	go PumpSortRevised(unsortedArray2[:],in,out)

	for i := 0; i < len(unsortedArray2); i++ {
		in <- unsortedArray2[i]
		<-out 
	}

	sortedArray2 := make([]float64,0)

	for i := 0; i < len(unsortedArray2); i++ {
		in <- math.MaxInt64
		sortedArray2 = append(sortedArray2,<-out)
	}
	fmt.Print("sorted float64\t> ")
	fmt.Println(sortedArray2)

}
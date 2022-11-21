package main

import (
	"fmt"
	"math"
	"sync"
)

// Takes a number and returns it
func Id(in, out chan int) {
	for {
		buf := <-in
		out <- buf
	}
}

// Takes a number, add by 1 and returns it
func Succ(in, out chan int) {
	for {
		buf := <-in
		buf++
		out <- buf
	}
}

// Takes a number and dumps it
func Sink(in chan int) {
	for {
		<-in
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

// Takes a number dumps it then, forever Takes a number and returns it
func Tail(in, out chan int) {
	<-in
	for {
		buff := <-in
		out <- buff
	}
}

//Takes 2 numbers, adds them then, returns the sum
func Plus(in_x, in_y, out chan int) {
	for {
		x, y := 0, 0
		Par(func() {
			x = <-in_x
		}, func() {
			y = <-in_y
		})
		out <- x + y
	}
}

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

//takes 1 number and returns the same number out twice
func delta(in, out_x, out_y chan int) {
	for {
		buf := 0
		buf = <-in
		Par(func() {
			out_x <- buf
		}, func() {
			out_y <- buf
		})
	}
}

// only returns numbers > pattern is the Natural Numbers
func Nos(out chan int) {

	a := make(chan int)
	b := make(chan int)
	c := make(chan int)

	go Prefix(0, a, b)  // 1 > Prefix(0) > b
	go delta(b, c, out) // b > delta() > c, out 
	go Succ(c, a)		// c > Succ() > a

}

// Takes a number and adds it the next number, each number sum is also returned > pattern is total sum of all numbers entered
func Int(in, out chan int) {

	a := make(chan int)
	b := make(chan int)
	c := make(chan int)

	go Plus(in, a, b) 	// in, a > Plus() > b
	go delta(b, c, out)	// b > delta() > c, out 
	go Prefix(0, c, a) 	// c > Prefix(0) > a

}

// takes a number stream of integers and returns the pairs of each
func Pairs(in, out chan int) {

	a := make(chan int)
	b := make(chan int)
	c := make(chan int)

	go delta(in, a, b) 	// in > delta() > a, b
	go Tail(a, c) 	   	// a > Tail() > c
	go Plus(b, c, out)	// b, c > Plus() > out 

}

//returns a fibonacci sequence
func Fib(out chan int) {

	a := make(chan int)
	b := make(chan int)
	c := make(chan int)
	d := make(chan int)

	go Prefix(1, d, a)  // d > pref(1) > a
	go Prefix(0, a, b)  // a > pref(0) > b
	go delta(b, c, out) // b > delt > c > out
	go Pairs(c, d)      // c > pair > d

}

//returns a set of squares
func Squares(out chan int) {

	a := make(chan int)
	b := make(chan int)

	Nos(a) 			// Nos() > a
	Int(a, b)		// a > Int() > b
	Pairs(b, out)	// b > Pairs() > out
}

//Compares 2 integers and returns the higher one in out_h and lower one in out_l
func Comp(in_x, in_y, out_h, out_l chan int) {
	for {

		bufhigh, bufLow := 0, 0

		// take the 2 in 
		Par(func() {
			bufhigh = <-in_x
		}, func() {
			bufLow = <-in_y
		})

		// compare and return accordingly 
		if bufhigh < bufLow { 	// other way
			Par(func() {
				out_l <- bufhigh
			}, func() {
				out_h <- bufLow
			})
		} else { 				// right way
			Par(func() {
				out_l <- bufLow
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

	go Comp(in, a, b, out)         // in,a > comp() > b,out
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

func main(){
	/* Squares Implementation */
	//out := make(chan int)

	go Squares(out)

	for i := 0; i < 20; i++ {
		fmt.Println(<-out)
	}
	

	/* Fib Implemenation
	//out := make(chan int)

	go Fib(out)

	for i := 0; i < 10; i++ {
		fmt.Println(<-out)
	}

	*/

	/* Pairs Implementation
	//out := make(chan int)
	//in := make(chan int)

	go Pairs(in, out)

	in <- 0
	for i := 1; i < 10; i++ {
		in <- i
		fmt.Println(<-out)
	}

	*/

	/* Int Implementation
	//out := make(chan int)
	//in := make(chan int)

	go Int(in, out)

	for i := 0; i < 10; i++ {
		in <- i
		fmt.Println(<-out)
	}

	*/

	/* Nos Implementation
	//out := make(chan int)

	go Nos(out)

	for i := 0; i < 10; i++ {
		fmt.Println(<-out)
	}

	*/
}
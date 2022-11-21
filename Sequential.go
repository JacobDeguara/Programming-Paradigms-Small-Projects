// current issue my go not connecting have to use this before running anything "export PATH=$PATH:/usr/local/go/bin"

package main

import (
	"fmt"
)

func Fact(n int) int{
	if(n < 2){ // Base Case
		return 1 // if 1 or less return 1
	}else {
		return n + Fact(n-1) // else n + (n-1)! 
	}
}

func Fact2(n int) int{
	if(n < 2){ // Base Case
		return 1 // return 1 if les then 2
	}else{
		return Fact3(n-1,n) // Indirect Recurtion
	}
}

func Fact3(n int, acc int) int{
	return acc + Fact2(n) // same as n + (n-1)!
}

func Fib(n int) int{
	if(n < 1){ // Base Case for less then 0
		return 0 // if less then 1 return 0
	}
	if(n == 1 || n ==2){ // Base Case for 1 & 2
		return 1 // if 1 or 2 return 1
	}else{
		return Fib(n-1) + Fib(n-2) // else recrution of fib(n-1) + fib(n-2)
	}
}

func Fib2(n int) int{
	if(n < 1){ // Base Case
		return 0 // return 0 if less then 1
	}else{
		return Fib3(n,n-2,n-1) // Indirect Recursion
	}
}

func Fib3(n int, n_minus2 int , n_minus1 int) int{
	if(n_minus1 == 1){ // base case for -1
		return 1
	}else if(n_minus2 ==1){ // base case for -2
		return 2
	}else {
		return Fib2(n_minus2) + Fib2(n_minus1) // recursive
	}
}

//Filters an array based on a function that takes an interger from the array and returns true if accepted
func Filter(f func(int) bool, elems []int) []int{
	result := make([]int,0) // create new array

	for i := 0 ; i < len(elems) ; i++{
		if(f(elems[i])){ // check on func f if true
			result = append(result,elems[i]) // keep element
		} 
	} //repeat

	return result //return result when done
}

//returns true if divisible by 2
func isDivisibleBy2(n int) bool{ 
	if(n%2 == 0){
		return true
	}
	return false
}

//returns the reverse of the array 
func Reverse(elems []int) []int{
	reverse := make([]int,0) // create new array

	for j := len(elems)-1 ; j > -1 ; j--{ // iterate backwards
		reverse = append(reverse, elems[j]) // append every element
	}

	return reverse //return reverse
}


func main() {
	fmt.Println(Fact(3)) 	// 6
	fmt.Println(Fact2(3)) 	// 6
	fmt.Println(Fib(10)) 	// 55
	fmt.Println(Fib2(10))	// 55

	x := []int{1,2,3,4,5,6,7,8,9,10}
	fmt.Println(x)

	y := Filter(isDivisibleBy2,x)
	fmt.Println(y)

	z := Reverse(x)
	fmt.Println(z)

}

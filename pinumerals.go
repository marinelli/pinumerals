package main

import "fmt"



type Name chan Name

var empty Name



func output_act (x Name, y Name) {
	x <- y
}



func input_act (x Name, y *Name) {
	*y = <-x
}



func numeral_from_int (n int, x Name, z Name) (func ()) {
	return func () {
		for i := 0; i < n; i++ {
			output_act (x, empty)
		}

		output_act (z, empty)
	}
}



func int_from_numeral (x Name, z Name) (func () int) {
	return func () int {

		var result int = 0

		loop: for {
			select {
			case <- x : result += 1
			case <- z : break loop
			}
		}

		return result
	}
}



func numeral_from_int_rec (n int, x Name, z Name) (func ()) {

	var aux func (result int, x Name, z Name)

	aux = func (result int, x Name, z Name) {
		switch {
		case n == 0 : output_act(z, empty)
		case n > 0  : output_act(x, empty) ; aux(result - 1, x, z)
		}
	}

	return func () {
		aux (n, x, z)
	}
}



func int_from_numeral_rec (x Name, z Name) (func () int) {

	var aux func (result int, x Name, z Name) int

	aux = func (result int, x Name, z Name) int {
		select {
		case <- x : return aux (result + 1, x, z)
		case <- z : return result
		}
	}

	return func () int {
		return aux (0, x, z)
	}
}



func Copy(x, z, y, w Name) {
	select {
	case <- x : Succ (x, z, y, w)
	case <- z : output_act (w, empty)
	}
}



func Succ(x, z, y, w Name) {
	output_act (y, empty)
	Copy (x, z, y, w)
}



func Add(x1, z1, x2, z2, y, w Name) {
	select {
	case <- x1 : output_act (y, empty) ; Add (x1, z1, x2, z2, y, w)
	case <- z1 : Copy (x2, z2, y, w)
	}
}



func main () {

	x1, z1 := make(Name), make(Name)
	x2, z2 := make(Name), make(Name)
	x3, z3 := make(Name), make(Name)

	n1 := numeral_from_int (3, x1, z1)
	n2 := numeral_from_int (7, x2, z2)

	go n1()
	go n2()

	i1 := int_from_numeral (x1, z1)
	i2 := int_from_numeral (x2, z2)

	fmt.Print (i1(), " + ", i2(), " = ")

	go numeral_from_int (3, x1, z1) ()
	go numeral_from_int (7, x2, z2) ()
	go Add (x1, z1, x2, z2, x3, z3)

	fmt.Println (int_from_numeral (x3, z3) ())
}


package main

import (
	"fmt"
	//"log"
	//"time"
)

type Person struct {
	firstName string
	lastName  string
	age       int
	address   string
}

type Option func(*Person)

func WithAddress(address string) Option {
	return func(p *Person) {
		p.address = address
	}
}

func WithAge(age int) Option {
	return func(p *Person) {
		p.age = age
	}
}

func NewPerson(firstName, lastName string, option ...Option) *Person {
	p := &Person{firstName: firstName, lastName: lastName}
	for _, option := range option {
		option(p)
	}
	return p
}

func main() {

	fmt.Println("start\n")

	p := NewPerson("John", "Doe", WithAddress("Tokyo"), WithAge(30))
	//fmt.Println(now.Format("2006/01/0215:04:05"))
	fmt.Printf("%+v\n", p)

	fmt.Println("\nend")
}

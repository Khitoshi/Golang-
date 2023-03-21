package main

import (
	"fmt"
)

type Person struct {
	Name    string
	Age     int
	Address string
}

type PersonBuilder struct {
	name    string
	age     int
	address string
}

func (pb *PersonBuilder) SetName(name string) *PersonBuilder {
	pb.name = name
	return pb
}

func (pb *PersonBuilder) SetAge(age int) *PersonBuilder {
	pb.age = age
	return pb
}

func (pb *PersonBuilder) SetAddress(address string) *PersonBuilder {
	pb.address = address
	return pb
}

func (pb *PersonBuilder) Build() *Person {
	return &Person{
		Name:    pb.name,
		Age:     pb.age,
		Address: pb.address,
	}
}

func main() {
	fmt.Println("start")

	pb := &PersonBuilder{}
	person := pb.SetName("Tom").SetAge(15).SetAddress("japan").Build()

	fmt.Printf("%#v\n", person)

	fmt.Println("end")

}

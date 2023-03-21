package main

import "fmt"

type Attr struct {
	Name string
	Age  int
}

type Teacher struct {
	Attr
	Subject string
}

type Student struct {
	Attr
	Score int
}

type AttrEX struct {
	Name string
}

type TeacherEx struct {
	Attr
	AttrEX
	Subject string
}

func (t *TeacherEx) Stringer() string {
	return fmt.Sprintf("TeacherEx: %v", t)
}

func (t *TeacherEx) OutPutName() string {
	return fmt.Sprintf("TeacherEx.Name: %v", t.Attr.Name)
}

func main() {
	fmt.Println("Starting\n")

	teacher := Teacher{
		Attr: Attr{
			Name: "John",
			Age:  20,
		},
		Subject: "Math",
	}
	fmt.Printf("%#v\n", teacher)

	student := Student{
		Attr: Attr{
			Name: "jack",
			Age:  12,
		},
		Score: 55,
	}
	fmt.Printf("%#v\n", student)

	teacherex := TeacherEx{
		Attr: Attr{
			Name: "john",
			Age:  42,
		},
		AttrEX: AttrEX{
			Name: "jack",
		},
		Subject: "Math",
	}
	fmt.Printf("%#v\n", teacherex)

	fmt.Println(teacherex.OutPutName())

	fmt.Println("\nDone")
}

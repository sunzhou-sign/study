package main

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"protoBuf/entity"
)

func main() {
	t := &entity.Student{
		Name:   "A",
		Gender: entity.Student_FEMALE,
		Scores: []int32{98, 84, 88},
	}

	data, err := proto.Marshal(t)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	nt := &entity.Student{}
	err = proto.Unmarshal(data, nt)
	if err != nil {
		log.Fatal("unmarshaline error: ", err)
	}

	if t.GetName() != nt.GetName() {
		log.Fatal("data mismatch %q != %q", t.GetName(), nt.GetName())
	} else {
		fmt.Println("success!")
	}
}

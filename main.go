package main

import (
	"github.com/evanphx/json-patch"
	"fmt"
	"log"
)

func main() {
	d1 := `{
   "a": [1, 2]
}`
	p1 := `[{"op":"add","path":"/a/0","value":3}]`

	patch, err := jsonpatch.DecodePatch([]byte(p1))
	if err != nil {
		log.Fatal(err)
	}

	o1, err := patch.Apply([]byte(d1))
	fmt.Println(err)
	fmt.Println(string(o1))
}

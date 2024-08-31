package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	s1 := `{"ggg": 1}`
	s2 := `{sdsds:dadaa, 2}`

	fmt.Println(json.Valid([]byte(s1)))
	fmt.Print(json.Valid([]byte(s2)))

}

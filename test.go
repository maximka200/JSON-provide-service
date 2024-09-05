package main

import (
	"errors"
	"fmt"
)

var (
	errInternal  = errors.New("internal err")
	errInternalT = errors.New("internal err")
)

func main() {
	errorOne := errInternal
	err := errInternalT
	fmt.Println(errorOne == err)
	fmt.Println(errors.Is(errorOne, err))
}

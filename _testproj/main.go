package main

import (
	"fmt"
	"grr/grr"
	"os"
)

func main() {
	num, err := ReturnSomethingOrError()

	if err != nil {
		grr.Errorf("MyError: this is your number: %v", num).
			AddError(err)
	}

	fmt.Println("Number: ")

}

func ReturnSomethingOrError() (int, error) {
	file, err := os.Open("filename")

	if err != nil {
		return 0, err
	}

	defer file.Close()

	return 1, nil
}

package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("start step 1");

	file, err := os.Open("./tests/step1/valid.json");

	if err != nil {
		fmt.Println("error opening valid file in step 1");
	}

	parse(file);
}

func parse(f *os.File) {


}

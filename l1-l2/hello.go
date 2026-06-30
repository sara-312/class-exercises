// Every executable Go Program should contain a package called main.
// This tells the Go compiler to compile the package into an executable
// program rather than a shared library.
package main

import (
	"fmt"
	"os"
)

// The entry point of a Go program should be the main function of main package.
// When the executable is run, main() is automatically called.
// func main() {
// 	fmt.Println("Hello World\n")
// }

func main() {
// input := "There once was a cat named Barry. He was a very good cat. This cat lived in Boston. He loved doing Boston-related activities (that were good for cats). He walked the esplanade. He shopped on Newbury. He ate at Tatte. He sometimes even went to TD Garden. Did you know that cats are not allowed in TD Garden?"
word := "cat"
data, err := os.ReadFile("dictionary.txt")
if err != nil {
	fmt.Printf("Error reading file: %v\n", err)
	return
}
input := string(data)

for i := 0; i < len(input) - len(word); i++ {
	if input[i:i+len(word)] == word {
		fmt.Printf(input[i:i+len(word)])
		fmt.Printf("Found '%s' at index %d\n", word, i)
	}
}
}
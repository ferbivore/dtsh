package main

import "bufio"
import "fmt"
import "os"

import "github.com/ferbivore/dtsh"

func main() {
	fmt.Printf("> ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := dtsh.Tokenize(line)
		for i, token := range tokens {
			fmt.Printf("%d", i)
			switch token.Type {
			case dtsh.TokenRegular:
				fmt.Printf(" reg ")
			case dtsh.TokenLiteral:
				fmt.Printf(" lit ")
			}
			fmt.Println(token.Value)
		}
		fmt.Printf("> ")
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

type Token struct {
	token   TokenType
	literal string
	// pos     uint64
}

type Stack struct {
	items []Token
}

func (s *Stack) Push(t Token) {
	s.items = append(s.items, t);
}

func (s *Stack) Pop() Token {
	lastItem := s.items[len(s.items)-1:];
	s.items = s.items[0:len(s.items) - 1];
	return lastItem[0];
}

func main() {
	file, err := os.Open("./tests/step1/valid.json")

	if err != nil {
		fmt.Println("error opening valid file in step 1")
	}

	tokens := scan(file);
	valid := isValidJSON(tokens);

	fmt.Println("is valid json?:", valid);

	if valid {
		os.Exit(0);
	} else {
		os.Exit(1);
	}
}

func createToken(char string, token TokenType) Token {
	return Token{ token: token, literal: char };
}

func scan(f *os.File) []Token {
	var tokens []Token;

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanBytes)

	for scanner.Scan() {
		char := scanner.Text()
		switch char {
			case "{":
				tokens = append(tokens, createToken(char, RBRACE));
			case "}":
				tokens = append(tokens, createToken(char, LBRACE));
		}
	}

	return tokens;
}

func isValidJSON(tokens []Token) bool {
	var s Stack;

	if len(tokens) == 0 {
		return false;
	}

	valid := true;
	for i := 0; i < len(tokens); i ++ {
		token := tokens[i];

		if token.token == LBRACE {
			corresponding := s.Pop();

			if corresponding.token != RBRACE {
				valid = false;
			}
		} else {
			s.Push(token);
		}
	}

	if len(s.items) > 0 {
		return false;
	}

	return valid;
}

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

func (s *Stack) PeekBack(steps int) Token {
	lastItem := s.items[len(s.items)-steps:];
	return lastItem[0];
}

func main() {
	file, err := os.Open("./tests/step2/invalid2.json")

	if err != nil {
		fmt.Println("error opening valid file in step 1")
	}

	tokens := scan(file);

	fmt.Println("tokens: ", tokens);

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
			case "\"":
				// continue scanning until you reach the corresponding "
				literal := "";
				for scanner.Scan() {
					c := scanner.Text();
					if c == "\"" {
						break;
					}
					literal += c;
				}

				tokens = append(tokens, createToken(literal, STRING));
			case ":":
				tokens = append(tokens, createToken(char, COLON));
			case ",":
				tokens = append(tokens, createToken(char, COMMA));
			case " ", "\n":
				// keep going, dont care about spaces or newline characters (yet?)
			default:
				fmt.Printf("Err: unknown token '%s'\n", char);
				os.Exit(1);
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
			// not handling nested objects yet

			// left brace only appears as last token
			peekedValue := s.PeekBack(1);

			if peekedValue.token != STRING && peekedValue.token != RBRACE {
				valid = false;
				break;
			}
		} else if token.token == STRING {
			peekedValue := s.PeekBack(1);

			if peekedValue.token == RBRACE || peekedValue.token == COMMA {
				s.Push(token);
			} else if peekedValue.token == COLON {
				// its a string value
				key := s.PeekBack(2);

				if key.token != STRING {
					valid = false;
				} else {
					s.Push(token);
				}
			}
		}  else {
			s.Push(token);
		}
	}

	return valid;
}

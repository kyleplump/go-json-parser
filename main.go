package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
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
	file, err := os.Open("./tests/step3/invalid.json")

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

	reader := bufio.NewReader(f)
	// scanner.Split(bufio.ScanBytes)

	for {
		byte, err := reader.ReadByte();

		if err != nil {
			break;
		}

		char := string(byte);

		switch char {
			case "{":
				tokens = append(tokens, createToken(char, RBRACE));
			case "}":
				tokens = append(tokens, createToken(char, LBRACE));
			case "\"":
				// continue scanning until you reach the corresponding "
				literal := "";
				for {
					b, err := reader.ReadByte();

					if err != nil {
						fmt.Println("encountered error: ", err);
						os.Exit(1);
					}

					c := string(b);

					if c == "\"" {
						break;
					}
					literal += c;
				}

				tokens = append(tokens, createToken(literal, STRING));
			// booleans
			case "t", "f":
				literal := char;

				for {
					b, _ := reader.ReadByte();
					c := string(b);
					literal += c

					if strings.Contains("true", literal) || strings.Contains("false", literal) {
						if literal == "true" || literal == "false" {
							tokens = append(tokens, createToken(literal, BOOL));
							break;
						}
					} else {
						fmt.Printf("Err: unknown token '%s'\n", literal);
						os.Exit(1);
					}
				}
			// null
			case "n":
				literal := char;

				for {
					b, _ := reader.ReadByte();
					c := string(b);
					literal += c

					if strings.Contains("null", literal) {
						if literal == "null" {
							tokens = append(tokens, createToken(literal, NULL));
							break;
						}
					} else {
						fmt.Printf("Err: unknown token '%s'\n", literal);
						os.Exit(1);
					}
				}
			case ":":
				tokens = append(tokens, createToken(char, COLON));
			case ",":
				tokens = append(tokens, createToken(char, COMMA));
			case " ", "\n":
				// keep going, dont care about spaces or newline characters (yet?)
			default:
				// numbers
				_, err := strconv.Atoi(char);
				createdNumericToken := false;

				if err != nil {
					fmt.Printf("Err: unknown token '%s'\n", char);
					os.Exit(1);
				}

				literal := char;
				for {
					b, _ := reader.Peek(1);
					c := string(b);
					_, err := strconv.Atoi(c);

					if err != nil {
						tokens = append(tokens, createToken(literal, NUM));
						createdNumericToken = true;
						break;
					} else {
						bb, _ := reader.ReadByte();
						cc := string(bb);

						literal += cc;
					}
				}

				if !createdNumericToken {
					fmt.Printf("Err: unknown token '%s'\n", char);
					os.Exit(1);
				}
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
	validValues := []TokenType{STRING, NUM, BOOL, NULL}

	for i := 0; i < len(tokens); i ++ {
		token := tokens[i];

		if token.token == LBRACE {
			// not handling nested objects yet

			// left brace only appears as last token
			peekedValue := s.PeekBack(1);

			if !slices.Contains(validValues, peekedValue.token) {
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
		} else if (token.token == BOOL || token.token == NULL || token.token == NUM) {
			peekedValue := s.PeekBack(1);

			if peekedValue.token != COLON {
				valid = false;
				break;
			} else {
				s.Push(token);
			}
		}  else {
			s.Push(token);
		}
	}

	return valid;
}

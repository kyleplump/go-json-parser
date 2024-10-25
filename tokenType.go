package main

type TokenType int

const (
	RBRACE TokenType = iota
	LBRACE
	STRING
	COLON
	COMMA
)

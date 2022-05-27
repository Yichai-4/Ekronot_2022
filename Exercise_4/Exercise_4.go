/*
Submitter 1: Yichaï Hazan
ID: 1669535
Submitter 2: Amitai Shmeeda
ID: 305361479
Group number: 150060.01.5782.41-42
*/

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Collection of keywords in the Jack language
var keywords = []string{"class", "constructor", "function", "method", "field", "static", "var",
	"int", "char", "boolean", "void", "true", "false", "null", "this", "let", "do", "if", "else", "while", "return"}

// Collection of symbols in the Jack language
var symbols = []string{"{", "}", "(", ")", "[", "]", ".", ",", ";", "+", "-", "*", "/", "&", "|", "<", ">", "=", "~"}

// Collection of binary operators in the Jack language
var operators = []string{"+", "-", "*", "/", "&", "|", "<", ">", "="}

// Collection of unary operators in the Jack language
var unaryOperators = []string{"-", "~"}

// Collection of non-terminals language elements in the Jack language
/*var nonTerminals = []string{"class", "classVarDec", "subroutineDec", "parameterList", "subroutineBody", "varDec",
"statements", "whileStatement", "ifStatement", "returnStatement", "letStatement", "doStatement",
"expression", "term", "expressionList"}*/

// For the different indentations
var indentation string

// Receive the path in program argument
var path = os.Args[1] // receiving the path as cli parameter
var pathArray = strings.Split(path, "\\")

// Create the output file with the according name
var directoryName = pathArray[len(pathArray)-1]
var fileName string

// Initializes a file which will contain the tokens of the program
var tokensFile *os.File

var parsedFile *os.File

func main() {
	// Close the file "outputFile" at the end of the main function
	defer tokensFile.Close()

	// Go through the file and performs some operations
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}
		fullFileName := info.Name()
		extension := filepath.Ext(fullFileName)
		if extension == ".jack" {
			fmt.Printf("File Name: %s\n", fullFileName)
			fileName = strings.TrimRight(fullFileName, extension) // removes the extension from the file name

			inputJackFile, _ := os.Open(path)
			defer inputJackFile.Close()

			Tokenize(path)

			Parse(tokensFile)
		}
		return nil
	})
}

func Parse(tokensFile *os.File) {
	parsedFile, _ = os.Create(directoryName + "_" + fileName + ".xml")
	tokensFile, _ = os.Open(directoryName + "_" + fileName + "T.xml")

	data := bufio.NewScanner(tokensFile)
	data.Scan()
	data.Scan() // skip the first line <token>
	// Program structure
	CompileClass(data)
}

func CompileClass(data *bufio.Scanner) {
	parsedFile.WriteString("<class>\n")

	indentation = "  "
	Eat(data, "class")
	CompileIdentifier(data) // checks class name
	Eat(data, "{")
	// Implements classVarDec*
	words := strings.Split(data.Text(), " ")
	nextToken := words[1]
	for nextToken == "static" || nextToken == "field" {
		CompileClassVarDec(data)
		words = strings.Split(data.Text(), " ")
		nextToken = words[1]
	}
	// Implements subroutineDec*
	words = strings.Split(data.Text(), " ")
	nextToken = words[1]
	for nextToken == "constructor" || nextToken == "function" || nextToken == "method" {
		CompileSubroutineDec(data)
		words = strings.Split(data.Text(), " ")
		nextToken = words[1]
	}
	Eat(data, "}")

	parsedFile.WriteString("</class>\n")
}

func CompileSubroutineDec(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<subroutineDec>\n")
	indentation += "  "

	var firstWords = []string{"constructor", "function", "method"}
	EatOptions(data, firstWords)
	var voidOrType = []string{"int", "char", "boolean", "void"}
	EatOptions(data, voidOrType)
	CompileIdentifier(data) // checks subroutine name
	Eat(data, "(")
	CompileParameterList(data)
	Eat(data, ")")
	CompileSubroutineBody(data)

	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</subroutineDec>\n")
}

func CompileSubroutineBody(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<subroutineBody>\n")
	indentation += "  "

	Eat(data, "{")
	// Implements varDec*
	words := strings.Split(data.Text(), " ")
	nextToken := words[1]
	for nextToken == "var" {
		CompileVarDec(data)
		words = strings.Split(data.Text(), " ")
		nextToken = words[1]
	}
	CompileStatements(data)
	Eat(data, "}")

	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</subroutineBody>\n")
}

func CompileStatements(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<statements>\n")
	indentation += "  "

	var firstWords = []string{"let", "if", "while", "do", "return"}
	words := strings.Split(data.Text(), " ")
	nextToken := words[1]
	for StringInList(nextToken, firstWords) {
		switch nextToken {
		// Statements
		case "let":
			CompileLetStatement(data)
		case "if":
			CompileIfStatement(data)
		case "while":
			CompileWhileStatement(data)
		case "do":
			CompileDoStatement(data)
		case "return":
			CompileReturnStatement(data)
		}
		words = strings.Split(data.Text(), " ")
		nextToken = words[1]
	}

	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</statements>\n")
}

func CompileReturnStatement(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<returnStatement>\n")
	indentation += "  "

	Eat(data, "return")
	// Implements "expression?"
	if IsTerm(data) {
		CompileExpression(data)
	}
	Eat(data, ";")

	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</returnStatement>\n")
}

func IsTerm(data *bufio.Scanner) bool {
	words := strings.Split(data.Text(), " ")
	tokenType := words[0]
	currentToken := words[1]
	if IsInteger(currentToken) || tokenType == "<stringConstant>" || StringInList(currentToken, keywords) ||
		tokenType == "<identifier" || currentToken == "(" || StringInList(currentToken, unaryOperators) {
		return true
	}
	return false
}

func CompileDoStatement(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<doStatement>\n")
	indentation += "  "

	Eat(data, "do") // code to handle 'do'
	// Implements subroutine call
	CompileIdentifier(data)
	words := strings.Split(data.Text(), " ")
	currentToken := words[1]
	if currentToken == "(" {
		Eat(data, "(")
		CompileExpressionList(data)
		Eat(data, ")")
	} else {
		Eat(data, ".")
		CompileIdentifier(data)
		Eat(data, "(")
		CompileExpressionList(data)
		Eat(data, ")")
	}
	Eat(data, ";")

	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</doStatement>\n")
}

func CompileExpressionList(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<expressionList>\n")
	indentation += "  "

	words := strings.Split(data.Text(), " ")
	if IsTerm(data) {
		CompileExpression(data)
		words = strings.Split(data.Text(), " ")
		nextToken := words[1]
		for nextToken == "," {
			Eat(data, ",")
			CompileExpression(data)
		}
	}

	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</expressionList>\n")
}

func CompileWhileStatement(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<whileStatement>\n")
	indentation += "  "

	Eat(data, "while") // code to handle 'while'
	CompileExpression(data)
	Eat(data, ")")
	Eat(data, "{")
	CompileStatements(data)
	Eat(data, "}")

	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</whileStatement>\n")
}

func CompileIfStatement(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<ifStatement>\n")
	indentation += "  "

	Eat(data, "if") // code to handle 'if'
	Eat(data, "(")
	CompileExpression(data)
	Eat(data, ")")
	Eat(data, "{")
	CompileStatements(data)
	Eat(data, "}")
	// Implements ('else' '{'statements'}')?
	words := strings.Split(data.Text(), " ")
	currentToken := words[1]
	if currentToken == "else" {
		parsedFile.WriteString(indentation + data.Text() + "\n")
		data.Scan()
		Eat(data, "{")
		CompileStatements(data)
		Eat(data, "}")
	}

	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</ifStatement>\n")
}

func CompileLetStatement(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<letStatement>\n")
	indentation += "  "

	Eat(data, "let")        // code to handle 'let'
	CompileIdentifier(data) // check var name
	// Implements ('[' expression ']')?
	words := strings.Split(data.Text(), " ")
	currentToken := words[1]
	if currentToken == "[" {
		Eat(data, "[")
		CompileExpression(data)
		Eat(data, "]")
	}
	Eat(data, "=")
	CompileExpression(data)
	Eat(data, ";")

	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</letStatement>\n")
}

func CompileExpression(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<expression>\n")
	indentation += "  "

	CompileTerm(data)
	words := strings.Split(data.Text(), " ")
	nextToken := words[1]
	for StringInList(nextToken, operators) {
		Eat(data, nextToken)
		CompileTerm(data)
	}

	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</expression>\n")
}

func CompileTerm(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<term>\n")
	indentation += "  "

	CompileIdentifier(data)

	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</term>\n")
}

func CompileVarDec(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<varDec>\n")
	indentation += "  "

	Eat(data, "var")
	CompileType(data)
	CompileIdentifier(data) // checks variable name
	// Implements (',' type varName)*
	words := strings.Split(data.Text(), " ")
	nextToken := words[1]
	for nextToken == "," {
		Eat(data, ",")
		CompileIdentifier(data) // checks var name
		words = strings.Split(data.Text(), " ")
		nextToken = words[1]
	}
	Eat(data, ";")

	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</varDec>\n")
}

func CompileIdentifier(data *bufio.Scanner) {
	words := strings.Split(data.Text(), " ")
	tokenType := words[0]
	if tokenType != "<identifier>" {
		println("error - expected identifier")
		data.Scan()
	} else {
		parsedFile.WriteString(indentation + data.Text() + "\n")
		data.Scan()
	}
}

func CompileParameterList(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<parameterList>\n")
	indentation += "  "

	// Implements ((type varName) (',' type varName)*)?
	words := strings.Split(data.Text(), " ")
	tokenType := words[0]
	nextToken := words[1]
	var typeList = []string{"int", "char", "boolean"}
	if StringInList(nextToken, typeList) || tokenType == "<identifier" {
		CompileType(data)
		CompileIdentifier(data)
		words = strings.Split(data.Text(), " ")
		nextToken = words[1]
		for nextToken == "," {
			CompileType(data)
			CompileIdentifier(data) // checks var name
			words = strings.Split(data.Text(), " ")
			nextToken = words[1]
		}
	}
	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</parameterList>\n")
}

func CompileClassVarDec(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<classVarDec>\n")
	indentation += "  "

	var firstWords = []string{"static", "field"}
	EatOptions(data, firstWords)
	CompileType(data)
	CompileIdentifier(data) // checks variable name
	// Implements (',' varName)*
	words := strings.Split(data.Text(), " ")
	nextToken := words[1]
	for nextToken == "," {
		CompileIdentifier(data)
		words = strings.Split(data.Text(), " ")
		nextToken = words[1]
	}
	Eat(data, ";")

	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</classVarDec>\n")
}

func CompileType(data *bufio.Scanner) {
	var types = []string{"int", "char", "boolean"}
	words := strings.Split(data.Text(), " ")
	tokenType := words[0]
	currentToken := words[1]
	if !StringInList(currentToken, types) && tokenType != "<identifier>" {
		print("error - expected one of them: ")
		fmt.Print(strings.Join(types, ", "))
		println(", <identifier>")
	} else {
		parsedFile.WriteString(indentation + data.Text() + "\n")
		data.Scan()
	}
}

func EatOptions(data *bufio.Scanner, start []string) {
	words := strings.Split(data.Text(), " ")
	currentToken := words[1]
	if !StringInList(currentToken, start) {
		print("error - expected one of them: ")
		fmt.Println(strings.Join(start, ", "))
	} else {
		parsedFile.WriteString(indentation + data.Text() + "\n")
		data.Scan()
	}
}

func Eat(data *bufio.Scanner, s string) {
	words := strings.Split(data.Text(), " ")
	currentToken := words[1]
	if currentToken != s {
		println("error - expected " + s)
	} else {
		parsedFile.WriteString(indentation + data.Text() + "\n")
		data.Scan()
	}
}

// Tokenize Reads the input file, goes through it character by character,
// and writes all the tokens it contains (according to the JACK language syntax) in the output file
func Tokenize(jackFilePath string) {
	tokensFile, _ = os.Create(directoryName + "_" + fileName + "T.xml")
	tokensFile.WriteString("<tokens>\n") // the xml file (xxxT.xml) has to begin by "<tokens>"

	// Reads the input file and converts it into string of characters
	fileBuffer, err := ioutil.ReadFile(jackFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	inputData := string(fileBuffer)
	data := bufio.NewScanner(strings.NewReader(inputData))
	data.Split(bufio.ScanRunes)

	var tokenClassification string

	for data.Scan() {
		char := data.Text()
		var token string
		var nextChar string
		switch {
		// Skip all comment types
		case char == "/": // "//" or "/*" or "/**"
			SkipCommentLines(data, token, char)

		// Handles keyword or identifier
		case char == "_", IsLetter(char):
			token, nextChar = GetKeywordIdentifierToken(data, token, char)
			if StringInList(token, keywords) {
				tokenClassification = "keyword"
			} else {
				tokenClassification = "identifier"
			}
			WriteToken(tokenClassification, token)
			// Case of a symbol next to keyword or identifier
			if nextChar != " " && StringInList(nextChar, symbols) {
				WriteToken("symbol", nextChar)
			}

		// Handles symbol
		case StringInList(char, symbols):
			tokenClassification = "symbol"
			token = GetSymbolToken(char, token)
			WriteToken(tokenClassification, token)

		// Handles integer constant
		case IsInteger(char):
			tokenClassification = "integerConstant"
			token, nextChar = GetIntegerToken(data, token, char)
			WriteToken(tokenClassification, token)
			// Case of a symbol next to integer
			if nextChar != " " && StringInList(nextChar, symbols) {
				WriteToken("symbol", nextChar)
			}

		// Handles string constant
		case char == "\"":
			tokenClassification = "stringConstant"
			token = GetStringToken(data, token)
			WriteToken(tokenClassification, token)

		default: // skips spaces
			break
		}
	}
	tokensFile.WriteString("</tokens>\n") // the xml file (xxxT.xml) has to end by "</tokens>"

	if err := data.Err(); err != nil {
		log.Fatal(err)
	}
}

func SkipCommentLines(data *bufio.Scanner, token string, char string) {
	data.Scan()
	nextChar := data.Text()
	switch nextChar {
	// Start comment
	case "/": // found "//"
		data.Scan()
		nextChar = data.Text()
		for nextChar != "\n" {
			data.Scan()
			nextChar = data.Text()
		}
	case "*": // found "/*"
		data.Scan()
		nextChar = data.Text()
		if nextChar == "*" { // found "/**"
			data.Scan()
			nextChar = data.Text()
		}
		for nextChar != "*" {
			data.Scan()
			nextChar = data.Text()
			if nextChar == "*" {
				data.Scan()
				nextChar = data.Text()
				if nextChar == "/" {
					break
				}
			}
		}
	default: // found single symbol "/"
		tokenClassification := "symbol"
		token = char
		WriteToken(tokenClassification, char)
	}
}

// GetKeywordIdentifierToken Returns the keyword or identifier of the current token and the next char
func GetKeywordIdentifierToken(data *bufio.Scanner, token string, char string) (string, string) {
	token += char
	data.Scan()
	nextChar := data.Text()
	for nextChar == "_" || IsLetter(nextChar) || IsInteger(nextChar) {
		token += nextChar
		data.Scan()
		nextChar = data.Text()
	}
	return token, nextChar
}

// GetSymbolToken Returns the according symbol of the current token
func GetSymbolToken(char string, token string) string {
	switch char {
	// Special symbol
	case "<":
		token = "&lt;"
	case ">":
		token = "&gt;"
	case "\"":
		token = "&quot;"
	case "&":
		token = "&amp;"
	// Regular symbol
	default:
		token = char
	}
	return token
}

// GetIntegerToken Returns the integer value of the current token and the next char
func GetIntegerToken(data *bufio.Scanner, token string, char string) (string, string) {
	token = char
	data.Scan()
	nextChar := data.Text()
	for IsInteger(nextChar) {
		token += nextChar
		data.Scan()
		nextChar = data.Text()
	}
	return token, nextChar
}

// GetStringToken Returns the string value of the current token
func GetStringToken(data *bufio.Scanner, token string) string {
	data.Scan()
	nextChar := data.Text()
	for nextChar != "\"" {
		token += nextChar
		data.Scan()
		nextChar = data.Text()
	}
	return token
}

// WriteToken Write the current token on the file according to the correct format
func WriteToken(tokenClassification string, token string) {
	tokensFile.WriteString("<" + tokenClassification + "> ")
	tokensFile.WriteString(token)
	tokensFile.WriteString(" </" + tokenClassification + ">\n")
}

// StringInList Checks if the list contains the string a or not
func StringInList(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// IsInteger Checks if the string is an integer or not
func IsInteger(s string) bool {
	if _, errInt := strconv.Atoi(s); errInt == nil {
		return true
	}
	return false
}

// IsLetter Checks if the string is a letter or not
func IsLetter(s string) bool {

	for _, char := range s {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') {
			return false
		}
	}
	return true
}

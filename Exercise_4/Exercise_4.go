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

// Collection of keywords in the JACK language
var keyword = []string{"class", "constructor", "function", "method", "field", "static", "var",
	"int", "char", "boolean", "void", "true", "false", "null", "this", "let", "do", "if", "else", "while", "return"}

// Collection of symbols in the JACK language
var symbol = []string{"{", "}", "(", ")", "[", "]", ".", ",", ";", "+", "-", "*", "/", "&", "|", "<", ">", "=", "~"}

// Receive the path in program argument
var path = os.Args[1] // receiving the path as cli parameter
var pathArray = strings.Split(path, "\\")

// Create the output file with the according name
var directoryName = pathArray[len(pathArray)-1]

// Initializes a file which will contain the tokens of the program
var tokensFile *os.File

//var parserFile, _  = os.Create("my" + directoryName + ".xml")

func main() {
	// Close the file "outputFile" at the end of the main function
	defer tokensFile.Close()

	// Go through the file and performs some operations
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}
		fileName := info.Name()
		extension := filepath.Ext(fileName)
		if extension == ".jack" {
			fmt.Printf("File Name: %s\n", fileName)
			// removes the extension from the file name
			name := strings.TrimRight(fileName, extension)
			tokensFile, _ = os.Create(directoryName + "_" + name + "T.xml")

			inputJackFile, _ := os.Open(path)
			defer inputJackFile.Close()

			Tokenize(tokensFile, path)

		}
		return nil
	})
}

// Tokenize Reads the input file, goes through it character by character,
// and writes all the tokens it contains (according to the JACK language syntax) in the output file
func Tokenize(outputFile *os.File, inputFilePath string) {
	outputFile.WriteString("<tokens>\n") // the xml file (xxxT.xml) has to begin by "<tokens>"

	// Reads the input file and converts it into string of characters
	fileBuffer, err := ioutil.ReadFile(inputFilePath)
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
			if StringInList(token, keyword) {
				tokenClassification = "keyword"
			} else {
				tokenClassification = "identifier"
			}
			WriteToken(tokenClassification, token)
			// Case of a symbol next to keyword or identifier
			if nextChar != " " && StringInList(nextChar, symbol) {
				WriteToken("symbol", nextChar)
			}

		// Handles symbol
		case StringInList(char, symbol):
			tokenClassification = "symbol"
			token = GetSymbolToken(char, token)
			WriteToken(tokenClassification, token)

		// Handles integer constant
		case IsInteger(char):
			tokenClassification = "integerConstant"
			token, nextChar = GetIntegerToken(data, token, char)
			WriteToken(tokenClassification, token)
			// Case of a symbol next to integer
			if nextChar != " " && StringInList(nextChar, symbol) {
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
	outputFile.WriteString("</tokens>\n") // the xml file (xxxT.xml) has to end by "</tokens>"

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

// GetKeywordIdentifierToken Gets the current keyword or identifier token and the next char
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

// GetSymbolToken Gets the current correct symbol token
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

// GetIntegerToken Gets the token of the current constant integer and the next char
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

// GetStringToken Gets the token of the current constant string
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

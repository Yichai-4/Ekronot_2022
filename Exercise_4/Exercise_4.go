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

// Constants
var keyword = []string{"class", "constructor", "function", "method", "field", "static", "var",
	"int", "char", "boolean", "void", "true", "false", "null", "this", "let", "do", "if", "else", "while", "return"}

var symbol = []string{"{", "}", "(", ")", "[", "]", ".", ",", ";", "+", "-", "*", "/", "&", "|", "<", ">", "=", "~"}

// Receive the path in program argument
var path = os.Args[1] // receiving the path as cli parameter
var pathArray = strings.Split(path, "\\")

// Create the output file with the according name
var directoryName = pathArray[len(pathArray)-1]

//var tokensFile, _ = os.Create(directoryName + "T.xml")
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

func Tokenize(outputFile *os.File, inputFilePath string) {
	outputFile.WriteString("<tokens>\n")

	//data := bufio.NewScanner(inputFilePath)
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
		_, errInt := strconv.Atoi(char)
		var token = ""
		var nextChar string
	switchAgain:
		switch {
		case char == "/": // "//" or "/*" or "/**"
			data.Scan()
			nextChar = data.Text()
			switch nextChar {
			case "/": // found "//"
				data.Scan()
				nextChar = data.Text()
				for nextChar != "\n" {
					data.Scan()
					nextChar = data.Text()
				}
			case "*": // found "/*" or "/**"
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
			default:
				tokenClassification = "symbol"
				token = char
				WriteToken(tokenClassification, char)
			}

		case char == "_", IsLetter(char): // keyword or identifier
			//KeywordFunc()
			token += char
			data.Scan()
			nextChar = data.Text()
			for nextChar == "_" || IsLetter(nextChar) || errInt == nil {
				token += nextChar
				data.Scan()
				nextChar = data.Text()
			}
			if stringInList(token, keyword) {
				tokenClassification = "keyword"
			} else {
				tokenClassification = "identifier"
			}
			WriteToken(tokenClassification, token)
			//println(token + " " + nextChar)
			if nextChar != " " && stringInList(nextChar, symbol) {
				WriteToken("symbol", nextChar)
			}
			break switchAgain

		case stringInList(char, symbol): // symbol
			//SymbolFunc()
			tokenClassification = "symbol"
			switch char { // Special characters
			case "<":
				token = "&lt;"
			case ">":
				token = "&gt;"
			case "\"":
				token = "&quot;"
			case "&":
				token = "&amp;"
			default:
				token = char
			}
			WriteToken(tokenClassification, token)

		case errInt == nil: // integer constant
			//IntegerConstantFunc()
			tokenClassification = "integerConstant"
			nextChar = char
			token += nextChar
			data.Scan()
			nextChar = data.Text()
			_, errInt = strconv.Atoi(nextChar)
			for errInt == nil {
				token += nextChar
				data.Scan()
				nextChar = data.Text()
				_, errInt = strconv.Atoi(nextChar)
			}
			WriteToken(tokenClassification, token)
			if nextChar != " " && stringInList(nextChar, symbol) {
				WriteToken("symbol", nextChar)
			}
			break switchAgain

		case char == "\"": // string constant
			//IdentifierFunc()
			tokenClassification = "stringConstant"
			data.Scan()
			char = data.Text()
			for char != "\"" {
				token += data.Text()
				data.Scan()
				char = data.Text()
			}
			WriteToken(tokenClassification, token)

		default: // skips spaces
			break
		}
	}
	outputFile.WriteString("</tokens>\n")

	if err := data.Err(); err != nil {
		log.Fatal(err)
	}
}

func WriteToken(tokenClassification string, token string) {
	tokensFile.WriteString("<" + tokenClassification + "> ")
	tokensFile.WriteString(token)
	tokensFile.WriteString(" </" + tokenClassification + ">\n")
}

func IntegerConstantFunc() {

}

func IdentifierFunc() {

}

func SymbolFunc() {

}

func KeywordFunc() {

}

func stringInList(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func IsLetter(s string) bool {
	for _, char := range s {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') {
			return false
		}
	}
	return true
}

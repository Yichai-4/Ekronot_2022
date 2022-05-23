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
var tokensFile, _ = os.Create("my" + directoryName + "T.xml")

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
			// removes the extension from the file name and prints it
			name := strings.TrimRight(fileName, extension)
			tokensFile.WriteString("//" + name + "\n")
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

	var skipLine bool

	for data.Scan() {
		char := data.Text()
		_, errInt := strconv.Atoi(char)
		var token = ""
		switch {
		case char == "/": // "//" or "/*" or "/**"
			data.Scan()
			nextChar := data.Text()
			switch nextChar {
			case "/": // found "//"
				skipLine = true
				data.Scan()
				nextChar = data.Text()
				for nextChar != "\n" {
					data.Scan()
					nextChar = data.Text()
				}
			case "*": // found "/*" or "/**"
				skipLine = true
				data.Scan()
				nextChar = data.Text()
				for nextChar != "*" {
					data.Scan()
					nextChar = data.Text()
					if nextChar == "*" {
						data.Scan()
						if data.Text() == "/" {
							break
						}
					}
				}
			default:
				tokenClassification = "symbol"
				token = char
			}

		//SkipCommentLines(*outputFile)
		case char == "_", IsLetter(char): // keyword or identifier
			skipLine = false
			//KeywordFunc()
			token += char
			data.Scan()
			nextChar := data.Text()
			for nextChar == "_" || IsLetter(nextChar) || errInt == nil {
				token += nextChar
				data.Scan()
				nextChar = data.Text()
			}
			//char = nextChar
			if stringInList(token, keyword) {
				tokenClassification = "keyword"
			} else {
				tokenClassification = "identifier"
			}

		case stringInList(char, symbol): // symbol
			skipLine = false
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

		case errInt == nil: // integer constant
			skipLine = false
			//IntegerConstantFunc()
			tokenClassification = "integerConstant"
			token += char
			data.Scan()
			_, errInt = strconv.Atoi(data.Text())
			for errInt == nil {
				char = data.Text()
				token += char
				data.Scan()
				_, errInt = strconv.Atoi(data.Text())
			}

		case char == "\"": // string constant
			skipLine = false
			//IdentifierFunc()
			tokenClassification = "stringConstant"
			data.Scan()
			char = data.Text()
			for char != "\"" {
				token += data.Text()
				data.Scan()
				char = data.Text()
			}
		default: // skips spaces
			continue
		}
		if !skipLine {
			outputFile.WriteString("<" + tokenClassification + "> ")
			outputFile.WriteString(token)
			outputFile.WriteString(" </" + tokenClassification + ">\n")
		}
	}
	outputFile.WriteString("</tokens>")

	if err := data.Err(); err != nil {
		log.Fatal(err)
	}
}

func SkipCommentLines(file os.File) {

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

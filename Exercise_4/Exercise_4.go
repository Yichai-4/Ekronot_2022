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
	"log"
	"os"
	"path/filepath"
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

			Tokenize(tokensFile, inputJackFile)

		}
		return nil
	})
}

func Tokenize(outputFile *os.File, inputFilePath *os.File) {
	outputFile.WriteString("<tokens>\n")

	data := bufio.NewScanner(inputFilePath)
	var tokenClassification string

	for data.Scan() {
		/*
			lines := strings.Split(data.Text(), "\n")
			for _, line := range lines {
				if line == "" {
				}
			}*/
		words := strings.Split(data.Text(), " ")
		firstWord := words[0]
		if firstWord == "//" || firstWord == "/*" || firstWord == "/**" {
			continue
		}
		for _, word := range words {
			currentToken := word
			if stringInList(word, keyword) {
				tokenClassification = "keyword"
			}
			if stringInList(word, symbol) {
				tokenClassification = "symbol"
			}
			outputFile.WriteString("<" + tokenClassification + "> ")
			outputFile.WriteString(currentToken)
			outputFile.WriteString(" </" + tokenClassification + ">\n")
		}

	}
	outputFile.WriteString("</tokens>")

	if err := data.Err(); err != nil {
		log.Fatal(err)
	}
}

func stringInList(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

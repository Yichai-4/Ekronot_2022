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

//var keywords = {"class"}

// Receive the path in program argument
var path = os.Args[1] // receiving the path as cli parameter
var pathArray = strings.Split(path, "\\")

// Create the output file with the according name
var directoryName = pathArray[len(pathArray)-1]
var outputFile, _ = os.Create(directoryName + ".xml")

func main() {
	// Close the file "outputFile" at the end of the main function
	defer outputFile.Close()

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
			outputFile.WriteString(name + "\n")

			inputFile, err := os.Open(path)
			check(err)
			defer inputFile.Close()

			scanner := bufio.NewScanner(inputFile)

			for scanner.Scan() {
				outputFile.WriteString("<test> ok </test>\n")
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}

		}
		return nil
	})
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

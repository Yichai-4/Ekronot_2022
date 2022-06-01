/*
Submitter 1: Yichaï Hazan
ID: 1669535
Submitter 2: Amitai Shmeeda
ID: 305361479
Group number: 150060.01.5782.41-42
*/

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Receive the path in program argument
var path = os.Args[1] // receiving the path as cli parameter
var pathArray = strings.Split(path, "\\")

// Create the output file with the according name
var fileName = pathArray[len(pathArray)-1]
var vmFile, _ = os.Create(fileName + ".vm")

func main() {
	// Close the file "vmFile" at the end of the main function
	defer vmFile.Close()

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
			vmFile.WriteString("// Program: " + fileName + ".jack\n")

			inputJackFile, _ := os.Open(path)
			defer inputJackFile.Close()

		}
		return nil
	})
}

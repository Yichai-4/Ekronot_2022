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

// Defining constants

// var PushPopCommand = [2]string{"push", "pop"}

// Receive the path in program argument
//var path = os.Args[1] // receiving the path as cli parameter
var path = "C:\\Ekronot_2022\\nand2tetris\\projects\\07\\StackArithmetic\\SimpleAdd"
var pathArray = strings.Split(path, "\\")

// Create the output file with the according name
var fileName = pathArray[len(pathArray)-1]
var outputFile, _ = os.Create(fileName + ".asm")

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
		if extension == ".vm" {
			fmt.Printf("File Name: %s\n", fileName)
			// removes the extension from the file name and prints it
			name := strings.TrimRight(fileName, extension)
			outputFile.WriteString(name + ":\n")

			inputFile, err := os.Open(path)
			check(err)
			defer inputFile.Close()

			scanner := bufio.NewScanner(inputFile)

			for scanner.Scan() {
				words := strings.Split(scanner.Text(), " ")
				command := words[0]
				if command == "push" || command == "pop" {
					segment := words[1]
					i := words[2]
					if command == "push" {
						PushTranslation(segment, i)
					}
					if command == "pop" {
						PopTranslation(segment, i)
					}
				}

			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}

		}
		return nil
	})

}

func PopTranslation(segment string, i string) {

}

func PushTranslation(segment string, i string) {
	// Translation for the command push constant i
	if segment == "constant" {
		outputFile.WriteString("@" + i + "\nD=A\n") // D=i
		outputFile.WriteString("@SP\nA=M\nM=D\n")   // *SP=D
		outputFile.WriteString("@SP\nM=M+1\n")      // SP++
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

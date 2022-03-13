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
			outputFile.WriteString("// " + name + ":\n")

			inputFile, err := os.Open(path)
			check(err)
			defer inputFile.Close()

			scanner := bufio.NewScanner(inputFile)

			for scanner.Scan() {
				words := strings.Split(scanner.Text(), " ")
				command := words[0]
				switch command {
				// Arithmetic commands
				case "add": // Integer addition (2's complement)
					AddTranslation()
				case "sub": // Integer subtraction (2's complement)
					SubTranslation()
				case "neg": // Arithmetic negation (2's complement)
					NegTranslation()
				// Boolean commands
				case "eq": // Equality
					EqTranslation()
				case "gt": // Greater than
					GtTranslation()
				case "lt": // Less than
					LtTranslation()
				case "and":
					AndTranslation()
				case "or":
					OrTranslation()
				case "not":
					NotTranslation()
				// Memory access commands
				case "push":
					segment := words[1]
					i := words[2]
					PushTranslation(segment, i)
				case "pop":
					segment := words[1]
					i := words[2]
					PopTranslation(segment, i)
				}

			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}

		}
		return nil
	})

}

func AddTranslation() {

}

func SubTranslation() {

}

func NegTranslation() {

}

func EqTranslation() {

}

func GtTranslation() {

}

func LtTranslation() {

}

func AndTranslation() {

}

func OrTranslation() {

}

func NotTranslation() {

}

// PopTranslation Translation of pop command to hack language
func PopTranslation(segment string, i string) {
	outputFile.WriteString("// pop " + segment + i + "\n") // general comment for the respective pop command
	switch segment {
	// Translation for the command pop local i
	case "local":
		outputFile.WriteString("@" + i + "\nD=A\n@LCL\nD=D+M\n@addr\nM=D\n") // addr=LCL+i
		outputFile.WriteString("@SP\nM=M-1\n")                               // SP--
		outputFile.WriteString("@SP\nA=M\nD=M\n@addr\nA=M\nM=D\n")           // *addr=*SP
	// Translation for the command pop argument i
	case "argument":
		outputFile.WriteString("@" + i + "\nD=A\n@ARG\nD=D+M\n@addr\nM=D\n") // addr=ARG+i
		outputFile.WriteString("@SP\nM=M-1\n")                               // SP--
		outputFile.WriteString("@SP\nA=M\nD=M\n@addr\nA=M\nM=D\n")           // *addr=*SP
	// Translation for the command pop this i
	case "this":
		outputFile.WriteString("@" + i + "\nD=A\n@THIS\nD=D+M\n@addr\nM=D\n") // addr=THIS+i
		outputFile.WriteString("@SP\nM=M-1\n")                                // SP--
		outputFile.WriteString("@SP\nA=M\nD=M\n@addr\nA=M\nM=D\n")            // *addr=*SP
	// Translation for the command pop that i
	case "that":
		outputFile.WriteString("@" + i + "\nD=A\n@THAT\nD=D+M\n@addr\nM=D\n") // addr=THAT+i
		outputFile.WriteString("@SP\nM=M-1\n")                                // SP--
		outputFile.WriteString("@SP\nA=M\nD=M\n@addr\nA=M\nM=D\n")            // *addr=*SP
	// Translation for the command pop static i
	case "static":
		outputFile.WriteString("@SP\nM=M-1\n")                       // SP--
		outputFile.WriteString("@SP\nA=M\nD=M\n")                    // D=*SP
		outputFile.WriteString("@" + fileName + "." + i + "\nM=D\n") // static i = D
	// Translation for the command pop temp i
	case "temp":
		outputFile.WriteString("@" + i + "\nD=A\n@5\nD=D+A\n@addr\nM=D\n") // addr=5+i
		outputFile.WriteString("@SP\nM=M-1\n")                             // SP--
		outputFile.WriteString("@SP\nA=M\nD=M\n@addr\nA=M\nM=D\n")         // *addr=*SP
	case "pointer":
		outputFile.WriteString("@SP\nM=M-1\n") // SP--
		if i == "0" {
			outputFile.WriteString("@THIS\nD=M\n@addr\nM=D\n") // addr=THIS
		} else { // i == "1"
			outputFile.WriteString("@THAT\nD=M\n@addr\nM=D\n") // addr=THAT
		}
		outputFile.WriteString("@SP\nA=M\nD=M\n@addr\nA=M\nM=D\n") // *addr=*SP
	}
}

// PushTranslation Translation of push command to hack language
func PushTranslation(segment string, i string) {
	outputFile.WriteString("// push " + segment + i + "\n") // general comment for the respective push command
	switch segment {
	// Translation for the command push local i
	case "local":
		outputFile.WriteString("@" + i + "\nD=A\n@LCL\nD=D+M\n@addr\nM=D\n") // addr=LCL+i
		outputFile.WriteString("A=M\nD=M\n@SP\nA=M\nM=D\n")                  // *SP=*addr
		outputFile.WriteString("@SP\nM=M+1\n")                               // SP++
	// Translation for the command push argument i
	case "argument":
		outputFile.WriteString("@" + i + "\nD=A\n@ARG\nD=D+M\n@addr\nM=D\n") // addr=ARG+i
		outputFile.WriteString("A=M\nD=M\n@SP\nA=M\nM=D\n")                  // *SP=*addr
		outputFile.WriteString("@SP\nM=M+1\n")                               // SP++
	// Translation for the command push this i
	case "this":
		outputFile.WriteString("@" + i + "\nD=A\n@THIS\nD=D+M\n@addr\nM=D\n") // addr=THIS+i
		outputFile.WriteString("A=M\nD=M\n@SP\nA=M\nM=D\n")                   // *SP=*addr
		outputFile.WriteString("@SP\nM=M+1\n")                                // SP++
	// Translation for the command push that i
	case "that":
		outputFile.WriteString("@" + i + "\nD=A\n@THAT\nD=D+M\n@addr\nM=D\n") // addr=THAT+i
		outputFile.WriteString("A=M\nD=M\n@SP\nA=M\nM=D\n")                   // *SP=*addr
		outputFile.WriteString("@SP\nM=M+1\n")                                // SP++
	// Translation for the command push constant i
	case "constant":
		outputFile.WriteString("@" + i + "\nD=A\n") // D=i
		outputFile.WriteString("@SP\nA=M\nM=D\n")   // *SP=D
		outputFile.WriteString("@SP\nM=M+1\n")      // SP++
	// Translation for the command push static i
	case "static":
		outputFile.WriteString("@" + fileName + "." + i + "\nD=M\n") // D = static i
		outputFile.WriteString("@SP\nA=M\nM=D\n")                    // *SP=D
		outputFile.WriteString("@SP\nM=M+1\n")                       // SP++
	// Translation for the command push temp i
	case "temp":
		outputFile.WriteString("@" + i + "\nD=A\n@5\nD=D+A\n@addr\nM=D\n") // addr=5+i
		outputFile.WriteString("A=M\nD=M\n@SP\nA=M\nM=D\n")                // *SP=*addr
		outputFile.WriteString("@SP\nM=M+1\n")                             // SP++
	// Translation for the command push pointer 0/1
	case "pointer":
		if i == "0" {
			outputFile.WriteString("@THIS\nD=M\n") // D=*THIS
		} else { // i == "1"
			outputFile.WriteString("@THAT\nD=M\n") // D=*THAT
		}
		outputFile.WriteString("@SP\nA=M\nM=D\n") // *SP=D
		outputFile.WriteString("@SP\nM=M+1\n")    // SP++
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

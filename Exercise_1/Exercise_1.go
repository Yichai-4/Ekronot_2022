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

// Receive the path in program argument
// var path = os.Args[1] // receiving the path as cli parameter
var path = "C:\\Ekronot_2022\\nand2tetris\\projects\\07\\MemoryAccess\\BasicTest"
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
			outputFile.WriteString("// Program: " + name + ".asm\n")

			inputFile, err := os.Open(path)
			check(err)
			defer inputFile.Close()

			scanner := bufio.NewScanner(inputFile)

			for scanner.Scan() {
				words := strings.Split(scanner.Text(), " ")
				command := words[0]
				switch command {
				// Arithmetic commands
				case "add":
				case "sub":
				case "neg":
					WriteArithmetic(command)
				// Boolean commands
				case "eq": // Equality
				case "gt": // Greater than
				case "lt": // Less than
					WriteBoolean(command)
				// Logical commands - Bit-wise
				case "and":
				case "or":
				case "not":
					WriteLogical(command)
				// Memory access commands
				case "pop":
					segment := words[1]
					i := words[2]
					WritePop(segment, i)
				case "push":
					segment := words[1]
					i := words[2]
					WritePush(segment, i)
				}
			}
			// Ends the program with an infinite loop
			outputFile.WriteString("(END)\n@END\n0;JMP\n")

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})
}

// WriteArithmetic Translation of arithmetic command (i.e. add, sub and neg) in VM language to Hack language
func WriteArithmetic(command string) {
	switch command {
	case "add": // Integer addition (2's complement)
		outputFile.WriteString("// add\n")
		outputFile.WriteString("@SP\nM=M-1\nA=M\nD=M\n@result\nM=D\n")   // SP--, result=y(*SP)
		outputFile.WriteString("@SP\nM=M-1\nA=M\nD=M\n@result\nM=D+M\n") // SP--, result=x(*SP)+y
	case "sub": // Integer subtraction (2's complement)
		outputFile.WriteString("// sub\n")
		outputFile.WriteString("@SP\nM=M-1\nA=M\nD=M\n@result\nM=D\n")   // SP--, result=y(*SP)
		outputFile.WriteString("@SP\nM=M-1\nA=M\nD=M\n@result\nM=D-M\n") // SP--, result=x(*SP)-y
	case "neg": // Arithmetic negation (2's complement)
		outputFile.WriteString("// neg\n")
		outputFile.WriteString("@SP\nM=M-1\nA=M\nD=M\n@result\nM=M-D\n") // SP--, result=0-x
	}
	outputFile.WriteString("@result\nD=M\n@SP\nA=M\nM=D\n") // *SP=result
	outputFile.WriteString("@SP\nM=M+1\n")                  // SP++
}

// WriteBoolean Translation of boolean command (i.e. eq, gt or lt) in VM language to Hack language
func WriteBoolean(command string) {
	switch command {
	case "eq": // Equality
		outputFile.WriteString("// eq\n")
		outputFile.WriteString("@SP\nM=M-1\nA=M\nD=M\n") // SP--, D=STACK[SP]
		outputFile.WriteString("@SP\nA=M-1\nD=D-M\n")    // D=D-STACK[SP-1]
		outputFile.WriteString("@IF_TRUE\nD;JEQ\n")      // jump to (IF_TRUE) if D=0
	case "gt": // Greater than
		outputFile.WriteString("// gt\n")
		outputFile.WriteString("@SP\nM=M-1\nA=M\nD=M\n") // SP--, D=STACK[SP]
		outputFile.WriteString("@SP\nA=M-1\nD=D-M\n")    // D=D-STACK[SP-1]
		outputFile.WriteString("@IF_TRUE\nD;JGT\n")      // jump to (IF_TRUE) if D>0
	case "lt": // Less than
		outputFile.WriteString("// lt\n")
		outputFile.WriteString("@SP\nM=M-1\nA=M\nD=M\n") // SP--, D=STACK[SP]
		outputFile.WriteString("@SP\nA=M-1\nD=D-M\n")    // D=D-STACK[SP-1]
		outputFile.WriteString("@IF_TRUE\nD;JLT\n")      // jump to (IF_TRUE) if D<0
	}
	// If the condition is not met
	outputFile.WriteString("@SP\nA=M-1\nM=0\n") // *SP=0
	outputFile.WriteString("@IF_FALSE\n")       // LABEL
	outputFile.WriteString("0 ; JMP\n")         // unconditional jump
	// Otherwise
	outputFile.WriteString("(IF_TRUE)\n")        // Declaring a label
	outputFile.WriteString("@SP\nA=M-1\nM=-1\n") // *SP=-1
	outputFile.WriteString("(IF_FALSE)\n")       // Declaring a label
}

// WriteLogical Translation of logical command (i.e. and, or and not) in VM language to Hack language
func WriteLogical(command string) {
	switch command {
	case "and":
		outputFile.WriteString("// and\n")
		outputFile.WriteString("@SP\nM=M-1\nA=M\nD=M\n") // D=STACK[SP]
		outputFile.WriteString("A=A-1\nM=D&M\n")         // STACK[SP]=x and y
	case "or":
		outputFile.WriteString("// or\n")
		outputFile.WriteString("@SP\nM=M-1\nA=M\nD=M\n") // D=STACK[SP]
		outputFile.WriteString("A=A-1\nM=D||M\n")        // STACK[SP]=x or y
	case "not":
		outputFile.WriteString("// not\n")
		outputFile.WriteString("@SP\nA=M-1\nM=!M\n") // STACK[SP]=not(x)
	}
}

// WritePop Translation of pop command (in VM language) to Hack language
func WritePop(segment string, i string) {
	outputFile.WriteString("// pop " + segment + " " + i + "\n") // general comment for the respective pop command
	switch segment {
	// Translation for the command pop local i
	case "local":
		outputFile.WriteString("@" + i + "\nD=A\n@LCL\nD=D+M\n@addr\nM=D\n") // addr=LCL+i
		outputFile.WriteString("@SP\nM=M-1\n")                               // SP--
	// Translation for the command pop argument i
	case "argument":
		outputFile.WriteString("@" + i + "\nD=A\n@ARG\nD=D+M\n@addr\nM=D\n") // addr=ARG+i
		outputFile.WriteString("@SP\nM=M-1\n")                               // SP--
	// Translation for the command pop this i
	case "this":
		outputFile.WriteString("@" + i + "\nD=A\n@THIS\nD=D+M\n@addr\nM=D\n") // addr=THIS+i
		outputFile.WriteString("@SP\nM=M-1\n")                                // SP--
	// Translation for the command pop that i
	case "that":
		outputFile.WriteString("@" + i + "\nD=A\n@THAT\nD=D+M\n@addr\nM=D\n") // addr=THAT+i
		outputFile.WriteString("@SP\nM=M-1\n")                                // SP--
	// Translation for the command pop static i
	case "static":
		outputFile.WriteString("@SP\nM=M-1\n")                       // SP--
		outputFile.WriteString("@SP\nA=M\nD=M\n")                    // D=*SP
		outputFile.WriteString("@" + fileName + "." + i + "\nM=D\n") // static i = D
		os.Exit(0)
	// Translation for the command pop temp i
	case "temp":
		outputFile.WriteString("@" + i + "\nD=A\n@5\nD=D+A\n@addr\nM=D\n") // addr=5+i
		outputFile.WriteString("@SP\nM=M-1\n")                             // SP--
	// Translation for the command pop pointer 0/1
	case "pointer":
		outputFile.WriteString("@SP\nM=M-1\n") // SP--
		if i == "0" {
			outputFile.WriteString("@THIS\nD=M\n@addr\nM=D\n") // addr=THIS
		} else { // i == "1"
			outputFile.WriteString("@THAT\nD=M\n@addr\nM=D\n") // addr=THAT
		}
	}
	// For all pop commands (except for static) add the value to the according address:
	outputFile.WriteString("@SP\nA=M\nD=M\n@addr\nA=M\nM=D\n") // *addr=*SP
}

// WritePush Translation of push command (in VM language) to Hack language
func WritePush(segment string, i string) {
	outputFile.WriteString("// push " + segment + " " + i + "\n") // general comment for the respective push command
	switch segment {
	// Translation for the command push local i
	case "local":
		outputFile.WriteString("@" + i + "\nD=A\n@LCL\nD=D+M\n@addr\nM=D\n") // addr=LCL+i
		outputFile.WriteString("A=M\nD=M\n@SP\nA=M\nM=D\n")                  // *SP=*addr
	// Translation for the command push argument i
	case "argument":
		outputFile.WriteString("@" + i + "\nD=A\n@ARG\nD=D+M\n@addr\nM=D\n") // addr=ARG+i
		outputFile.WriteString("A=M\nD=M\n@SP\nA=M\nM=D\n")                  // *SP=*addr
	// Translation for the command push this i
	case "this":
		outputFile.WriteString("@" + i + "\nD=A\n@THIS\nD=D+M\n@addr\nM=D\n") // addr=THIS+i
		outputFile.WriteString("A=M\nD=M\n@SP\nA=M\nM=D\n")                   // *SP=*addr
	// Translation for the command push that i
	case "that":
		outputFile.WriteString("@" + i + "\nD=A\n@THAT\nD=D+M\n@addr\nM=D\n") // addr=THAT+i
		outputFile.WriteString("A=M\nD=M\n@SP\nA=M\nM=D\n")                   // *SP=*addr
	// Translation for the command push constant i
	case "constant":
		outputFile.WriteString("@" + i + "\nD=A\n") // D=i
		outputFile.WriteString("@SP\nA=M\nM=D\n")   // *SP=D
	// Translation for the command push static i
	case "static":
		outputFile.WriteString("@" + fileName + "." + i + "\nD=M\n") // D = static i
		outputFile.WriteString("@SP\nA=M\nM=D\n")                    // *SP=D
	// Translation for the command push temp i
	case "temp":
		outputFile.WriteString("@" + i + "\nD=A\n@5\nD=D+A\n@addr\nM=D\n") // addr=5+i
		outputFile.WriteString("A=M\nD=M\n@SP\nA=M\nM=D\n")                // *SP=*addr
	// Translation for the command push pointer 0/1
	case "pointer":
		if i == "0" {
			outputFile.WriteString("@THIS\nD=M\n") // D=*THIS
		} else { // i == "1"
			outputFile.WriteString("@THAT\nD=M\n") // D=*THAT
		}
		outputFile.WriteString("@SP\nA=M\nM=D\n") // *SP=D
	}
	// For all push commands increment stack pointer at the end
	outputFile.WriteString("@SP\nM=M+1\n") // SP++
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

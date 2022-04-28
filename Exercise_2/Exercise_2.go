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
	"strconv"
	"strings"
)

// Receive the path in program argument
var path = os.Args[1] // receiving the path as cli parameter
var pathArray = strings.Split(path, "\\")

// Create the output file with the according name
var fileName = pathArray[len(pathArray)-1]
var outputFile, _ = os.Create(fileName + ".asm")

// Global variable in order to distinguish the different labels (for the boolean and call operations)
var labelCount = 0
var callCount = 0

func main() {
	// Close the file "outputFile" at the end of the main function
	defer outputFile.Close()

	WriteInit() // Bootstrap code

	// Go through the file and performs some operations
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}
		fileName := info.Name()
		extension := filepath.Ext(fileName)
		if extension == ".vm" {
			fmt.Printf("File Name: %s\n", fileName)
			// removes the extension from the file name and writes it
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
				case "add", "sub", "neg":
					WriteArithmetic(command)
				// Boolean commands
				case "eq", "gt", "lt":
					WriteBoolean(command)
				// Logical commands - Bit-wise
				case "and", "or", "not":
					WriteLogical(command)
				// Memory access commands
				case "pop":
					segment := words[1]
					i := words[2]
					WritePop(segment, i, name)
				case "push":
					segment := words[1]
					i := words[2]
					WritePush(segment, i, name)
				// Program flow commands
				case "label":
					labelName := words[1]
					WriteLabel(labelName)
				case "goto":
					label := words[1]
					WriteGoto(label)
				case "if-goto":
					label := words[1]
					WriteIf(label)
				// Functions calling commands
				case "function": // declaration
					functionName := words[1]
					numLocals, err := strconv.Atoi(words[2])
					check(err)
					WriteFunction(functionName, numLocals)
				case "call":
					functionName := words[1]
					numArgs, err := strconv.Atoi(words[2])
					check(err)
					WriteCall(functionName, numArgs)
				case "return":
					WriteReturn()
				}
			}
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
		outputFile.WriteString("@SP\nM=M-1\nA=M\nD=M\n@result\nM=0\nM=M-D\n") // SP--, result=0-x
	}
	outputFile.WriteString("@result\nD=M\n@SP\nA=M\nM=D\n") // *SP=result
	outputFile.WriteString("@SP\nM=M+1\n")                  // SP++
}

// WriteBoolean Translation of boolean command (i.e. eq, gt or lt) in VM language to Hack language
func WriteBoolean(command string) {
	labelCount = labelCount + 1
	labelCountStr := strconv.Itoa(labelCount)
	switch command {
	case "eq": // Equality
		outputFile.WriteString("// eq\n")
		outputFile.WriteString("@SP\nM=M-1\nA=M\nD=M\n")                  // SP--, D=STACK[SP]
		outputFile.WriteString("@SP\nA=M-1\nD=M-D\n")                     // D=D-STACK[SP-1]
		outputFile.WriteString("@IF_TRUE_" + labelCountStr + "\nD;JEQ\n") // jump to (IF_TRUE) if D=0
	case "gt": // Greater than
		outputFile.WriteString("// gt\n")
		outputFile.WriteString("@SP\nM=M-1\nA=M\nD=M\n")                  // SP--, D=STACK[SP]
		outputFile.WriteString("@SP\nA=M-1\nD=M-D\n")                     // D=D-STACK[SP-1]
		outputFile.WriteString("@IF_TRUE_" + labelCountStr + "\nD;JGT\n") // jump to (IF_TRUE) if D>0
	case "lt": // Less than
		outputFile.WriteString("// lt\n")
		outputFile.WriteString("@SP\nM=M-1\nA=M\nD=M\n")                  // SP--, D=STACK[SP]
		outputFile.WriteString("@SP\nA=M-1\nD=M-D\n")                     // D=D-STACK[SP-1]
		outputFile.WriteString("@IF_TRUE_" + labelCountStr + "\nD;JLT\n") // jump to (IF_TRUE) if D<0
	}
	// If the condition is not met
	outputFile.WriteString("@SP\nA=M-1\nM=0\n")                        // *SP=0
	outputFile.WriteString("@IF_FALSE_" + labelCountStr + "\n0;JMP\n") // unconditional jump to (IF_FALSE)
	// Otherwise
	outputFile.WriteString("(IF_TRUE_" + labelCountStr + ")\n")  // declaring a label
	outputFile.WriteString("@SP\nA=M-1\nM=-1\n")                 // *SP=-1
	outputFile.WriteString("(IF_FALSE_" + labelCountStr + ")\n") // declaring a label
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
		outputFile.WriteString("A=A-1\nM=D|M\n")         // STACK[SP]=x or y
	case "not":
		outputFile.WriteString("// not\n")
		outputFile.WriteString("@SP\nA=M-1\nM=!M\n") // STACK[SP]=not(x)
	}
}

// WritePop Translation of pop command (in VM language) to Hack language
func WritePop(segment string, i string, programName string) {
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
		outputFile.WriteString("@SP\nM=M-1\n")                                           // SP--
		outputFile.WriteString("@SP\nA=M\nD=M\n")                                        // D=*SP
		outputFile.WriteString("@" + programName + "." + fileName + "." + i + "\nM=D\n") // static i = D
		return
	// Translation for the command pop temp i
	case "temp":
		outputFile.WriteString("@" + i + "\nD=A\n@5\nD=D+A\n@addr\nM=D\n") // addr=5+i
		outputFile.WriteString("@SP\nM=M-1\n")                             // SP--
	// Translation for the command pop pointer 0/1
	case "pointer":
		outputFile.WriteString("@SP\nM=M-1\n") // SP--
		if i == "0" {
			outputFile.WriteString("@THIS\nD=A\n@addr\nM=D\n") // addr=THIS
		} else { // i == "1"
			outputFile.WriteString("@THAT\nD=A\n@addr\nM=D\n") // addr=THAT
		}
	}
	// For all pop commands (except for static) add the value to the according address:
	outputFile.WriteString("@SP\nA=M\nD=M\n@addr\nA=M\nM=D\n") // *addr=*SP
}

// WritePush Translation of push command (in VM language) to Hack language
func WritePush(segment string, i string, programName string) {
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
		outputFile.WriteString("@" + programName + "." + fileName + "." + i + "\nD=M\n") // D = static i
		outputFile.WriteString("@SP\nA=M\nM=D\n")                                        // *SP=D
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

// WriteInit Writes the assembly code that effects the VM initialization also called "bootstrap code"
func WriteInit() {
	outputFile.WriteString("// Bootstrap code\n")
	outputFile.WriteString("@256\nD=A\n@SP\nM=D\n") // SP=256
	WriteCall("Sys.init", 0)
}

// WriteLabel Translation of label command in VM language to Hack language
func WriteLabel(labelName string) {
	outputFile.WriteString("// label " + labelName + "\n")
	outputFile.WriteString("(" + labelName + ")\n")
}

// WriteGoto Translation of goto command in VM language to Hack language
func WriteGoto(label string) {
	outputFile.WriteString("// goto " + label + "\n")
	outputFile.WriteString("@" + label + "\n0;JMP\n") // unconditional jump
}

// WriteIf Translation of if-goto command in VM language to Hack language
func WriteIf(label string) {
	outputFile.WriteString("// if-goto " + label + "\n")
	outputFile.WriteString("@SP\nM=M-1\nA=M\nD=M\n")  // SP--, D=*SP
	outputFile.WriteString("@" + label + "\nD;JNE\n") // if D!=0 jump to label
}

// WriteFunction Writes the assembly code that is the translation of the given function command
func WriteFunction(functionName string, numLocals int) {
	numLocalsStr := strconv.Itoa(numLocals)
	outputFile.WriteString("// function " + functionName + " " + numLocalsStr + "\n")
	WriteLabel(functionName) // declares a label for the function entry

	// Repeat numLocals times: push 0
	for i := 0; i < numLocals; i++ {
		WritePush("constant", "0", "") // initializes the local variables to 0
	}
}

// WriteCall Writes the assembly code that is the translation of the call command
func WriteCall(functionName string, numArgs int) {
	callCount = callCount + 1
	callCountStr := strconv.Itoa(callCount)
	numArgsStr := strconv.Itoa(numArgs)
	outputFile.WriteString("// call " + functionName + " " + numArgsStr + "\n")

	// Saving the caller's frame
	// push returnAddress
	outputFile.WriteString("@" + functionName + ".returnAddress_" + callCountStr + "\nD=A\n")
	outputFile.WriteString("@SP\nA=M\nM=D\n") // *SP=*returnAddress
	outputFile.WriteString("@SP\nM=M+1\n")    // SP++
	// push LCL - Saves LCL of the caller
	outputFile.WriteString("@LCL\nD=M\n@SP\nA=M\nM=D\n") // *SP=*returnAddress
	outputFile.WriteString("@SP\nM=M+1\n")               // SP++
	// push ARG - Saves ARG of the caller
	outputFile.WriteString("@ARG\nD=M\n@SP\nA=M\nM=D\n") // *SP=*returnAddress
	outputFile.WriteString("@SP\nM=M+1\n")               // SP++
	// push THIS - Saves THIS of the caller
	outputFile.WriteString("@THIS\nD=M\n@SP\nA=M\nM=D\n") // *SP=*returnAddress
	outputFile.WriteString("@SP\nM=M+1\n")                // SP++
	// push THAT - Saves THAT of the caller
	outputFile.WriteString("@THAT\nD=M\n@SP\nA=M\nM=D\n") // *SP=*returnAddress
	outputFile.WriteString("@SP\nM=M+1\n")                // SP++

	// Repositions ARG: ARG = SP-5-nArgs
	outputFile.WriteString("@SP\nD=M\n@5\nD=D-A\n")        // D = SP-5
	outputFile.WriteString("@" + numArgsStr + "\nD=D-A\n") // D = D-nArgs
	outputFile.WriteString("@ARG\nM=D\n")                  // ARG = D

	// Repositions LCL: LCL = SP
	outputFile.WriteString("@SP\nD=M\n@LCL\nM=D\n")

	// Transfers control to the called function
	WriteGoto(functionName) // goto functionName
	// Declares a label for the return-address
	WriteLabel(functionName + ".returnAddress_" + callCountStr)
}

// WriteReturn Writes the assembly code that is the translation of the return command
func WriteReturn() {
	outputFile.WriteString("// return\n")

	// endFrame is a temporary variable
	outputFile.WriteString("@LCL\nD=M\n@endFrame\nM=D\n") // endFrame=LCL
	// Gets the return address: retAddr = *(endFrame-5)
	outputFile.WriteString("@endFrame\nD=M\n@5\nD=D-A\n") // D = endFrame-5
	outputFile.WriteString("A=D\nD=M\n@retAddr\nM=D\n")   // retAddr = *(endFrame-5)

	WritePop("argument", "0", "")                     // *ARG=pop()
	outputFile.WriteString("@ARG\nD=M+1\n@SP\nM=D\n") // SP = ARG + 1
	// Restores caller's frame
	outputFile.WriteString("@endFrame\nA=M-1\nD=M\n@THAT\nM=D\n")          // THAT=*(endFrame-1)
	outputFile.WriteString("@2\nD=A\n@endFrame\nA=M-D\nD=M\n@THIS\nM=D\n") // THIS=*(endFrame-2)
	outputFile.WriteString("@3\nD=A\n@endFrame\nA=M-D\nD=M\n@ARG\nM=D\n")  // ARG=*(endFrame-3)
	outputFile.WriteString("@4\nD=A\n@endFrame\nA=M-D\nD=M\n@LCL\nM=D\n")  // LCL=*(endFrame-4)

	// goes to return address in the caller's code
	outputFile.WriteString("@retAddr\nA=M\n0;JMP\n")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

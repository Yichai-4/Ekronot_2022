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
	"unicode"
)

// Collection of keywords in the Jack language
var keywords = []string{"class", "constructor", "function", "method", "field", "static", "var",
	"int", "char", "boolean", "void", "true", "false", "null", "this", "let", "do", "if", "else", "while", "return"}

// Collection of keywords which could appear in an expression in the Jack language
var keywordsConstant = []string{"true", "false", "null", "this"}

// Collection of symbols in the Jack language
var symbols = []string{"{", "}", "(", ")", "[", "]", ".", ",", ";", "+", "-", "*", "/", "&", "|", "<", ">", "=", "~"}

// Collection of binary operators in the Jack language
var operators = []string{"+", "-", "*", "/", "&", "&amp;", "|", "<", "&lt;", ">", "&gt;", "="}

// Collection of unary operators in the Jack language
var unaryOperators = []string{"-", "~"}

// Collection of built-in types in the Jack language
var types = []string{"int", "char", "boolean"}

// Symbol table which saves the information about the variables of a class
var classSymbolTable [][]string

// Variables categories in the class-level symbol table
var classVarName, classVarType, classVarKind, classVarCount string
var generalClassName, className string

// Symbol table which saves the information about the variables of a subroutine
var subroutineSymbolTable [][]string

// Variables categories in the subroutine-level symbol table
var subVarName, subVarType, subVarKind, subVarCount string
var subroutineName, subroutineType, subroutineReturn string

// Label counter for compiling if and while statements
var labelCountIf, labelCountWhile int

// Saves the number of arguments of a specific function
var numOfArgs int

// For the different indentations according to the current scope
var indentation string

// Receive the path in program argument
var path = os.Args[1] // receiving the path as cli parameter
var pathArray = strings.Split(path, "\\")

// Create the output file with the according name
var directoryName = pathArray[len(pathArray)-1]
var fileName string

// Defines the file which will contain the tokens of the program
var tokensFile *os.File

// Defines the file which will contain the tokens of the program with the hierarchy
var parsedFile *os.File

// Defines the file which will contain the program written in the vm language
var vmFile *os.File

func main() {
	// Close the file "tokensFile" at the end of the main function
	defer tokensFile.Close()

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
			vmFile, _ = os.Create(directoryName + "_" + fileName + ".vm")
			vmFile.WriteString("// Program: " + fileName + ".jack\n")

			inputJackFile, _ := os.Open(path)
			defer inputJackFile.Close()

			Tokenize(path)

			Parse(tokensFile)
		}
		return nil
	})
}

// Tokenize Reads the input file, goes through it character by character,
// and writes all the tokens it contains (according to the Jack language syntax) in the output file
func Tokenize(jackFilePath string) {
	tokensFile, _ = os.Create(directoryName + "_" + fileName + "T.xml")
	tokensFile.WriteString("<tokens>\n") // the xml file (xxxT.xml) has to begin by "<tokens>"

	// Reads the input file and converts it into string of characters
	fileBuffer, err := ioutil.ReadFile(jackFilePath)
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
			if StringInList(token, keywords) {
				tokenClassification = "keyword"
			} else {
				tokenClassification = "identifier"
			}
			WriteToken(tokenClassification, token)
			// Case of a symbol next to keyword or identifier
			if nextChar != " " && StringInList(nextChar, symbols) {
				WriteToken("symbol", nextChar)
			}

		// Handles symbol
		case StringInList(char, symbols):
			tokenClassification = "symbol"
			token = GetSymbolToken(char, token)
			WriteToken(tokenClassification, token)

		// Handles integer constant
		case IsInteger(char):
			tokenClassification = "integerConstant"
			token, nextChar = GetIntegerToken(data, token, char)
			WriteToken(tokenClassification, token)
			// Case of a symbol next to integer
			if nextChar != " " && StringInList(nextChar, symbols) {
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
	tokensFile.WriteString("</tokens>\n") // the xml file (xxxT.xml) has to end by "</tokens>"

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

// GetKeywordIdentifierToken Returns the keyword or identifier of the current token and the next char
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

// GetSymbolToken Returns the according symbol of the current token
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

// GetIntegerToken Returns the integer value of the current token and the next char
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

// GetStringToken Returns the string value of the current token
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

// Parse Reads the input file with the tokens, goes through it line by line,
// and rewrites it with hierarchy and indentations (according to the rules of the Jack grammar) in the output file
func Parse(tokensFile *os.File) {
	parsedFile, _ = os.Create(directoryName + "_" + fileName + ".xml")
	tokensFile, _ = os.Open(directoryName + "_" + fileName + "T.xml")

	data := bufio.NewScanner(tokensFile)
	data.Scan() // scans the first line <token>
	data.Scan() // skips the first line <token> (doesn't appear in the output file)
	// Program structure
	CompileClass(data)
}

// CompileClass Compiles the class structure in the Jack language
// Syntax: 'class' className '{'classVarDec* subroutineDec*'}'
func CompileClass(data *bufio.Scanner) {
	parsedFile.WriteString("<class>\n")
	indentation = "  "

	Eat(data, "class")
	words := strings.Split(data.Text(), " ")
	generalClassName = words[1]
	CompileIdentifier(data) // checks class name
	Eat(data, "{")
	classSymbolTable = nil // resets the symbol table for the new class
	classVarCount = "0"
	// Compiles classVarDec*
	words = strings.Split(data.Text(), " ")
	nextToken := words[1]
	for nextToken == "static" || nextToken == "field" {
		CompileClassVarDec(data)
		words = strings.Split(data.Text(), " ")
		nextToken = words[1]
	}
	// Compiles subroutineDec*
	labelCountIf = 0
	labelCountWhile = 0
	words = strings.Split(data.Text(), " ")
	nextToken = words[1]
	for nextToken == "constructor" || nextToken == "function" || nextToken == "method" {
		CompileSubroutineDec(data)
		words = strings.Split(data.Text(), " ")
		nextToken = words[1]
	}
	Eat(data, "}")

	parsedFile.WriteString("</class>\n")
}

// CompileClassVarDec Compiles declaration of variables in a class in the Jack language
// Syntax: ('static'|'field') type varName (','varName)* ';'
func CompileClassVarDec(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<classVarDec>\n")
	indentation += "  "

	var firstWords = []string{"static", "field"}
	words := strings.Split(data.Text(), " ")
	oldKind := classVarKind
	classVarKind = words[1]
	// If the actual kind is different from the older kinder, then reset the counter
	if classVarCount != "0" && oldKind != classVarKind {
		classVarCount = "0"
	}
	EatOptions(data, firstWords)
	words = strings.Split(data.Text(), " ")
	classVarType = words[1]
	CompileType(data)
	words = strings.Split(data.Text(), " ")
	classVarName = words[1]
	CompileIdentifier(data) // checks variable name
	newRow := []string{classVarName, classVarType, classVarKind, classVarCount}
	classSymbolTable = append(classSymbolTable, newRow)
	classVarCount = IncStrCount(classVarCount)
	// Implements (',' varName)*
	words = strings.Split(data.Text(), " ")
	nextToken := words[1]
	for nextToken == "," { // same kind and type of variable
		Eat(data, ",")
		words = strings.Split(data.Text(), " ")
		classVarName = words[1]
		CompileIdentifier(data)
		newRow = []string{classVarName, classVarType, classVarKind, classVarCount}
		classSymbolTable = append(classSymbolTable, newRow)
		words = strings.Split(data.Text(), " ")
		nextToken = words[1]
		classVarCount = IncStrCount(classVarCount)
	}
	Eat(data, ";")

	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</classVarDec>\n")
}

// CompileType Compiles a type in the Jack language
// Syntax: 'int'|'char'|'boolean'|className
func CompileType(data *bufio.Scanner) {
	words := strings.Split(data.Text(), " ")
	tokenType := words[0]
	currentToken := words[1]
	if !StringInList(currentToken, types) && tokenType != "<identifier>" {
		print("error - expected one of them: ")
		fmt.Print(strings.Join(types, ", "))
		println(", <identifier>")
	} else {
		parsedFile.WriteString(indentation + data.Text() + "\n")
		data.Scan()
	}
}

// CompileSubroutineDec Compile declaration of subroutine in the Jack language
// Syntax: ('constructor', 'function', 'method') ('void'|type) subroutineName '('parameterList')' subroutineBody
func CompileSubroutineDec(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<subroutineDec>\n")
	indentation += "  "

	subroutineSymbolTable = nil // resets the symbol table for the new subroutine

	var firstWords = []string{"constructor", "function", "method"}
	words := strings.Split(data.Text(), " ")
	subroutineType = words[1]
	EatOptions(data, firstWords)
	words = strings.Split(data.Text(), " ")
	tokenType := words[0]
	subroutineReturn = words[1]
	if tokenType == "<identifier>" { // checks className
		CompileIdentifier(data)
	} else {
		var voidOrType = []string{"int", "char", "boolean", "void"}
		EatOptions(data, voidOrType)
	}
	words = strings.Split(data.Text(), " ")
	subroutineName = words[1]
	CompileIdentifier(data) // checks subroutine name
	Eat(data, "(")
	CompileParameterList(data)
	Eat(data, ")")
	CompileSubroutineBody(data)

	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</subroutineDec>\n")
}

// CompileParameterList Compile list of parameters in a subroutine in the Jack language
// Syntax: ((type varName) (',' type varName)*)?
func CompileParameterList(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<parameterList>\n")
	indentation += "  "

	var newRow []string
	subVarType = className
	subVarKind = "argument"
	subVarCount = "0"
	// Adds the special line for the method - current object
	if subroutineType == "method" {
		newRow = []string{"this", subVarType, subVarKind, subVarCount}
		subroutineSymbolTable = append(subroutineSymbolTable, newRow)
		subVarCount = IncStrCount(subVarCount)
	}
	// Compiles ((type varName) (',' type varName)*)?
	words := strings.Split(data.Text(), " ")
	tokenType := words[0]
	nextToken := words[1]
	if StringInList(nextToken, types) || tokenType == "<identifier>" {
		subVarType = nextToken
		CompileType(data)
		words = strings.Split(data.Text(), " ")
		subVarName = words[1]
		CompileIdentifier(data)
		newRow = []string{subVarName, subVarType, subVarKind, subVarCount}
		subroutineSymbolTable = append(subroutineSymbolTable, newRow)
		words = strings.Split(data.Text(), " ")
		nextToken = words[1]
		for nextToken == "," { // same kind of variable - argument
			subVarCount = IncStrCount(subVarCount)
			Eat(data, ",")
			words = strings.Split(data.Text(), " ")
			subVarType = words[1]
			CompileType(data)
			words = strings.Split(data.Text(), " ")
			subVarName = words[1]
			CompileIdentifier(data) // checks var name
			newRow = []string{subVarName, subVarType, subVarKind, subVarCount}
			subroutineSymbolTable = append(subroutineSymbolTable, newRow)
			words = strings.Split(data.Text(), " ")
			nextToken = words[1]
		}
	}

	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</parameterList>\n")
}

// CompileSubroutineBody Compiles the body of a subroutine in the Jack language
// Syntax: '{' varDec* statements '}'
func CompileSubroutineBody(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<subroutineBody>\n")
	indentation += "  "

	Eat(data, "{")
	subVarCount = "0"
	// Compiles varDec*
	words := strings.Split(data.Text(), " ")
	nextToken := words[1]
	for nextToken == "var" {
		CompileVarDec(data)
		words = strings.Split(data.Text(), " ")
		nextToken = words[1]
	}
	numOfLocals := GetCountOf(subroutineSymbolTable, "local")
	numOfLocalsStr := strconv.Itoa(numOfLocals)
	vmFile.WriteString("function " + generalClassName + "." + subroutineName + " " + numOfLocalsStr + "\n")
	// In case of constructor, finds a memory block of the required size
	// and returns its base address
	if subroutineType == "constructor" {
		subroutineReturn = generalClassName
		numOfWords := GetCountOf(classSymbolTable, "field")
		numOfWordsStr := strconv.Itoa(numOfWords)
		vmFile.WriteString("push constant " + numOfWordsStr + "\n")
		vmFile.WriteString("call Memory.alloc 1\n")
		vmFile.WriteString("pop pointer 0\n") // anchors this at the base address
	} else
	// In case of method, associates the "this" memory with the object on which the method is called to operate
	if subroutineType == "method" {
		vmFile.WriteString("push argument 0\n")
		vmFile.WriteString("pop pointer 0\n") // THIS = argument 0
	}
	CompileStatements(data)
	Eat(data, "}")

	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</subroutineBody>\n")
}

// IncStrCount Increments string counter using conversion between
func IncStrCount(countStr string) string {
	countInt, _ := strconv.Atoi(countStr)
	countInt += 1
	countStr = strconv.Itoa(countInt)
	return countStr
}

// GetCountOf Gets the number of occurrences of some kind in a symbol table
func GetCountOf(symbolTable [][]string, kind string) int {
	var count = 0
	for _, line := range symbolTable {
		if line[2] == kind {
			count += 1
		}
	}
	return count
}

// CompileVarDec Compiles declaration of variable/s in the Jack language
// Syntax: 'var' type varName (',' varName)* ';'
func CompileVarDec(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<varDec>\n")
	indentation += "  "

	subVarKind = "local"
	Eat(data, "var")
	words := strings.Split(data.Text(), " ")
	subVarType = words[1]
	CompileType(data)
	words = strings.Split(data.Text(), " ")
	subVarName = words[1]
	CompileIdentifier(data) // checks variable name
	newRow := []string{subVarName, subVarType, subVarKind, subVarCount}
	subroutineSymbolTable = append(subroutineSymbolTable, newRow)
	subVarCount = IncStrCount(subVarCount)
	// Compiles (',' varName)*
	words = strings.Split(data.Text(), " ")
	nextToken := words[1]
	for nextToken == "," { // same type and same kind of variable
		Eat(data, ",")
		words = strings.Split(data.Text(), " ")
		subVarName = words[1]
		CompileIdentifier(data) // checks var name
		newRow = []string{subVarName, subVarType, subVarKind, subVarCount}
		subroutineSymbolTable = append(subroutineSymbolTable, newRow)
		words = strings.Split(data.Text(), " ")
		nextToken = words[1]
		subVarCount = IncStrCount(subVarCount)
	}
	Eat(data, ";")

	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</varDec>\n")
}

// CompileStatements Compiles statements in the Jack language
// Syntax: (letStatement|ifStatement|whileStatement|doStatement|returnStatement)*
func CompileStatements(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<statements>\n")
	indentation += "  "

	var firstWords = []string{"let", "if", "while", "do", "return"}
	words := strings.Split(data.Text(), " ")
	nextToken := words[1]
	for StringInList(nextToken, firstWords) {
		switch nextToken {
		// Statements
		case "let":
			CompileLetStatement(data)
		case "if":
			CompileIfStatement(data)
		case "while":
			CompileWhileStatement(data)
		case "do":
			CompileDoStatement(data)
		case "return":
			CompileReturnStatement(data)
		}
		words = strings.Split(data.Text(), " ")
		nextToken = words[1]
	}

	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</statements>\n")
}

// CompileLetStatement Compiles a let statement in the Jack language
// Syntax: 'let' varName ('[' expression']')? '=' expression ';'
func CompileLetStatement(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<letStatement>\n")
	indentation += "  "

	numOfArgs = 0
	Eat(data, "let") // code to handle 'let'
	words := strings.Split(data.Text(), " ")
	currentToken := words[1]
	lineIndex, flag := IsInTable(currentToken, subroutineSymbolTable)
	var kind, count string
	if flag {
		kind = subroutineSymbolTable[lineIndex][2]
		count = subroutineSymbolTable[lineIndex][3]
	} else {
		lineIndex, flag = IsInTable(currentToken, classSymbolTable)
		kind = classSymbolTable[lineIndex][2]
		count = classSymbolTable[lineIndex][3]
		if flag && kind == "field" {
			kind = "this"
		}
	}
	CompileIdentifier(data) // checks var name
	// Compiles ('[' expression ']')?
	words = strings.Split(data.Text(), " ")
	currentToken = words[1]
	// arr[expression1] = expression2
	if currentToken == "[" {
		vmFile.WriteString("push " + kind + " " + count + "\n") // base address of the new object
		Eat(data, "[")
		CompileExpression(data)
		Eat(data, "]")
		vmFile.WriteString("add\n") // top stack value = RAM address of arr[expression1]
		Eat(data, "=")
		CompileExpression(data)
		vmFile.WriteString("pop temp 0\n") // temp 0 = the value of expression2
		vmFile.WriteString("pop pointer 1\npush temp 0\npop that 0\n")
	} else {
		Eat(data, "=")
		CompileExpression(data)
		vmFile.WriteString("pop " + kind + " " + count + "\n") // base address of the new object
	}
	Eat(data, ";")

	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</letStatement>\n")
}

// CompileIfStatement Compiles an if (and else) statement in the Jack language
// Syntax: 'if' '('expression')' '{'statements'}' ('else' '{'statements'}')?
func CompileIfStatement(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<ifStatement>\n")
	indentation += "  "

	labelCountStr := strconv.Itoa(labelCountIf)
	Eat(data, "if") // code to handle 'if'
	Eat(data, "(")
	CompileExpression(data)
	Eat(data, ")")
	vmFile.WriteString("not\n")
	vmFile.WriteString("if-goto IF_FALSE" + labelCountStr + "\n")
	labelCountIf += 1
	Eat(data, "{")
	CompileStatements(data)
	Eat(data, "}")
	vmFile.WriteString("goto IF_END" + labelCountStr + "\n")
	vmFile.WriteString("label IF_FALSE" + labelCountStr + "\n")
	// Compiles ('else' '{'statements'}')?
	words := strings.Split(data.Text(), " ")
	currentToken := words[1]
	if currentToken == "else" {
		parsedFile.WriteString(indentation + data.Text() + "\n")
		data.Scan()
		Eat(data, "{")
		CompileStatements(data)
		Eat(data, "}")
	}
	vmFile.WriteString("label IF_END" + labelCountStr + "\n")

	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</ifStatement>\n")
}

// CompileWhileStatement Compiles a while statement in the Jack language
// Syntax: 'while' '('expression')' '{'statements'}'
func CompileWhileStatement(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<whileStatement>\n")
	indentation += "  "

	labelCountStr := strconv.Itoa(labelCountWhile)
	vmFile.WriteString("label WHILE_EXP" + labelCountStr + "\n")
	Eat(data, "while") // code to handle 'while'
	Eat(data, "(")
	CompileExpression(data)
	Eat(data, ")")
	vmFile.WriteString("not\n")
	vmFile.WriteString("if-goto WHILE_END" + labelCountStr + "\n")
	labelCountWhile += 1
	Eat(data, "{")
	CompileStatements(data)
	vmFile.WriteString("goto WHILE_EXP" + labelCountStr + "\n")
	Eat(data, "}")
	vmFile.WriteString("label WHILE_END" + labelCountStr + "\n")

	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</whileStatement>\n")
}

// CompileDoStatement Compiles a do statement in the Jack language
// Syntax: 'do' subroutineCall ';'
func CompileDoStatement(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<doStatement>\n")
	indentation += "  "

	numOfArgs = 0
	Eat(data, "do") // code to handle 'do'
	// Compiles subroutine call
	words := strings.Split(data.Text(), " ")
	subroutineName = words[1]
	CompileIdentifier(data)
	words = strings.Split(data.Text(), " ")
	currentToken := words[1]
	if currentToken == "(" {
		// First push a reference to the object on which the method is supposed to operate
		vmFile.WriteString("push pointer 0\n")
		Eat(data, "(")
		CompileExpressionList(data)
		Eat(data, ")")
		numOfArgs += 1
		numOfArgsStr := strconv.Itoa(numOfArgs)
		vmFile.WriteString("call " + generalClassName + "." + subroutineName + " " + numOfArgsStr + "\n")
	} else {
		className = subroutineName
		if unicode.IsLower(rune(className[0])) { // obj.foo(x1, x2, ...)
			objectName := className
			lineIndex, flag := IsInTable(objectName, subroutineSymbolTable)
			var kind, count, objectType string
			if flag {
				objectType = subroutineSymbolTable[lineIndex][1]
				kind = subroutineSymbolTable[lineIndex][2]
				count = subroutineSymbolTable[lineIndex][3]
			} else {
				lineIndex, flag = IsInTable(objectName, classSymbolTable)
				objectType = classSymbolTable[lineIndex][1]
				kind = classSymbolTable[lineIndex][2]
				count = classSymbolTable[lineIndex][3]
				if flag && kind == "field" {
					kind = "this"
				}
			}
			vmFile.WriteString("push " + kind + " " + count + "\n")
			numOfArgs += 1
			className = objectType
		}
		Eat(data, ".")
		words = strings.Split(data.Text(), " ")
		subroutineName = words[1]
		CompileIdentifier(data)
		Eat(data, "(")
		CompileExpressionList(data)
		Eat(data, ")")
		numOfArgsStr := strconv.Itoa(numOfArgs)
		vmFile.WriteString("call " + className + "." + subroutineName + " " + numOfArgsStr + "\n")
	}
	// The caller of a void method must dump the returned value
	vmFile.WriteString("pop temp 0\n")
	Eat(data, ";")

	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</doStatement>\n")
}

// CompileReturnStatement Compiles a return statement in the Jack language to vm language
// Syntax: 'return' expression? ';'
func CompileReturnStatement(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<returnStatement>\n")
	indentation += "  "

	Eat(data, "return")
	// Compiles "expression?"
	if IsTerm(data) {
		CompileExpression(data)
	}
	Eat(data, ";")
	if subroutineReturn == "void" {
		// Methods must return a value
		vmFile.WriteString("push constant 0\n")
	}
	vmFile.WriteString("return\n")

	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</returnStatement>\n")
}

// IsTerm Checks if the current token is a term according to the Jack language
func IsTerm(data *bufio.Scanner) bool {
	words := strings.Split(data.Text(), " ")
	tokenType := words[0]
	currentToken := words[1]
	if IsInteger(currentToken) || tokenType == "<stringConstant>" || StringInList(currentToken, keywordsConstant) ||
		tokenType == "<identifier>" || currentToken == "(" || StringInList(currentToken, unaryOperators) {
		return true
	}
	return false
}

// CompileExpression Compiles an expression in the Jack language
// Syntax: term (op term)*
func CompileExpression(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<expression>\n")
	indentation += "  "

	CompileTerm(data)
	words := strings.Split(data.Text(), " ")
	nextToken := words[1]
	for StringInList(nextToken, operators) {
		var operator string
		if StringInList(nextToken, operators) {
			operator = nextToken
		}
		Eat(data, nextToken)
		CompileTerm(data)
		WriteOperator(operator)
		words = strings.Split(data.Text(), " ")
		nextToken = words[1]
	}

	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</expression>\n")
}

// WriteOperator Translates an operator in the Jack language to the vm language and writes it in the vm file
func WriteOperator(operator string) {
	var vmOperator string
	switch operator {
	case "+":
		vmOperator = "add"
	case "-":
		vmOperator = "sub"
	case "=":
		vmOperator = "eq"
	case ">", "&gt;":
		vmOperator = "gt"
	case "<", "&lt;":
		vmOperator = "lt"
	case "&", "&amp;":
		vmOperator = "and"
	case "|":
		vmOperator = "or"
	case "/":
		vmOperator = "call Math.divide 2"
	case "*":
		vmOperator = "call Math.multiply 2"
	}
	vmFile.WriteString(vmOperator + "\n")
}

// CompileTerm Compiles a term in the Jack language
// Syntax: integerConstant| stringConstant|keywordConstant|varName|varName'['expression']'|
//			subroutineCall|'('expression')'|unaryOp term
func CompileTerm(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<term>\n")
	indentation += "  "

	words := strings.Split(data.Text(), " ")
	tokenType := words[0]
	currentToken := words[1]
	switch {
	case tokenType == "<integerConstant>":
		vmFile.WriteString("push constant " + currentToken + "\n")
		parsedFile.WriteString(indentation + data.Text() + "\n")
		data.Scan()
	case tokenType == "<stringConstant>":
		tokenLine := data.Text()
		stringToken := strings.TrimLeft(tokenLine, "<stringConstant> ")
		stringToken = strings.TrimRight(stringToken, "</stringConstant>")
		//strLength := strings.Count(stringToken, "")
		strLength := len(stringToken)
		stringToken = stringToken[:strLength-1] // removes one space in right
		strLength -= 1                          // updates the length
		strLengthStr := strconv.Itoa(strLength)
		vmFile.WriteString("push constant " + strLengthStr + "\n")
		vmFile.WriteString("call String.new 1\n")
		i := 0
		for i < strLength {
			temp := int(stringToken[i])
			tempStr := strconv.Itoa(temp)
			vmFile.WriteString("push constant " + tempStr + "\n")
			vmFile.WriteString("call String.appendChar 2\n")
			i += 1
		}
		parsedFile.WriteString(indentation + data.Text() + "\n")
		data.Scan()
	// Compiling constants
	case StringInList(currentToken, keywordsConstant):
		switch currentToken {
		case "null", "false":
			vmFile.WriteString("push constant 0\n")
		case "true":
			vmFile.WriteString("push constant 1\nneg\n") // push constant -1
		case "this":
			vmFile.WriteString("push pointer 0\n")
		}
		parsedFile.WriteString(indentation + data.Text() + "\n")
		data.Scan()
	case tokenType == "<identifier>": // varName | varName '['expression']' | subroutineCall
		words = strings.Split(data.Text(), " ")
		currentToken = words[1]
		lineIndex, flag := IsInTable(currentToken, subroutineSymbolTable)
		var kind, count, objectType string
		if flag {
			objectType = subroutineSymbolTable[lineIndex][1]
			kind = subroutineSymbolTable[lineIndex][2]
			count = subroutineSymbolTable[lineIndex][3]
			if kind == "field" {
				kind = "this"
			}
			vmFile.WriteString("push " + kind + " " + count + "\n")
		} else {
			lineIndex, flag = IsInTable(currentToken, classSymbolTable)
			if flag {
				objectType = classSymbolTable[lineIndex][1]
				kind = classSymbolTable[lineIndex][2]
				count = classSymbolTable[lineIndex][3]
				if kind == "field" {
					kind = "this"
				}
				vmFile.WriteString("push " + kind + " " + count + "\n")
			}
		}
		subClassTemp := currentToken
		CompileIdentifier(data)
		words = strings.Split(data.Text(), " ")
		nextToken := words[1]
		if nextToken == "[" { // varName '['expression']'
			// push base address already has been written
			Eat(data, "[")
			CompileExpression(data) // offset
			Eat(data, "]")
			vmFile.WriteString("add\n")
			vmFile.WriteString("pop pointer 1\npush that 0\n")
			// Compiles subroutine call
			// Syntax: subroutineName '('expressionList')' | (className|varName)'.'subroutineName '('expressionList')'
		} else if nextToken == "(" { // subroutineName '('expressionList')'
			subroutineName = subClassTemp
			Eat(data, "(")
			CompileExpressionList(data)
			Eat(data, ")")
			vmFile.WriteString("call " + subroutineName + "\n") // output "call f"
		} else if nextToken == "." { // (className | varName)'.'subroutineName '('expressionList')'
			className = subClassTemp
			if unicode.IsLower(rune(subClassTemp[0])) {
				className = objectType
				numOfArgs += 1
			}
			Eat(data, ".")
			words = strings.Split(data.Text(), " ")
			subroutineName = words[1]
			CompileIdentifier(data)
			Eat(data, "(")
			CompileExpressionList(data)
			Eat(data, ")")
			numOfArgsStr := strconv.Itoa(numOfArgs)
			vmFile.WriteString("call " + className + "." + subroutineName + " " + numOfArgsStr + "\n")
		}
	case currentToken == "(": // '('expression')'
		Eat(data, "(")
		CompileExpression(data)
		Eat(data, ")")
	case StringInList(currentToken, unaryOperators): // unaryOp term
		var vmOperator string
		operator := currentToken
		Eat(data, currentToken)
		CompileTerm(data)
		if operator == "~" {
			vmOperator = "not"
		} else if operator == "-" {
			vmOperator = "neg"
		}
		vmFile.WriteString(vmOperator + "\n") // output "op"
	}

	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</term>\n")
}

// IsInTable Checks if an item is in a symbol table and returns the matched line in the table
func IsInTable(item string, table [][]string) (int, bool) {
	rowIndex := 0
	for _, line := range table {
		if line[0] == item { // checks only in the column of item name
			return rowIndex, true
		}
		rowIndex += 1
	}
	return rowIndex, false
}

// CompileExpressionList Compiles an expression list in the Jack language
// Syntax: (expression(',' expression)*)?
func CompileExpressionList(data *bufio.Scanner) {
	parsedFile.WriteString(indentation + "<expressionList>\n")
	indentation += "  "

	words := strings.Split(data.Text(), " ")
	if IsTerm(data) {
		numOfArgs += 1
		CompileExpression(data)
		words = strings.Split(data.Text(), " ")
		nextToken := words[1]
		for nextToken == "," {
			numOfArgs += 1
			Eat(data, ",")
			CompileExpression(data)
			words = strings.Split(data.Text(), " ")
			nextToken = words[1]
		}
	}

	indentation = indentation[2:]
	parsedFile.WriteString(indentation + "</expressionList>\n")
}

// CompileIdentifier Compiles an identifier in the Jack language and writes it (token) in the output file.
// At the end it advances the scanner to the next line in the tokens file
func CompileIdentifier(data *bufio.Scanner) {
	words := strings.Split(data.Text(), " ")
	tokenType := words[0]
	currentToken := words[1]
	if tokenType != "<identifier>" {
		println("error - expected identifier " + "got " + currentToken)
	} else {
		parsedFile.WriteString(indentation + data.Text() + "\n")
		data.Scan()
	}
}

// Eat Reads the current token, checks if it's matching to the string s and if it's so, writes it in the output file.
// At the end it advances the scanner to the next line in the tokens file
func Eat(data *bufio.Scanner, s string) {
	words := strings.Split(data.Text(), " ")
	currentToken := words[1]
	if currentToken != s {
		println("error - expected " + s + " got " + currentToken)
	} else {
		parsedFile.WriteString(indentation + data.Text() + "\n")
		data.Scan()
	}
}

// EatOptions Like Eat function but with several options of matching
func EatOptions(data *bufio.Scanner, start []string) {
	words := strings.Split(data.Text(), " ")
	currentToken := words[1]
	if !StringInList(currentToken, start) {
		print("error - expected one of them: ")
		fmt.Print(strings.Join(start, ", "))
		println(" got " + currentToken)
	} else {
		parsedFile.WriteString(indentation + data.Text() + "\n")
		data.Scan()
	}
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

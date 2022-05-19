package main

import "os"

var keyword = [...]string{"class", "constructor", "function", "method", "field", "static", "var",
	"int", "char", "boolean", "void", "true", "false", "null", "this", "let", "do", "if", "else", "while", "return"}

var symbol = [...]string{"{", "}", "(", ")", "[", "]", ".", ",", ";", "+", "-", "*", "/", "&", "|", "<", ">", "=", "~"}

//var integerConstant = range [32768]int{}

type JackTokenizer interface {
	func advance()
	func hasMoreTokens()
}

func New(fileName string) JackTokenizer {
	var outputFile, _ = os.Create(fileName + "T.xml")
	inputFile, err := os.Open(fileName)
	check(err)
}

func advance() {

}

func hasMoreTokens() {
	//return io.EOF
}
func main() {
	tknzr = new JackTokenizer()
	tknzr.advance()
	flag = hasMoreTokens()
	outputFile.WriteString("<tokens>")
	for flag {
		tokenClassification := currentToken
		outputFile.WriteString("< " + tokenClassification + " >")
		outputFile.WriteString(currentToken.value)
		outputFile.WriteString("</ " + tokenClassification + " >\n")
		tknzr.advance()
	}
	outputFile.WriteString("</tokens>")
}

package main

type JackParser interface {
	func compileStatements()
	func compileIfStatement()
	func compileWhileStatement()
	func compileLetStatement()
	func compileExpression()
	func compileTerm()
	func compileVarName()
}

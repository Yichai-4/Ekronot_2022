// Submitters: Yichaï Hazan . ID: 1669535
//			   Amitai Shmeeda. ID: 305361479

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

// A function error that check immediately before proceeding to the next steps.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Global variables
// Receive the path in program argument
var path = os.Args[1]                                 // receiving the path as cli parameter
var pathArray = strings.Split(path, "\\")             // split the path to words in the array
var directoryName = pathArray[len(pathArray)-1]       // take the last word of "Tar0"
var outputFile, _ = os.Create(directoryName + ".asm") // change the output of the fil to be .asm

var totalBuy = float64(0)
var totalCell = float64(0)

func main() {
	fmt.Println("Hello World !")

	defer outputFile.Close() // will close the file when we are done

	// function Walk for read a file - need the path to the file, info of file, error if there is nothing.
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// if the file is empty
		if err != nil {
			log.Fatalf(err.Error())
		}

		var fileName = info.Name()
		var extension = filepath.Ext(fileName) // take the end of the file type name

		if extension == ".vm" {
			fmt.Printf("File Name: %s\n", fileName) // print the name of the file

			var name = strings.TrimRight(fileName, extension) // take off the end of type file

			outputFile.WriteString(name + "\n")

			inputFile, err := os.Open(path)
			check(err)
			defer inputFile.Close() // will close the file when we are done

			scanner := bufio.NewScanner(inputFile) // scan the file (read the file)

			// loop for reading the file
			for scanner.Scan() {
				var words = strings.Split(scanner.Text(), " ") // split the words in the file
				var firstWord = words[0]                       // take the first word in the line
				var productName = words[1]                     // take the second word in the line
				var amount, err = strconv.Atoi(words[2])       // take the third word in the line and convert it to int
				check(err)
				price, err := strconv.ParseFloat(words[3], 64) //// take the fourth word in the line and convert to float

				// check the first word and call the func HandleBuy\Sell
				if firstWord == "buy" {
					HandleBuy(productName, amount, price)
				}

				if firstWord == "cell" {
					HandleSell(productName, amount, price)
				}
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}

		}
		return nil
	})

	// print
	outputFile.WriteString("TOTAL BUY: " + fmt.Sprintf("%f\n", totalBuy))
	outputFile.WriteString("TOTAL CELL: " + fmt.Sprintf("%f\n", totalCell))

}

func HandleBuy(ProductName string, Amount int, Price float64) {

	outputFile.WriteString("### BUY " + ProductName + " ###\n")
	var totalPrice = float64(Amount) * Price
	var priceStr = fmt.Sprintf("%f", totalPrice)
	outputFile.WriteString(priceStr + "\n")
	totalBuy = totalBuy + totalPrice
}

func HandleSell(ProductName string, Amount int, Price float64) {

	outputFile.WriteString("$$$ CELL " + ProductName + " $$$\n")
	var totalPrice = float64(Amount) * Price
	var n = fmt.Sprintf("%f", totalPrice)
	outputFile.WriteString(n + "\n")
	totalCell = totalCell + totalPrice
}

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
var directoryName = pathArray[len(pathArray)-1]
var outputFile, _ = os.Create(directoryName + ".asm")

// Define global variables
var totalBuy = float64(0)
var totalCell = float64(0)

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
			outputFile.WriteString(name + "\n")

			inputFile, err := os.Open(path)
			check(err)
			defer inputFile.Close()

			scanner := bufio.NewScanner(inputFile)

			for scanner.Scan() {
				words := strings.Split(scanner.Text(), " ")
				firstWord := words[0]
				productName := words[1]
				// Converts string to int
				amount, err := strconv.Atoi(words[2])
				check(err)
				// Converts string to float
				price, _ := strconv.ParseFloat(words[3], 64)

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
	outputFile.WriteString("TOTAL BUY: " + fmt.Sprintf("%.1f\n", totalBuy))
	outputFile.WriteString("TOTAL CELL: " + fmt.Sprintf("%.1f\n", totalCell))

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func HandleBuy(ProductName string, Amount int, Price float64) {
	outputFile.WriteString("### BUY " + ProductName + " ###\n")
	totalPrice := float64(Amount) * Price
	priceInStr := fmt.Sprintf("%.1f", totalPrice)
	outputFile.WriteString(priceInStr + "\n")
	totalBuy = totalBuy + totalPrice
}

func HandleSell(ProductName string, Amount int, Price float64) {
	outputFile.WriteString("$$$ CELL " + ProductName + " $$$\n")
	totalPrice := float64(Amount) * Price
	priceInStr := fmt.Sprintf("%.1f", totalPrice)
	outputFile.WriteString(priceInStr + "\n")
	totalCell = totalCell + totalPrice
}

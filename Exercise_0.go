// Submitter: Yichaï Hazan
// ID: 1669535

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

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Receive the path in program argument
var path = os.Args[1] // receiving the path as cli parameter
var pathArray = strings.Split(path, "\\")
var directoryName = pathArray[len(pathArray)-1]
var outputFile, _ = os.Create(directoryName + ".asm")

var totalBuy = float64(0)
var totalCell = float64(0)

func main() {
	fmt.Println("Hello World !")

	defer outputFile.Close()
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}
		var fileName = info.Name()
		var extension = filepath.Ext(fileName)
		if extension == ".vm" {
			fmt.Printf("File Name: %s\n", fileName)
			var name = strings.TrimRight(fileName, extension)
			outputFile.WriteString(name + "\n")

			inputFile, err := os.Open(path)
			check(err)
			defer inputFile.Close()

			scanner := bufio.NewScanner(inputFile)

			for scanner.Scan() {
				var words = strings.Split(scanner.Text(), " ")
				var firstWord = words[0]
				var productName = words[1]
				var amount, err = strconv.Atoi(words[2])
				check(err)
				price, err := strconv.ParseFloat(words[3], 64)
				if firstWord == "buy" {
					HandleBuy(productName, amount, price)
				}
				if firstWord == "sell" {
					HandleSell(productName, amount, price)
				}
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}

		}
		return nil
	})
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

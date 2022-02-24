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
func main() {
	fmt.Println("Hello World !")
	/* 1) Receive the path from the user
	fmt.Println("Enter a path: ")
	var path string
	fmt.Scanln(&path)
	*/

	// 2) Receive the path in program argument
	path := os.Args[1]
	pathArray := strings.Split(path, "\\")
	directoryName := pathArray[len(pathArray)-1]
	fmt.Println(directoryName)
	outputFile, err := os.Create(directoryName + ".asm")
	check(err)
	defer outputFile.Close()
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}
		var fileName = info.Name()
		var extension = filepath.Ext(fileName)
		if extension == ".vm" {
			fmt.Printf("File Name: %s\n", info.Name())
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

}

func HandleBuy(ProductName string, Amount int, Price float64) {

}

func HandleSell(ProductName string, Amount int, Price float64) {

}

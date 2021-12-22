package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strings"
)

// Declaring constants
const lowerCase = "abcdefghijklmnopqrstuvwxyz"
const upperCase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const digits = "0123456789"
const chars = "\"'|@#&-!;,?.:/%[]{}()<>\\~=+*_$ "

// Declaring variable
var allowedChars string
var calculations int

// Returns a listOfChars that contains all of the allowed characters.
// For example: getChars("LUDC"), returns string with all lowercase letters, uppercase letters, digits and other characters
func getChars(charPattern string) string {
	var listOfChars string
	if strings.Contains(charPattern, "L") {
		listOfChars += lowerCase
	}
	if strings.Contains(charPattern, "U") {
		listOfChars += upperCase
	}
	if strings.Contains(charPattern, "D") {
		listOfChars += digits
	}
	if strings.Contains(charPattern, "C") {
		listOfChars += chars
	}
	return listOfChars
}

// Returns the next sequence based on the given sequence
func next(str string) string {
	if (len(str) <= 0) {
		str += string(allowedChars[0])
	} else {
		str = str[:0] + string(allowedChars[(strings.Index(allowedChars, string(str[0])) + 1 ) % len(allowedChars)]) + str [1:] // Increases the first character of the string by 1
		if (strings.Index(allowedChars, string(str[0])) == 0) { // Checks if the character is equal to the first character of the allowedChars string
			ret := string(str[0]) // Stores the first value of the string
			nxt := next(str[1:]) // Preforms the next function on a slice of the full string
			return ret + nxt // Returns the first value of current string + the next sequence of the sliced off part of the full string
		}
	}
	return str
}

// Gets sequence from jobs channel and returns the next sequence on the results channel
func worker(workerId int, workerAmt *int, verbose bool, hash *string, result chan<- string) {
	// Declaring necessary variables
	var bytes [32]byte
	var sequence string 

	// Initializing sequence variable
	for i := 0; i < workerId + 1; i++ {
		sequence = next(sequence)
	}

	if (verbose) { // Checks for verbose flag
		for { 
			calculations += 1 // Adds 1 to calculations variable 
			fmt.Println(sequence) // Prints every generated sequence
			bytes = sha256.Sum256([]byte(sequence)) // Assign SHA256 sum of sequence to bytes variable
			if (hex.EncodeToString(bytes[:]) == *hash) { // Checks if hashes match and breaks out of for loop
				break
			}
			for i := 0; i < *workerAmt + 1; i++ { // Generates the (workerAmt + 1st) next sequence
				sequence = next(sequence)
			}
		}
	} else {
		for { 
			calculations += 1 // Adds 1 to calculations variable
			bytes = sha256.Sum256([]byte(sequence)) // Assign SHA256 sum of sequence to bytes variable
			if (hex.EncodeToString(bytes[:]) == *hash) { // Checks if hashes match and breaks out of for loop
				break
			}
			for i := 0; i < *workerAmt + 1; i++ { // Generates the (workerAmt + 1st) next sequence
				sequence = next(sequence)
			}
		}
	}

	result <- sequence

	close(result)
}

func main() {
	verbosePtr := flag.Bool("verbose", false, "set to true if you want command line output")
	flag.Parse() // Checks for verbose flags and relays info to workers when program starts

	allowedChars = getChars("LUDC")

	fmt.Print("Enter hash: ") // Writes "Enter hash: " to the command line and waits for input
	reader := bufio.NewReader(os.Stdin)
    hash, _ := reader.ReadString('\n') // Waits for \n character
    hash = strings.Replace(hash, "\n", "", -1) // Replace "\n" character with ""

	// Declares the amount of workers
	var workers int = 2

	// Creating channels
	result := make(chan string, 1)

	// Creating workers
	for i := 0; i < workers; i++ {
		go worker(i, &workers, *verbosePtr, &hash, result)
	}
	
	// Prints results to terminal
	fmt.Println("Cracked hash with: " + <-result)
	fmt.Println("Performed", calculations, "calculations to crack hash")
}


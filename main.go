package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
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

// Returns a listOfChars that contains all the allowed characters.
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

// Returns the i-th next sequence based on the given sequence
func next(str string, i int) string {
	if len(str) <= 0 {
		str += string(allowedChars[i])
	} else {
		index := strings.Index(allowedChars, string(str[0])) + i
		if index >= len(allowedChars) {
			index = index - len(allowedChars)
		}
		str = str[:0] + string(allowedChars[index]) + str[1:] // Increases the first character of the string by 1
		if strings.Index(allowedChars, string(str[0])) == 0 { // Checks if the character is equal to the first character of the allowedChars string
			ret := string(str[0])   // Stores the first value of the string
			nxt := next(str[1:], 1) // Preforms the next function on a slice of the full string
			return ret + nxt        // Returns the first value of current string + the next sequence of the sliced off part of the full string
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
	sequence = next(sequence, workerId)

	if verbose { // Checks for verbose flag
		for {
			calculations += 1                          // Adds 1 to calculations variable
			fmt.Println(sequence)                      // Prints every generated sequence
			bytes = sha256.Sum256([]byte(sequence))    // Assign SHA256 sum of sequence to bytes variable
			if hex.EncodeToString(bytes[:]) == *hash { // Checks if hashes match and breaks out of for loop
				break
			}
			// Generates the workerAmt-th next sequence
			sequence = next(sequence, *workerAmt)
		}
	} else {
		for {
			calculations += 1                          // Adds 1 to calculations variable
			bytes = sha256.Sum256([]byte(sequence))    // Assign SHA256 sum of sequence to bytes variable
			if hex.EncodeToString(bytes[:]) == *hash { // Checks if hashes match and breaks out of for loop
				break
			}
			// Generates the workerAmt-th next sequence
			sequence = next(sequence, *workerAmt)
		}
	}

	// Send result gracefully to result channel
	defer func() {
		if recover() != nil {
			return
		}
	}()

	result <- sequence
	close(result)
}

func crack(hash string, charPattern string, workerAmt int, verbose bool) string {
	// Declaring necessary variables
	// Channel for getting results from workers
	var result = make(chan string)

	// Initializing allowedChars variable
	allowedChars = getChars(charPattern)

	// Creating workers
	for i := 0; i < workerAmt; i++ {
		go worker(i, &workerAmt, verbose, &hash, result)
	}

	// Print result to terminal
	resultStr := <-result
	fmt.Println("Result:", resultStr)
	return resultStr
}

// Gets the hash and charPattern from the terminal/user
func main() {
	// Declaring necessary variables
	var hash string

	// Get variable values from terminal flags and save them to designated variables
	verbose := flag.Bool("verbose", false, "Prints every generated sequence")
	charPattern := flag.String("charPattern", "LUDC", "Specifies which characters to use. L = lowercase, U = uppercase, D = digits, C = other characters")
	workerAmt := flag.Int("workerAmt", 3, "Specifies the amount of workers to use")

	// Parse flags
	flag.Parse()

	// Getting hash from user
	fmt.Print("Enter hash: ")
	_, err := fmt.Scanln(&hash)
	if err != nil {
		return
	}

	// Cracking the hash
	crack(hash, *charPattern, *workerAmt, *verbose)

	// Print calculations to terminal
	fmt.Println("Calculations:", calculations)
}

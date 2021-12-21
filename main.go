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

// Matches hash of sequence to the given hash, if they don't match the function tries to match the next sequence and so on...
func crack(hash string, verbose bool) {
	var sequence string
	var bytes [32]byte // Had to declare this variable to be able to compare both hashes
	if (verbose) {
		for hex.EncodeToString(bytes[:]) != hash {
			sequence = next(sequence)
			fmt.Println(sequence)
			bytes = sha256.Sum256([]byte(sequence)) // Assigns the SHA256 sum of the sequence to the bytes variable
		}
	} else {
		for hex.EncodeToString(bytes[:]) != hash {
			sequence = next(sequence)
			bytes = sha256.Sum256([]byte(sequence)) // Assigns the SHA256 sum of the sequence to the bytes variable
		}
		fmt.Println(sequence) // If the hashes match, it prints the sequence to the terminal
	}
}

func main() {
	verbosePtr := flag.Bool("verbose", false, "set to true if you want command line output")
	flag.Parse()
	allowedChars = getChars("LUDC")
	fmt.Print("Enter hash: ") // Writes "Enter hash: " to the command line and waits for input
	reader := bufio.NewReader(os.Stdin)
    hash, _ := reader.ReadString('\n') // Waits for \n character
    hash = strings.Replace(hash, "\n", "", -1) // Replace "\n" character with ""
	crack(hash, *verbosePtr) // Calls crack function
}


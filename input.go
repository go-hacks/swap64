// Everything related to user input

package main

import (
	"fmt"
	"golang.org/x/term"
	"syscall"
	"os"
	"github.com/DavidGamba/go-getoptions"
)

func parseParams () (string, string, int, int) {
	var threadCnt int
	var roundCnt int
	opt := getoptions.New()
	opt.Bool("help", false, opt.Alias("h", "?"))
	opt.IntVar(&threadCnt, "threads", 0, opt.Alias("t"))
	opt.IntVar(&roundCnt, "rounds", 7, opt.Alias("r"))
	opt.StringVar(&opMode, "opmode", "", opt.Alias("o"))
	_, err := opt.Parse(os.Args[1:])
	synString := opt.Help(getoptions.HelpSynopsis)
	synData := []byte(synString)
	synString = string(synData[0:len(synData)-11]) + " fileName\n"
	if opt.Called("help") || len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, synString)
		usage()
		os.Exit(0)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n\n", err)
		fmt.Fprintf(os.Stderr, synString)
		usage()
		os.Exit(0)
	}
	var baseFileName, outFileName, subErr string
	baseFileName = os.Args[len(os.Args)-1]
	revCheck, _ := substr(baseFileName,len(baseFileName)-3,3)
	if revCheck == ".sp" {
		if opMode != "fwd" {
			opMode = "rev"
		} else {
			// So you can encrypt a .sp extension file again if desired
			opMode = "fwd"
		}
	}
	if opMode == "" {
		opMode = "fwd"
	}
	switch opMode {
		case "fwd":
			outFileName = fmt.Sprintf("%s.sp", baseFileName)
		case "rev":
			outFileName, subErr = substr(baseFileName,0,-3)
			if subErr != "" {
				fmt.Println("Filename substring error:", subErr)
				os.Exit(1)
			}
			if threadCnt != 0 {
				fmt.Println("Threads cannot be set during decryption")
				fmt.Println("as thread count must be the same as it")
				fmt.Println("was encrypted with and will be looked up.")
			}
	}
	if threadCnt == 0 {
		threadCnt = threads()
	} else if threadCnt < 0 || threadCnt > 255 {
		fmt.Fprintf(os.Stderr, synString)
		usage()
		os.Exit(0)
	}
	return outFileName, opMode, threadCnt, roundCnt
}

// gets passphrase from user
func getPassphrase () string {
	var bytePassword []byte
	var checkPassword []byte
	var err error
	for {
		fmt.Print("Enter Password:")
		bytePassword, err = term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println("\nTerm read error! ", err)
		} else {
			fmt.Println()
		}
		fmt.Print("Re-enter Password:")
		checkPassword, err = term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println("\nTerm read error! ", err)
		} else {
			fmt.Println()
		}
		if string(bytePassword) == string(checkPassword) && len(checkPassword) >= 4 {
			break
		} else {
			fmt.Println("Password mismatch or shorter than 4 chars!")
		}
	}
	return string(bytePassword)
}

// prints cipher usage
func usage () {
	fmt.Println("Operating modes are fwd and rev. Defaults to fwd unless file has .sp ext.")
	fmt.Println("Thread count is optional (0-255) and defaults to 0 which uses number of CPUs.")
	fmt.Println("Rounds default to 7 but can be set to anything you want HOWEVER this value")
	fmt.Println("is currently NOT stored and MUST be specified @ encryption AND decryption time.")
	return
}

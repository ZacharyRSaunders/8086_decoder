package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// Decode asm files
	// Get the filename and decode it to a destination file.
	asm_file := os.Args[1]
	destination := os.Args[2]

	// Open File to read.
	fmt.Printf("Reading File: %s\n", asm_file)
	buffer := readFile(asm_file)

	// Create file to write and open it.
	fmt.Printf("Writing to file: %s\n", destination)

	// Decode file
	fmt.Println("Decoding...")
	n := len(buffer)
	i := 0
  instructionLen := 0
	for n > i {
		fmt.Printf("Index %d contains %b \n", i, buffer[i])
		for _, encoding := range Encodings {
      if encoding.opcode == buffer[i] >> encoding.bitShift {
        instructionLen = encoding.minLen
        fmt.Printf("Opcode: %b\n", encoding.opcode)
      }
		}
		i += instructionLen
	}

	// Close files
	fmt.Println("Cleaning Up...")
}

func getRegister(w byte, reg byte) string {
  for _, reg := range Registers {
    if reg.w == w {
      return reg.name
    }
  }
  return ""
}

func readFile(filepath string) []byte {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Failure to open file.")
	}

	buffer := make([]byte, 4096)
	bytesRead, err := file.Read(buffer)
	if err != nil {
		log.Fatalf("Failure to read from the file.")
	}
	fmt.Printf("Read %d bytes from file.\n", bytesRead)

	return buffer[:bytesRead]
   }

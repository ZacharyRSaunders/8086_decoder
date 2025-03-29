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
			if encoding.opcode == buffer[i]>>byte(encoding.identifiers[0].shift) {
				instructionLen = encoding.minLen
				identifiersDecoded, lengthMod := decodeInstruction(encoding.identifiers, buffer[i:i+8])
				// TODO use length mod to determine if more data or displacement bytes were added

				fmt.Printf("Opcode: %b, Instruction Length: %d\n", encoding.opcode, lengthMod)
				for key, value := range identifiersDecoded {
					fmt.Printf("Identifier: %s, Value: %b\n", key, value)
				}

				var suffix1, suffix2 string
				_, ok := identifiersDecoded["reg"]
				if ok {
					suffix2 = getRegister(identifiersDecoded["w"], identifiersDecoded["reg"])
				}
				_, ok = identifiersDecoded["r/m"]
				if ok {
					if identifiersDecoded["mod"] == 0b11 {
						suffix1 = getRegister(identifiersDecoded["w"], identifiersDecoded["r/m"])
					}
				}
				fmt.Printf("%s %s, %s\n", encoding.asmRep, suffix1, suffix2)
			}
		}
		i += instructionLen
	}

	// Close files
	fmt.Println("Cleaning Up...")
}

func decodeInstruction(identifiers []Identifier, buffer []byte) (map[string]byte, int) {
	// Determines instructions length and assigns values to the identifiers.
	values := map[string]byte{}
	for _, identifier := range identifiers {
		value := buffer[identifier.byteIndex] & identifier.mask >> identifier.shift
		values[identifier.name] = value
		// TODO Add other modifiers for mod which will add disp-lo/disp-hi bytes to instructions
		if identifier.name == "mod" {
			if value == 0b11 {
				fmt.Println("mod is 0b11")
			}
		}
	}
	instructionLen := 0
	return values, instructionLen
}

func getRegister(w byte, reg byte) string {
	for _, r := range Registers {
		if r.w == w {
			if r.reg == reg {
				return r.name
			}
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

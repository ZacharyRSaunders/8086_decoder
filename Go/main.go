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

	// Create temp file to write and open it.
	fmt.Printf("Writing to file: %s\n", destination)
	tempFile, err := os.CreateTemp("./", fmt.Sprintf("%s-*.asm", destination))
	if err != nil {
		log.Fatalf("Failed to create temporary file: %v", err)
	}
	fileName := tempFile.Name()
	fmt.Printf("Created temporary file: %s\n", fileName)
	defer tempFile.Close()

	// Add header text to file
	if _, err := tempFile.Write([]byte("\nbits 16\n\n")); err != nil {
		// Attempt to close before logging fatal, resource leak is less critical than crash
		tempFile.Close()
		log.Fatalf("Failed to write to temporary file %s: %v", fileName, err)
	}

	// Decode file
	fmt.Println("Decoding...")
	n := len(buffer)
	i := 0
	instructionLen := 0
	for n > i {
		fmt.Printf("Index %d contains %b \n", i, buffer[i])
		for _, encoding := range Encodings {
			if encoding.opcode == buffer[i]>>byte(encoding.identifiers[0].shift) {
				var instructionString string
				instructionString, instructionLen = getInstructionString(encoding, buffer[i:i+8])

				if _, err := tempFile.Write([]byte(instructionString)); err != nil {
					// Attempt to close before logging fatal, resource leak is less critical than crash
					tempFile.Close()
					log.Fatalf("Failed to write to temporary file %s: %v", fileName, err)
				}
			}
		}
		i += instructionLen
	}

	// Close files
	fmt.Println("Cleaning Up...")
}

func getInstructionString(encoding Encoding, buffer []byte) (string, int) {
	instructionLen := encoding.minLen
	identifiersDecoded, lengthMod := decodeInstruction(encoding.identifiers, buffer)
	// TODO use length mod to determine if more data or displacement bytes were added

	fmt.Printf("Opcode: %b, Instruction Length: %d\n", encoding.opcode, lengthMod)
	for key, value := range identifiersDecoded {
		fmt.Printf("Identifier: %s, Value: %b\n", key, value)
	}

	var suffix1, suffix2 string
	_, ok := identifiersDecoded["reg"]
	if ok {
		suffix1 = getRegister(identifiersDecoded["w"], identifiersDecoded["reg"])
	}
	_, ok = identifiersDecoded["r/m"]
	if ok {
		if identifiersDecoded["mod"] == 0b11 {
			suffix2 = getRegister(identifiersDecoded["w"], identifiersDecoded["r/m"])
		}
		if identifiersDecoded["mod"] == 0b00 {
			if identifiersDecoded["r/m"] == 0b110 {
				// Implement this
				suffix2 = "DIRECT TO ADDRESS"
			} else {
				suffix2 = fmt.Sprintf("[%s]", getEffAddr(identifiersDecoded["r/m"]))
			}
		}
		if identifiersDecoded["mod"] == 0b01 {
			suffix2 = fmt.Sprintf("[%s + %d]", getEffAddr(identifiersDecoded["r/m"]), identifiersDecoded["data"])
		}
		if identifiersDecoded["mod"] == 0b10 {
			suffix2 = fmt.Sprintf("[%s + %d]", getEffAddr(identifiersDecoded["r/m"]), (uint16(identifiersDecoded["data if w = 1"])<<8)|uint16(identifiersDecoded["data"]))
		}
	}

	var instructionString string
	if identifiersDecoded["d"] == 0 {
		instructionString = fmt.Sprintf("%s %s, %s\n", encoding.asmRep, suffix2, suffix1)
	} else {
		instructionString = fmt.Sprintf("%s %s, %s\n", encoding.asmRep, suffix1, suffix2)
	}
	return instructionString, instructionLen
}

func decodeInstruction(identifiers []Identifier, buffer []byte) (map[string]byte, int) {
	// Determines instructions length and assigns values to the identifiers.
	values := map[string]byte{}
	for _, identifier := range identifiers {
		var value byte
		if identifier.mask == 0b00000000 {
			value = buffer[identifier.byteIndex]
			values[identifier.name] = value
		} else {
			value = buffer[identifier.byteIndex] & identifier.mask >> identifier.shift
			values[identifier.name] = value
		}
		// TODO Add other modifiers for mod which will add disp-lo/disp-hi bytes to instructions
		if identifier.name == "data" {
			if values["w"] == 0b1 {
				identifiers = append(identifiers, Identifier{"data if w = 1", 0b00000000, 0, identifier.byteIndex + 1})
			}
		}
		if identifier.name == "mod" {
			if value == 0b01 {
				// 8bit displacement
				identifiers = append(identifiers, Identifier{"disp-lo", 0b00000000, 0, 3})
			}
			if value == 0b10 {
				// 16bit displacement
				identifiers = append(identifiers, Identifier{"disp-lo", 0b00000000, 0, 3})
				identifiers = append(identifiers, Identifier{"disp-hi", 0b00000000, 0, 4})
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

func getEffAddr(rm byte) string {
	for _, r := range EffectiveAddress {
		if r.rm == rm {
			return r.name
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

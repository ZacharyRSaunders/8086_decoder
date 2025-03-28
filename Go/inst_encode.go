package main

type Identifier struct {
	name      string
	mask byte
}

type Encoding struct {
	name        string
	opcode      byte
  bitShift    int
  minLen      int
  maxLen      *int
	identifiers []Identifier
	asmRep      string
}

var Encodings []Encoding = []Encoding{
	// MOV = Move
	{
    name: "Register/memory to/from register", opcode: 0b100010, bitShift: 2,
    asmRep: "mov", minLen: 2,
		identifiers: []Identifier{
      {"d", 0b00000010},
      {"w", 0b00000001},
			{"mod", 0b11000000},
			{"reg", 0b00111000},
			{"r/m", 0b00000111},
      // Only add disp-lo/disp-hi if the mod calls for it
		},
	},
	{
    name: "Immediate to register/memory", opcode: 0b11000, bitShift: 3,
    asmRep: "mov", minLen: 3,
		identifiers: []Identifier{
			{"?", 0b00000100},
			{"d", 0b00000010},
			{"w", 0b00000001},
			{"mod", 0b11000000},
			{"reg", 0b00111000},
			{"r/m", 0b00000111},
			{"data", 0b00000000},
      // only add data if w=1 if data exists and w=1
		},
	},
}

type Register struct {
  w byte
  reg byte
  name string
}

var Registers []Register = []Register{
  {w: 0b0, reg: 0b000, name: "al"},
  {w: 0b0, reg: 0b001, name: "cl"},
  {w: 0b0, reg: 0b010, name: "dl"},
  {w: 0b0, reg: 0b011, name: "bl"},
  {w: 0b0, reg: 0b100, name: "ah"},
  {w: 0b0, reg: 0b101, name: "ch"},
  {w: 0b0, reg: 0b110, name: "dh"},
  {w: 0b0, reg: 0b111, name: "bh"},
  {w: 0b1, reg: 0b000, name: "ax"},
  {w: 0b1, reg: 0b001, name: "cx"},
  {w: 0b1, reg: 0b010, name: "dx"},
  {w: 0b1, reg: 0b011, name: "bx"},
  {w: 0b1, reg: 0b100, name: "sp"},
  {w: 0b1, reg: 0b101, name: "bp"},
  {w: 0b1, reg: 0b110, name: "si"},
  {w: 0b1, reg: 0b111, name: "di"},
}




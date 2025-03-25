package main

type Identifier struct {
	name      string
	mask byte
}

type Encoding struct {
	name        string
	opcode      byte
	identifiers []Identifier
	asmRep      string
}

var Encodings []Encoding = []Encoding{
	// MOV = Move
	{
		name: "Register/memory to/from register", opcode: 0b100010, asmRep: "mov",
		identifiers: []Identifier{
			{"opcode", 0b11111100},
			{"d", 0b00000010},
			{"w", 0b00000001},
			{"mod", 0b11000000},
			{"reg", 0b00111000},
			{"r/m", 0b00000111},
			{"DISP-LO", 0b111111111},
			{"DISP-HI", 0b111111111},
		},
	},
	{
		name: "Immediate to register/memory", opcode: 0b11000, asmRep: "mov",
		identifiers: []Identifier{
			{"opcode", 0b11111000},
			{"?", 0b00000100},
			{"d", 0b00000010},
			{"w", 0b00000001},
			{"mod", 0b11000000},
			{"reg", 0b00111000},
			{"r/m", 0b00000111},
			{"DISP-LO", 0b111111111},
			{"DISP-HI", 0b111111111},
			{"data", 0b111111111},
			{"data if w=1", 0b111111111},
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




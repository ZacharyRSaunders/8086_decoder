package main

type Identifier struct {
	name      string
	mask      byte
	shift     int
	byteIndex int
}

type Encoding struct {
	name        string
	opcode      byte
	minLen      int
	maxLen      *int
	identifiers []Identifier
	asmRep      string
}

var Encodings []Encoding = []Encoding{
	// MOV = Move
	{
		name: "Register/memory to/from register", opcode: 0b100010,
		asmRep: "mov", minLen: 2,
		identifiers: []Identifier{
			{"opcode", 0b11111100, 2, 0},
			{"d", 0b00000010, 1, 0},
			{"w", 0b00000001, 0, 0},
			{"mod", 0b11000000, 6, 1},
			{"reg", 0b00111000, 3, 1},
			{"r/m", 0b00000111, 0, 1},
		},
	},
	{
		name: "Immediate to register/memory", opcode: 0b11000,
		asmRep: "mov", minLen: 3,
		identifiers: []Identifier{
			{"opcode", 0b11111000, 3, 0},
			{"?", 0b00000100, 2, 0},
			{"d", 0b00000010, 1, 0},
			{"w", 0b00000001, 0, 0},
			{"mod", 0b11000000, 6, 1},
			{"reg", 0b00111000, 3, 1},
			{"r/m", 0b00000111, 0, 1},
			// only add data if w=1 if data exists and w=1
		},
	},
	{
		name: "Immediate to register", opcode: 0b1011,
		asmRep: "mov", minLen: 2,
		identifiers: []Identifier{
			{"opcode", 0b11110000, 4, 0},
			{"w", 0b00001000, 3, 0},
			{"reg", 0b00000111, 0, 0},
			{"data", 0b00000000, 0, 1},
		},
	},
}

type Register struct {
	w    byte
	reg  byte
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

type EffectiveAddr struct {
	rm   byte
	name string
}

var EffectiveAddress []EffectiveAddr = []EffectiveAddr{
	{rm: 0b000, name: "bx + si"},
	{rm: 0b001, name: "bx + di"},
	{rm: 0b010, name: "bp + si"},
	{rm: 0b011, name: "bp + di"},
	{rm: 0b100, name: "si"},
	{rm: 0b101, name: "di"},
	{rm: 0b110, name: "bp"},
	{rm: 0b111, name: "bx"},
}

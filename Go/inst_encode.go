package main

type Identifier struct {
	name      string
	bitLength int
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
			{"opcode", 6},
			{"d", 1},
			{"w", 1},
			{"mod", 2},
			{"reg", 3},
			{"r/m", 3},
		},
	},
	{
		name: "Immediate to register/memory", opcode: 0b11000, asmRep: "mov",
		identifiers: []Identifier{
			{"opcode", 5},
			{"?", 1},
			{"d", 1},
			{"w", 1},
			{"mod", 2},
			{"reg", 3},
			{"r/m", 3},
			{"DISP-LO", 8},
			{"DISP-HI", 8},
			{"data", 8},
			{"data if w=1", 8},
		},
	},
}

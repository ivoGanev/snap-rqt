package internal

const (
	CREATE = iota
	READ
	UPDATE
	DELETE
)

type CrudOp int

var CrudOperations = []CrudOp{CREATE, READ, UPDATE, DELETE}

const (
	ONE = iota
	MANY
)

type Multiplicity int

var Multiplicities = []Multiplicity{ONE, MANY}
package models

type DbItem struct {
	Lseq  string
	Key   string
	Value string
}

type ValidateItem struct {
	Lseq          *string
	LseqItemValid string
	Hash          string
}

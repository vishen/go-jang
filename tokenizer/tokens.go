package tokenizer

const (
	OpenTag TokenType = iota
	CloseTag
	Attribute
	Value
	Text
	Assign
	ForwardSlash
	Tag
)

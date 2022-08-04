package common

type Enum interface {
	name() string
	ordinal() int
	values() *[]string
}

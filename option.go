package collapsar

import "collapsar/policy"

type Option struct {
	Length     int
	Calculator HashInterface
	FailHandler FailHandlerFunc
	RemoveHandler RemoveHandlerFunc
	EliminateHandler policy.EliminateInterface
}

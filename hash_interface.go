package collapsar

type HashInterface interface {
	Hash(key string) uint32
}

package conf

type Block interface {
	Call(name string, args ...string) error
	CallBlock(name string, args ...string) (Block, error)
}

type BlockWrapper struct {
	Parent  *BlockWrapper
	Current Block
}

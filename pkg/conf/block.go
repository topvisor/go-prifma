package conf

type Block interface {
	Call(command Command) error
	CallBlock(command Command) (Block, error)
}

type BlockWrapper struct {
	Parent  *BlockWrapper
	Current Block
}

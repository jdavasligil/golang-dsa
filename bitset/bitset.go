package bitset

type BitSet8 uint8

func (b *BitSet8) Set(idx uint8) {
	*b |= (1 << idx)
}
func (b *BitSet8) Unset(idx uint8) {
	*b &= ^(1 << idx)
}
func (b *BitSet8) Get(idx uint8) bool {
	return (*b & (1 << idx)) != 0
}

type BitSet16 uint16

func (b *BitSet16) Set(idx uint8) {
	*b |= (1 << idx)
}
func (b *BitSet16) Unset(idx uint8) {
	*b &= ^(1 << idx)
}
func (b *BitSet16) Get(idx uint8) bool {
	return (*b & (1 << idx)) != 0
}

type BitSet32 uint32

func (b *BitSet32) Set(idx uint8) {
	*b |= (1 << idx)
}
func (b *BitSet32) Unset(idx uint8) {
	*b &= ^(1 << idx)
}
func (b *BitSet32) Get(idx uint8) bool {
	return (*b & (1 << idx)) != 0
}

type BitSet64 uint64

func (b *BitSet64) Set(idx uint8) {
	*b |= (1 << idx)
}
func (b *BitSet64) Unset(idx uint8) {
	*b &= ^(1 << idx)
}
func (b *BitSet64) Get(idx uint8) bool {
	return (*b & (1 << idx)) != 0
}

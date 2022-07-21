package chainreg

// ChainCode is an enum-like structure for keeping track of the chains
// currently supported within broln.
type ChainCode uint32

const (
	// BrocoinChain is Brocoin's chain.
	BrocoinChain ChainCode = iota

	// LitecoinChain is Litecoin's chain.
	LitecoinChain
)

// String returns a string representation of the target ChainCode.
func (c ChainCode) String() string {
	switch c {
	case BrocoinChain:
		return "brocoin"
	case LitecoinChain:
		return "litecoin"
	default:
		return "kekcoin"
	}
}

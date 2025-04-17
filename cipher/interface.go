package cipher

type Cipher interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}

const (
	AlgoAesCbc    = "CAFBCBAD-B6E7-4CAB-8A67-14D39F00CE1E"
	AlgoAesEcb    = "A474B1C2-3DE0-4EA2-8C5F-7093409CE6C4"
	AlgoDesEdeCbc = "5BFBA864-BBA9-42DB-8EAD-49B5F412BD81"
	AlgoDesEdeEcb = "6E0B65FF-0B5B-459C-8FCE-EC7F2BEA9FF5"
	AlgoZUC       = "B809531F-0007-4B5B-923B-4BD560398113"
	AlgoSm4Cbc    = "F3974434-C0DD-4C20-9E87-DDB6814A1C48"
	AlgoSm4Ecb    = "ED382482-F72C-4C41-A76D-28EEA0F1F2AF"
	AlgoXTea      = "B3047D4E-67DF-4864-A6A5-DF9B9E525C79"
	AlgoXTeaIv    = "C32C68F9-CA81-4260-A329-BBAFD1A9CCD1"
)

func NewCipher(algoID string) Cipher {
	switch algoID {
	case AlgoAesCbc:
		return new(AesCbc)
	case AlgoAesEcb:
		return new(AesEcb)
	case AlgoDesEdeCbc:
		return new(DesEdeCbc)
	case AlgoDesEdeEcb:
		return new(DesEdeEcb)
	case AlgoZUC:
		return new(Zuc)
	case AlgoSm4Cbc:
		return new(Sm4Cbc)
	case AlgoSm4Ecb:
		return new(Sm4Ecb)
	case AlgoXTea:
		return new(XTea)
	case AlgoXTeaIv:
		return new(XTeaIv)
	default:
		return nil
	}
}

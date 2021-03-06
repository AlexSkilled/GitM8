package emoji

import (
	"bytes"
)

const (
	Loudspeaker = 0x1F4E2

	BlackLargeSquare  = 0x2B1B
	YellowLargeSquare = 0x1F7E8
	RedLargeSquare    = 0x1F7E5
	GreenLargeSquare  = 0x1F7E9

	GrayLargeSquare = "ā¬"
	Man             = "šš¼\u200dāļø"
	Branches        = "š"
	EyeWatch        = "šāšØ"
	New             = "š"
	Gear            = "ā"

	WhiteLargeCircle  = 0x26AA
	OrangeLargeCircle = 0x1F7E0
	BlueLargeCircle   = 0x1F535
	BlackLargeCircle  = 0x26AB
	GreenLargeCircle  = 0x1F7E2

	WhiteCheckMark = 0x2705
)

func GetEmoji(r rune) string {
	buff := bytes.Buffer{}
	buff.WriteRune(r)
	return buff.String()
}

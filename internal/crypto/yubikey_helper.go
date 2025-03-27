// File: internal/crypto/yubikey_helper.go

package crypto

import (
	"github.com/billchurch/pica/internal/yubikey"
)

// FromYubiKeySlot converts a yubikey.PIVSlot to a crypto.Slot
func FromYubiKeySlot(pivSlot yubikey.PIVSlot) Slot {
	return Slot(pivSlot)
}

// ToYubiKeySlot converts a crypto.Slot to a yubikey.PIVSlot
func ToYubiKeySlot(slot Slot) yubikey.PIVSlot {
	return yubikey.PIVSlot(slot)
}

//go:build !amd64

package hashtable

func matchMask(a uint8, b [16]uint8) uint16 {
	return matchMaskLoop(a, b)
}

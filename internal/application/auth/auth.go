// Package auth provides the implementation of the authentication layer.
package auth

func ZeroBytes(bts []byte) {
	for i := range bts {
		bts[i] = 0
	}
}

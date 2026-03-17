// Package auth provides the implementation of the authentication layer.
package auth

func ZeroBytes(bts []byte) {
	if len(bts) == 0 {
		return
	}

	for i := range bts {
		bts[i] = 0
	}
}

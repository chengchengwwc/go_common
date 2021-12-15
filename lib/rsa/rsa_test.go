package rsa

import "testing"

func TestNewRsaServer(t *testing.T) {
	mm, err := NewRsaServer("ca.pem", "ca-key.pem")
	if err != nil {
		t.Log(err)
		return
	}
	tmpString := "abcd"

	ss, err := mm.Encrypt(tmpString)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(ss)
	vv, err := mm.Decryption(ss)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(vv)

}

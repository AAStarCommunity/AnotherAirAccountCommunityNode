package signature

import (
	"encoding/base64"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/herumi/bls-eth-go-binary/bls"
)

func TestRandSplit(t *testing.T) {
	data := []byte("asdfasdfasdfasdf314")
	n := rand.IntN(3) + 1
	groups := randSplit(data, n)
	for _, g := range groups {
		fmt.Print(string(g))
	}
}

func itorSigners(arr []string, k int) [][]string {
	var helper func(start int, combo []string)
	res := [][]string{}
	combo := []string{}

	helper = func(start int, combo []string) {
		if len(combo) == k {
			// Make a copy of combo since combo will be reused
			c := make([]string, len(combo))
			copy(c, combo)
			res = append(res, c)
			return
		}
		for i := start; i <= len(arr)-(k-len(combo)); i++ {
			// Add current element
			combo = append(combo, arr[i])
			// Move to the next element
			helper(i+1, combo)
			// Backtrack to try next candidate
			combo = combo[:len(combo)-1]
		}
	}

	helper(0, combo)
	return res
}

func TestBls(t *testing.T) {
	threshold := 2
	totalSigners := 5
	val := "dfabcasdfasf"
	data := []byte(val)
	verifyData := []byte(val)

	sig, msk, err := BlsTss(threshold, totalSigners, data)

	if err != nil {
		t.Error(err)
	}

	pub := bls.PublicKey{}
	pub.Deserialize(msk)
	signValidator, err := RecoverSignerGroup(threshold, &pub, totalSigners)
	if err != nil {
		t.Error(err)
	}

	sigObj := &bls.Sign{}
	sigObj.Deserialize(sig)

	if !signValidator.Verify(sigObj, verifyData) {
		t.Error("Signature verification failed")
	}
}

func TestSign(t *testing.T) {
	allId := []string{"1", "2", "3", "4", "5"}

	threshold := len(allId) - 2

	msg := "abc"
	oirgVal := []byte(msg)
	verifyVal := oirgVal

	grp, err := NewSignerGroup(threshold, allId...)

	if err != nil {
		t.Error(err)
	}

	comb := itorSigners(allId, threshold)

	signValidator, err := RecoverSignerGroup(threshold, grp.GetPublicKeys(), len(allId))
	if err != nil {
		t.Error(err)
	}

	for _, c := range comb {
		subGrp, err := grp.PickUpSigners(c...)
		if err != nil {
			t.Error(err)
		}
		sig, err := subGrp.Sign(oirgVal)
		if err != nil {
			t.Error(err)
		}
		s := sig.Serialize()
		dsig := &bls.Sign{}
		dsig.Deserialize(s)

		if !signValidator.Verify(dsig, verifyVal) {
			t.Error("Signature verification failed")
		}

		tmp := append(verifyVal, oirgVal...)
		if signValidator.Verify(dsig, tmp) {
			t.Error("Signature verification failed")
		}

		pubkey := grp.GetPublicKeys().Serialize()
		m1 := base64.StdEncoding.EncodeToString(pubkey)

		msg := base64.StdEncoding.EncodeToString([]byte(msg))

		sigMsg := base64.StdEncoding.EncodeToString(s)

		fmt.Printf("pubkey: %s\n msg: %s\n sig: %s\n\n\n", m1, msg, sigMsg)
		break
	}
}

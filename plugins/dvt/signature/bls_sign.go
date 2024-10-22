package signature

import (
	"crypto/sha256"
	"fmt"
	"math/rand/v2"
)

func randSplit(data []byte, n int) [][]byte {
	lengths := make([]int, n)
	total := len(data)

	// 生成随机分组方案
	for i := 0; i < n-1; i++ {
		lengths[i] = rand.IntN(total-(n-i-1)) + 1
		total -= lengths[i]
	}
	lengths[n-1] = total

	groups := make([][]byte, n)
	start := 0
	for i, length := range lengths {
		groups[i] = data[start : start+length]
		start += length
	}

	return groups
}

// Bls sign data using BLS signature scheme
func Bls(data []byte) (blsSignature []byte, blsPublickey []byte, err error) {
	msgHash := sha256.Sum256(data)
	groups := randSplit(msgHash[:], rand.IntN(3)+1)

	for _, g := range groups {
		// TODO: request DVT webapi to sign data
		fmt.Print(string(g))
	}

	panic("not implemented")
}

// BlsTss sign data using BLS threshold signature scheme
func BlsTss(threshold, totalSigners int, data []byte) (blsSignature []byte, blsPublickey []byte, err error) {
	allId := make([]string, totalSigners)
	for i := 0; i < totalSigners; i++ {
		allId[i] = fmt.Sprint(i)
	}
	grp, err := NewSignerGroup(threshold, allId...)
	if err != nil {
		return nil, nil, err
	}

	subGrp, err := grp.PickUpSigners(allId...)
	if err != nil {
		return nil, nil, err
	}

	sig, err := subGrp.Sign(data)
	if err != nil {
		return nil, nil, err
	}

	blsSignature = sig.Serialize()
	blsPublickey = grp.GetPublicKeys().Serialize()
	return
}

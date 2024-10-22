package signature

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"sync"
)

func randSplit(data string, n int) []string {
	lengths := make([]int, n)
	total := len(data)

	// 生成随机分组方案
	for i := 0; i < n-1; i++ {
		lengths[i] = rand.IntN(total-(n-i-1)) + 1
		total -= lengths[i]
	}
	lengths[n-1] = total

	groups := make([]string, n)
	start := 0
	for i, length := range lengths {
		groups[i] = data[start : start+length]
		start += length
	}

	return groups
}

// Bls sign data using BLS signature scheme
func Bls(data []byte) (blsSignature []byte, blsPublickey []byte, err error) {
	dvtNodes := []string{
		"http://localhost:8081",
		// "http://localhost:8082",
		// "http://localhost:8083",
	}
	msgHash := fmt.Sprintf("%x", sha256.Sum256(data))[0:31]
	groups := randSplit(msgHash, len(dvtNodes))
	mapGroups := make(map[string]string)
	for i, g := range groups {
		mapGroups[dvtNodes[i]] = g
	}

	var wg sync.WaitGroup
	wg.Add(len(groups))
	for dvt, g := range mapGroups {
		go func(d string, group string) {
			defer wg.Done()
			body := struct {
				Domain  string `json:"domain"`
				Message string `json:"message"`
			}{
				Domain:  "dvt",
				Message: group,
			}
			jsonData, err := json.Marshal(body)
			if err != nil {
				fmt.Println("Error encoding JSON:", err)
				return
			}

			resp, err := http.Post(d+"/sign", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == 200 {
				fmt.Println("response Status:", resp.Status)
				return
			}
		}(dvt, g)
	}
	wg.Wait()
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

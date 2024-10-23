package signature

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"sync"
)

type signResponse struct {
	Signature []string `json:"sig"`
	PublicKey []string `json:"pubkeys"`
}

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

func requestSign(d string, group *string, domain *string) (*signResponse, error) {
	body := struct {
		Domain  *string `json:"domain"`
		Message *string `json:"message"`
	}{
		Domain:  domain,
		Message: group,
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return nil, err
	}

	resp, err := http.Post(d+"/sign", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var signResponse signResponse
	if err := json.Unmarshal(respBody, &signResponse); err != nil {
		return nil, err
	}
	return &signResponse, nil
}

func aggrSign(dvt string, signGroup map[string]*signResponse) (*[2]string, error) {
	var sigs [][2]string
	for _, sign := range signGroup {
		sigs = append(sigs, [2]string{sign.Signature[0], sign.Signature[1]})
	}
	body := struct {
		Signatures [][2]string `json:"sigs"`
	}{
		Signatures: sigs,
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return nil, err
	}

	resp, err := http.Post(dvt+"/aggr", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var aggrResp struct {
		Signature [2]string `json:"sig"`
	}

	if err := json.Unmarshal(respBody, &aggrResp); err != nil {
		return nil, err
	}

	ret := aggrResp.Signature
	return &ret, nil
}

func verifyAggr(dvt string, pubkeys [][4]string, aggrSigs *[2]string, messages []string, domain *string) (bool, error) {
	payload := struct {
		PublicKey    [][4]string `json:"pubkeys"`
		AggregateSig [2]string   `json:"aggrSig"`
		Messages     []string    `json:"messages"`
		Domain       *string     `json:"domain"`
	}{
		PublicKey:    pubkeys,
		AggregateSig: *aggrSigs,
		Messages:     messages,
		Domain:       domain,
	}

	if jsonData, err := json.Marshal(payload); err != nil {
		return false, err
	} else {
		resp, err := http.Post(dvt+"/aggr/verify", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return false, err
		}
		return resp.StatusCode == http.StatusAccepted, nil
	}
}

// Bls sign data using BLS signature scheme
func Bls(data []byte) (ok bool, err error) {
	dvtNodes := []string{
		"http://localhost:8081",
		// "http://localhost:8082",
		// "http://localhost:8083",
	}
	msgHash := fmt.Sprintf("%x", sha256.Sum256(data))[0:31]
	groups := randSplit(msgHash, len(dvtNodes))
	mapMsgGroups := make(map[string]string)
	mapSignatures := make(map[string]*signResponse)
	for i, g := range groups {
		mapMsgGroups[dvtNodes[i]] = g
	}
	domain := groups[0]

	var wg sync.WaitGroup
	wg.Add(len(groups))
	for dvt, g := range mapMsgGroups {
		go func(d string, group string) {
			defer wg.Done()
			if signResult, err := requestSign(d, &group, &domain); err == nil {
				mapSignatures[d] = signResult
			}
		}(dvt, g)
	}
	wg.Wait()

	if len(mapSignatures) != len(dvtNodes) {
		return false, fmt.Errorf("not all nodes signed")
	} else {
		if aggrSign, err := aggrSign(dvtNodes[0], mapSignatures); err != nil {
			return false, err
		} else {
			pubkeys := make([][4]string, 0)
			for _, sign := range mapSignatures {
				pubkeys = append(pubkeys, [4]string{sign.PublicKey[0], sign.PublicKey[1], sign.PublicKey[2], sign.PublicKey[3]})
			}

			if ok, err := verifyAggr(dvtNodes[0], pubkeys, aggrSign, groups, &domain); err != nil {
				return false, err
			} else {
				return ok, nil
			}
		}
	}
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

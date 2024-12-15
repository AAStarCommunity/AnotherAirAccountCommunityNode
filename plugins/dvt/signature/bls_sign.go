package signature

import (
	dvtSeedworks "another_node/plugins/dvt/seedworks"
	"another_node/plugins/passkey_relay_party/seedworks"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
)

type signResponse struct {
	Signature []string `json:"sig"`
	PublicKey []string `json:"pubkeys"`
	Message   []string `json:"msg"`
}

type signGroup map[string]*signResponse

func (s signGroup) first() string {
	for k := range s {
		return k
	}
	return ""
}

func requestSign(host string, message []byte, passkeyPubkey []byte, passkey *protocol.ParsedCredentialAssertionData) (*signResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	body := struct {
		Message       string                                  `json:"message"`
		PasskeyPubkey []byte                                  `json:"passkeyPubkey"`
		Passkey       *protocol.ParsedCredentialAssertionData `json:"passkey"`
	}{
		Message:       string(message),
		PasskeyPubkey: passkeyPubkey,
		Passkey:       passkey,
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request body: %w", err)
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest("POST", host+"/sign", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server returned non-200 status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var signResponse signResponse
	if err := json.Unmarshal(respBody, &signResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &signResponse, nil
}

func aggrSign(host string, signGroup signGroup) ([]string, [][4]string, error) {
	var sigs [][2]string
	pubkeys := make([][4]string, 0)
	for _, sign := range signGroup {
		sigs = append(sigs, [2]string{sign.Signature[0], sign.Signature[1]})
		pubkeys = append(pubkeys, [4]string{sign.PublicKey[0], sign.PublicKey[1], sign.PublicKey[2], sign.PublicKey[3]})
	}
	body := struct {
		Signatures [][2]string `json:"sigs"`
		Pubkeys    [][4]string `json:"pubkeys"`
	}{
		Signatures: sigs,
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return nil, nil, err
	}

	resp, err := http.Post(host+"/aggr", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return nil, nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	var aggrResp struct {
		Signature []string `json:"sig"`
	}

	if err := json.Unmarshal(respBody, &aggrResp); err != nil {
		return nil, nil, err
	}

	ret := aggrResp.Signature
	return ret, pubkeys, nil
}

//lint:ignore U1000 ignore unused
func verifyAggr(host string, pubkeys [][4]string, aggrSigs *[2]string, messages []string, domain *string) (bool, error) {
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
		resp, err := http.Post(host+"/aggr/verify", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return false, err
		}
		return resp.StatusCode == http.StatusAccepted, nil
	}
}

// uniqueNodes removes duplicate nodes from the input slice
func uniqueNodes(nodes []string) []string {
	uniqueMap := make(map[string]struct{})
	var result []string
	for _, node := range nodes {
		if _, exists := uniqueMap[node]; !exists {
			uniqueMap[node] = struct{}{}
			result = append(result, node)
		}
	}
	return result
}

// Bls sign data using BLS signature scheme
func Bls(data []byte, threshold, timeoutSeconds int, dvtNodes []string, passkeyCA *protocol.ParsedCredentialAssertionData, passkeyCAPubKey []byte) (*seedworks.DvtResult, error) {
	dvtNodes = uniqueNodes(dvtNodes)
	if len(dvtNodes) < threshold {
		return nil, dvtSeedworks.ErrNotEnoughSigners{}
	}
	mapSignatures := make(signGroup)
	var mu sync.Mutex
	var doneOnce sync.Once
	done := make(chan struct{})

	var messagePt [2]string
	for _, host := range dvtNodes {
		go func() {
			if signResult, err := requestSign(host, data, passkeyCAPubKey, passkeyCA); err == nil {
				mu.Lock()
				mapSignatures[host] = signResult
				sigCount := len(mapSignatures)
				mu.Unlock()

				if sigCount >= threshold {
					doneOnce.Do(func() {
						messagePt[0] = signResult.Message[0]
						messagePt[1] = signResult.Message[1]
						close(done)
					})
				}
			} else {
				fmt.Println(err)
			}
		}()
	}

	timeout := time.After(time.Duration(timeoutSeconds) * time.Second)

	select {
	case <-done:
		firstNode := mapSignatures.first()
		if aggr, pubkeys, err := aggrSign(firstNode, mapSignatures); err != nil {
			return nil, err
		} else {
			return &seedworks.DvtResult{
				Signatures: aggr,
				MessagePt:  messagePt,
				PublicKeys: pubkeys,
			}, nil
		}
	case <-timeout:
		mu.Lock()
		sigCount := len(mapSignatures)
		mu.Unlock()
		if sigCount >= threshold {
			firstNode := mapSignatures.first()
			if aggr, pubkeys, err := aggrSign(firstNode, mapSignatures); err != nil {
				return nil, err
			} else {
				return &seedworks.DvtResult{
					Signatures: aggr,
					MessagePt:  messagePt,
					PublicKeys: pubkeys,
				}, nil
			}
		}
		return nil, dvtSeedworks.ErrNotEnoughSigners{}
	}
}

// BlsTss sign data using BLS threshold signature scheme
// [deprecated]
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

package bls_tss

import "github.com/herumi/bls-eth-go-binary/bls"

func init() {
	bls.Init(bls.BLS12_381)
	bls.SetETHmode(bls.EthModeDraft07)
}

type Signer struct {
	Id        bls.ID
	secretKey bls.SecretKey
	publicKey bls.PublicKey
}

func newSigner(id *bls.ID, secretKey *bls.SecretKey) *Signer {
	return &Signer{
		Id:        *id,
		secretKey: *secretKey,
		publicKey: *secretKey.GetPublicKey(),
	}
}

type SignerGroup struct {
	Signers   []Signer
	mpk       *bls.PublicKey
	threshold int
}

func NewSignerGroup(threshold int, id ...string) (*SignerGroup, error) {
	total := len(id)
	if threshold <= 0 || total <= 0 {
		return nil, &ErrThresholdGreaterThanZero{}
	}
	if threshold > total {
		return nil, &ErrThresholdGreaterThanTotal{}
	}

	var msk []bls.SecretKey = make([]bls.SecretKey, threshold)
	for i := 0; i < threshold; i++ {
		msk[i].SetByCSPRNG()
	}

	var secs []bls.SecretKey = make([]bls.SecretKey, total)

	blsId := make([]bls.ID, total)
	for i := 0; i < total; i++ {
		blsId[i].SetLittleEndian([]byte(id[i]))
	}

	// share secret key
	for i := 0; i < total; i++ {
		secs[i].Set(msk, &blsId[i])
	}

	signers := make([]Signer, total)
	for i := 0; i < total; i++ {
		signers[i] = *newSigner(&blsId[i], &secs[i])
	}

	// get master public key
	mpk := msk[0].GetPublicKey()

	return &SignerGroup{
		Signers:   signers,
		mpk:       mpk,
		threshold: threshold,
	}, nil
}

func (sg *SignerGroup) PickUpSigners(ids ...string) (*SignerGroup, error) {
	if len(ids) < sg.threshold {
		return nil, &ErrNotEnoughSigners{}
	}

	picked := make([]Signer, 0)
	for i := 0; i < len(sg.Signers); i++ {
		for j := 0; j < len(ids); j++ {
			tmpId := bls.ID{}
			tmpId.SetLittleEndian([]byte(ids[j]))
			if sg.Signers[i].Id.IsEqual(&tmpId) {
				picked = append(picked, sg.Signers[i])
			}
		}
	}

	if len(picked) < sg.threshold {
		return nil, &ErrNotEnoughSigners{}
	}

	return &SignerGroup{
		Signers:   picked,
		mpk:       sg.mpk,
		threshold: sg.threshold,
	}, nil
}

// Sign signs the message with all the signers in the (sub)SignerGroup
func (sg *SignerGroup) Sign(msg []byte) (*bls.Sign, error) {
	if len(sg.Signers) < sg.threshold {
		return nil, &ErrNotEnoughSigners{}
	}

	sig := make([]bls.Sign, 0)
	ids := make([]bls.ID, 0)
	m := string(msg)
	for i := 0; i < len(sg.Signers); i++ {
		sig = append(sig, *sg.Signers[i].secretKey.Sign(m))
		ids = append(ids, sg.Signers[i].Id)
	}

	var s bls.Sign
	if err := s.Recover(sig, ids); err != nil {
		return nil, err
	}

	return &s, nil
}

func (sg *SignerGroup) Verify(sig *bls.Sign, msg []byte) bool {
	return sig.Verify(sg.mpk, string(msg))
}

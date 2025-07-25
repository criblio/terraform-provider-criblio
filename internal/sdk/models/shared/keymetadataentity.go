// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
	"github.com/criblio/terraform-provider-criblio/internal/sdk/internal/utils"
)

type EncryptionAlgorithm string

const (
	EncryptionAlgorithmAes256Cbc EncryptionAlgorithm = "aes-256-cbc"
	EncryptionAlgorithmAes256Gcm EncryptionAlgorithm = "aes-256-gcm"
)

func (e EncryptionAlgorithm) ToPointer() *EncryptionAlgorithm {
	return &e
}
func (e *EncryptionAlgorithm) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "aes-256-cbc":
		fallthrough
	case "aes-256-gcm":
		*e = EncryptionAlgorithm(v)
		return nil
	default:
		return fmt.Errorf("invalid value for EncryptionAlgorithm: %v", v)
	}
}

type KMSForThisKey string

const (
	KMSForThisKeyLocal KMSForThisKey = "local"
)

func (e KMSForThisKey) ToPointer() *KMSForThisKey {
	return &e
}
func (e *KMSForThisKey) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "local":
		*e = KMSForThisKey(v)
		return nil
	default:
		return fmt.Errorf("invalid value for KMSForThisKey: %v", v)
	}
}

// InitializationVectorSize - Length of the initialization vector, in bytes
type InitializationVectorSize int64

const (
	InitializationVectorSizeTwelve   InitializationVectorSize = 12
	InitializationVectorSizeThirteen InitializationVectorSize = 13
	InitializationVectorSizeFourteen InitializationVectorSize = 14
	InitializationVectorSizeFifteen  InitializationVectorSize = 15
	InitializationVectorSizeSixteen  InitializationVectorSize = 16
)

func (e InitializationVectorSize) ToPointer() *InitializationVectorSize {
	return &e
}
func (e *InitializationVectorSize) UnmarshalJSON(data []byte) error {
	var v int64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case 12:
		fallthrough
	case 13:
		fallthrough
	case 14:
		fallthrough
	case 15:
		fallthrough
	case 16:
		*e = InitializationVectorSize(v)
		return nil
	default:
		return fmt.Errorf("invalid value for InitializationVectorSize: %v", v)
	}
}

type KeyMetadataEntity struct {
	KeyID       string               `json:"keyId"`
	Description *string              `json:"description,omitempty"`
	Algorithm   *EncryptionAlgorithm `default:"aes-256-cbc" json:"algorithm"`
	Kms         *KMSForThisKey       `default:"local" json:"kms"`
	Keyclass    *float64             `default:"0" json:"keyclass"`
	Created     *float64             `json:"created,omitempty"`
	Expires     *float64             `json:"expires,omitempty"`
	PlainKey    *string              `json:"plainKey,omitempty"`
	CipherKey   *string              `json:"cipherKey,omitempty"`
	// Seed encryption with a [nonce](https://en.wikipedia.org/wiki/Cryptographic_nonce) to make the key more random and unique. Must be enabled with the aes-256-gcm algorithm.
	UseIV *bool `default:"false" json:"useIV"`
	// Length of the initialization vector, in bytes
	IvSize *InitializationVectorSize `default:"12" json:"ivSize"`
	// Name of the Group/Fleet that created this key
	Group *string `json:"group,omitempty"`
}

func (k KeyMetadataEntity) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(k, "", false)
}

func (k *KeyMetadataEntity) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &k, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *KeyMetadataEntity) GetKeyID() string {
	if o == nil {
		return ""
	}
	return o.KeyID
}

func (o *KeyMetadataEntity) GetDescription() *string {
	if o == nil {
		return nil
	}
	return o.Description
}

func (o *KeyMetadataEntity) GetAlgorithm() *EncryptionAlgorithm {
	if o == nil {
		return nil
	}
	return o.Algorithm
}

func (o *KeyMetadataEntity) GetKms() *KMSForThisKey {
	if o == nil {
		return nil
	}
	return o.Kms
}

func (o *KeyMetadataEntity) GetKeyclass() *float64 {
	if o == nil {
		return nil
	}
	return o.Keyclass
}

func (o *KeyMetadataEntity) GetCreated() *float64 {
	if o == nil {
		return nil
	}
	return o.Created
}

func (o *KeyMetadataEntity) GetExpires() *float64 {
	if o == nil {
		return nil
	}
	return o.Expires
}

func (o *KeyMetadataEntity) GetPlainKey() *string {
	if o == nil {
		return nil
	}
	return o.PlainKey
}

func (o *KeyMetadataEntity) GetCipherKey() *string {
	if o == nil {
		return nil
	}
	return o.CipherKey
}

func (o *KeyMetadataEntity) GetUseIV() *bool {
	if o == nil {
		return nil
	}
	return o.UseIV
}

func (o *KeyMetadataEntity) GetIvSize() *InitializationVectorSize {
	if o == nil {
		return nil
	}
	return o.IvSize
}

func (o *KeyMetadataEntity) GetGroup() *string {
	if o == nil {
		return nil
	}
	return o.Group
}

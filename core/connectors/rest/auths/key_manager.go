package auths

import "errors"

type KeyManager interface {
	Latest() (version string, key []byte, err error)
	Get(version string) ([]byte, error)
}

type StaticKeyManager struct {
	keys          map[string][]byte
	latestVersion string
}

func NewStaticKeyManager(keys map[string][]byte,latest string) *StaticKeyManager {

	return &StaticKeyManager{
		keys:          keys,
		latestVersion: latest,
	}
}

func (k *StaticKeyManager) Latest() (string,[]byte,error) {

	key, ok := k.keys[k.latestVersion]
	if !ok {
		return "", nil, errors.New("latest key not found")
	}

	return k.latestVersion, key, nil
}

func (k *StaticKeyManager) Get(version string) ([]byte, error) {

	key, ok := k.keys[version]
	if !ok {
		return nil, errors.New("key version not found")
	}

	return key, nil
}
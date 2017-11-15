package jorm

import (
	"github.com/setekhid/jormungand/pkg/jorm/proto"
)

type Key struct {
	ID     string
	Cipher string
}

func (k *Key) Encoder() proto.Encoder {
	// TODO
	return nil
}

func (k *Key) Decoder() proto.Decoder {
	// TODO
	return nil
}

type AuthDB struct {
}

func NewAuthDB() *AuthDB {
	return nil
}

func (db *AuthDB) GenerateKey() (id string) {
	// TODO
	return ""
}

func (db *AuthDB) ReferToKey(id string) (*Key, error) {
	// TODO
	return nil, nil
}

func (db *AuthDB) ReleaseKey(id string) error {
	// TODO
	return nil
}

func (db *AuthDB) GC() {
	// TODO
}

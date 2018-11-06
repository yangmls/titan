package db

import (
	"encoding/json"
	"fmt"
)

const (
	// Separator of the key segment
	Separator = ":"
)

// ObjectEncoding is the encoding type of an object
type ObjectEncoding byte

// Encoding values, see https://github.com/antirez/redis/blob/unstable/src/server.h#L581
const (
	ObjectEncodingRaw = ObjectEncoding(iota)
	ObjectEncodingInt
	ObjectEncodingHT
	ObjectEncodingZipmap
	ObjectEncodingLinkedlist
	ObjectEncodingZiplist
	ObjectEncodingIntset
	ObjectEncodingSkiplist
	ObjectEncodingEmbstr
	ObjectEncodingQuicklist
)

// String representation of ObjectEncoding
func (enc ObjectEncoding) String() string {
	switch enc {
	case ObjectEncodingRaw:
		return "raw"
	case ObjectEncodingInt:
		return "int"
	case ObjectEncodingHT:
		return "hashtable"
	case ObjectEncodingZipmap:
		return "zipmap"
	case ObjectEncodingLinkedlist:
		return "linkedlist"
	case ObjectEncodingZiplist:
		return "ziplist"
	case ObjectEncodingIntset:
		return "intset"
	case ObjectEncodingSkiplist:
		return "skiplist"
	case ObjectEncodingEmbstr:
		return "embstr"
	case ObjectEncodingQuicklist:
		return "quicklist"
	default:
		return "unknown"
	}
}

// ObjectType is the type of a data structure
type ObjectType byte

// String representation of object type
func (t ObjectType) String() string {
	switch t {
	case ObjectString:
		return "string"
	case ObjectList:
		return "list"
	case ObjectSet:
		return "set"
	case ObjectZset:
		return "zset"
	case ObjectHash:
		return "hash"
	}
	return "none"
}

// Object types, see https://github.com/antirez/redis/blob/unstable/src/server.h#L461
const (
	ObjectString = ObjectType(iota)
	ObjectList
	ObjectSet
	ObjectZset
	ObjectHash
)

// Object meta schema
//   Layout {DB}:{TAG}:{Key}
//   DB     [0-255]
//   Key    Usersapce key
//   TAG    M(Meta), D(Data)
// Object data schema
//   Layout: {DB}:{TAG}:{ID}:{Others}
//   ID     Object ID, ID is not used for meta
// String schema (associated value with meta)
//   Layout: {DB}:M:{key}
type Object struct {
	ID        []byte
	Type      ObjectType     //refer to redis
	Encoding  ObjectEncoding //refer to redis
	CreatedAt int64
	UpdatedAt int64
	ExpireAt  int64
}

// String representation of an object
func (obj *Object) String() string {
	return fmt.Sprintf("ID:%s type:%s encoding:%s createdat:%d updatedat:%d expireat:%d",
		UUIDString(obj.ID), obj.Type, obj.Encoding, obj.CreatedAt, obj.UpdatedAt, obj.ExpireAt)
}

func (txn *Transaction) Object(key []byte) (*Object, error) {
	obj := &Object{}
	mkey := MetaKey(txn.db, key)

	meta, err := txn.txn.Get(mkey)
	if err != nil {
		if IsErrNotFound(err) {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}

	if err := json.Unmarshal(meta, obj); err != nil {
		return nil, err
	}

	return obj, nil
}

// Destory the object
func (txn *Transaction) Destory(obj *Object, key []byte) error {
	mkey := MetaKey(txn.db, key)
	dkey := DataKey(txn.db, obj.ID)

	if err := txn.txn.Delete(mkey); err != nil {
		return err
	}
	if obj.Type != ObjectString {
		return gc(txn, dkey)
	}
	return nil
}
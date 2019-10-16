package storage

import (
	"bytes"
	"fmt"
	"time"

	packed "github.com/EndCrystal/PackedIO"
	"github.com/EndCrystal/Server/world/chunk"
	"go.etcd.io/bbolt"
)

// Storage struct
type Storage struct{ db *bbolt.DB }

// Open open storage
func Open(path string) (ret *Storage, err error) {
	db, err := bbolt.Open(path, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err == nil {
		ret = &Storage{db}
	}
	return
}

// Close close storage
func (s Storage) Close() {
	s.db.Close()
}

// ForConfig create config lens
func (s *Storage) ForConfig(path ...string) SForConfig {
	if len(path) == 0 {
		panic("path is required")
	}
	return SForConfig{s, path}
}

// SForConfig custom lens for config
type SForConfig struct {
	*Storage
	path []string
}

// Get get key
func (s SForConfig) Get(key string) (ret []byte) {
	s.db.View(func(tx *bbolt.Tx) error {
		var ifc interface{ Bucket([]byte) *bbolt.Bucket } = tx
		for _, item := range s.path {
			bkt := ifc.Bucket([]byte(item))
			if bkt == nil {
				return nil
			}
			ifc = bkt
		}
		bkt := ifc.(*bbolt.Bucket)
		tmp := bkt.Get([]byte(key))
		ret = make([]byte, len(tmp))
		copy(ret, tmp)
		return nil
	})
	return
}

// Set set key
func (s SForConfig) Set(key string, value []byte) error {
	return s.db.Batch(func(tx *bbolt.Tx) error {
		var ifc interface {
			CreateBucketIfNotExists([]byte) (*bbolt.Bucket, error)
		} = tx
		for _, item := range s.path {
			var err error
			ifc, err = ifc.CreateBucketIfNotExists([]byte(item))
			if err != nil {
				return err
			}
		}
		bkt := ifc.(*bbolt.Bucket)
		return bkt.Put([]byte(key), value)
	})
}

// ForDim get dimension storage
func (s *Storage) ForDim(name string) SForDim {
	return SForDim{s, name}
}

// SForDim dimension storage
type SForDim struct {
	*Storage
	dim string
}

func dump(obj packed.Serializable) []byte {
	o, buf := packed.NewOutput()
	obj.Save(o)
	return buf.Bytes()
}

// LoadChunk load chunk
func (s SForDim) LoadChunk(pos chunk.CPos) (chk *chunk.Chunk, err error) {
	dimchk := []byte(fmt.Sprintf("dim-%s-chunks", s.dim))
	err = s.db.View(func(tx *bbolt.Tx) error {
		bkt := tx.Bucket(dimchk)
		if bkt == nil {
			return chunk.ErrChunkNotFound
		}
		data := bkt.Get(dump(&pos))
		var chkdata chunk.Chunk
		chkdata.Load(packed.MakeInput(bytes.NewReader(data)))
		chk = &chkdata
		return nil
	})
	return
}

// SaveChunk save chunk
func (s SForDim) SaveChunk(pos chunk.CPos, data *chunk.Chunk) error {
	dimchk := []byte(fmt.Sprintf("dim-%s-chunks", s.dim))
	posd := dump(&pos)
	datad := dump(data)
	err := s.db.Batch(func(tx *bbolt.Tx) error {
		bkt, e := tx.CreateBucketIfNotExists(dimchk)
		if e != nil {
			return e
		}
		return bkt.Put(posd, datad)
	})
	return err
}

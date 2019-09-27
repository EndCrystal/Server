package storage

import (
	"bytes"
	"fmt"
	"time"

	packed "github.com/EndCrystal/PackedIO"
	"github.com/EndCrystal/Server/world/chunk"
	"go.etcd.io/bbolt"
)

type Storage struct{ db *bbolt.DB }

func Open(path string) (ret *Storage, err error) {
	db, err := bbolt.Open(path, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err == nil {
		ret = &Storage{db}
	}
	return
}

func (s Storage) Close() {
	s.db.Close()
}

func (s *Storage) ForDim(name string) StorageForDim {
	return StorageForDim{s, name}
}

type StorageForDim struct {
	*Storage
	dim string
}

func dump(obj packed.Serializable) []byte {
	o, buf := packed.NewOutput()
	obj.Save(o)
	return buf.Bytes()
}

func (s StorageForDim) LoadChunk(pos chunk.ChunkPos) (chk *chunk.Chunk, err error) {
	dimchk := []byte(fmt.Sprintf("dim-%s-chunks", s.dim))
	err = s.db.View(func(tx *bbolt.Tx) error {
		bkt := tx.Bucket(dimchk)
		if bkt == nil {
			return chunk.EChunkNotFound
		}
		data := bkt.Get(dump(&pos))
		var chkdata chunk.Chunk
		chkdata.Load(packed.MakeInput(bytes.NewReader(data)))
		chk = &chkdata
		return nil
	})
	return
}

func (s StorageForDim) SaveChunk(pos chunk.ChunkPos, data *chunk.Chunk) error {
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

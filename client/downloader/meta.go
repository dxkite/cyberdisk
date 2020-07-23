package downloader

import (
	"dxkite.cn/storage/bitset"
	"dxkite.cn/storage/meta"
	"encoding/gob"
	"os"
)

const xor = 0x14

type ProcessMeta struct {
	meta.Info
	Index         bitset.BitSet
	Downloaded    int
	DownloadTotal int
}

func EncodeToFile(path string, info *ProcessMeta) error {
	f, er := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if er != nil {
		return er
	}
	defer func() { _ = f.Close() }()
	b := gob.NewEncoder(meta.NewXORWriter(xor, f))
	return b.Encode(info)
}

func DecodeToFile(path string) (*ProcessMeta, error) {
	f, er := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if er != nil {
		return nil, er
	}
	defer func() { _ = f.Close() }()
	b := gob.NewDecoder(meta.NewXORReader(xor, f))
	info := new(ProcessMeta)
	der := b.Decode(&info)
	if der != nil {
		return nil, der
	}
	return info, nil
}

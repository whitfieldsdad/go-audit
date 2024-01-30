package monitor

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"github.com/zeebo/xxh3"
)

type Hashes struct {
	MD5    string `json:"md5,omitempty"`
	SHA1   string `json:"sha1,omitempty"`
	SHA256 string `json:"sha256,omitempty"`
	XXH3   uint64 `json:"xxh3,omitempty"`
}

func (h *Hashes) Empty() bool {
	return h.MD5 == "" && h.SHA1 == "" && h.SHA256 == ""
}

func GetHashes(rd io.Reader) (*Hashes, error) {
	md5h := md5.New()
	sha1h := sha1.New()
	sha256h := sha256.New()
	xxh3h := xxh3.New()

	pagesize := os.Getpagesize()
	reader := bufio.NewReaderSize(rd, pagesize)
	multiWriter := io.MultiWriter(md5h, sha1h, sha256h, xxh3h)
	_, err := io.Copy(multiWriter, reader)
	if err != nil {
		return nil, err
	}
	hashes := &Hashes{
		MD5:    fmt.Sprintf("%x", md5h.Sum(nil)),
		SHA1:   fmt.Sprintf("%x", sha1h.Sum(nil)),
		SHA256: fmt.Sprintf("%x", sha256h.Sum(nil)),
		XXH3:   xxh3h.Sum64(),
	}
	return hashes, nil
}

func GetFileHashes(path string) (*Hashes, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return GetHashes(f)
}

func GetXXH3(data []byte) uint64 {
	return xxh3.Hash(data)
}

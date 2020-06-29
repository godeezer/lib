package deezer

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/erebid/go-deezer/deezer/crypto/ecb"
	"golang.org/x/crypto/blowfish"
)

const blowfishSecret = "g4el58wc0zvf9na1"
const filenameKey = "jo6aey6haid2Teih"

func songDownloadURL(s SongData, preferred Quality) (string, error) {
	key, err := songFilename(s, preferred)
	if err != nil {
		return "", err
	}
	cdn := string(s.MD5Origin[0])
	return "https://e-cdns-proxy-" + cdn + ".dzcdn.net/mobile/1/" + key, err
}

func songFilename(s SongData, preferred Quality) (string, error) {
	quality, err := GetValidSongQuality(s, preferred)
	if err != nil {
		return "", err
	}
	q := strconv.Itoa(int(quality))
	step1 := strings.Join(
		[]string{
			s.MD5Origin,
			q,
			s.ID,
			s.MediaVersion,
		},
		"\xa4",
	)
	sum := md5.Sum([]byte(step1))
	step2 := fmt.Sprintf("%x\xa4%s\xa4", sum[:], step1)
	for len(step2)%16 > 0 {
		step2 += " "
	}
	key := []byte(filenameKey)
	ciphertext := encryptAes128ECB([]byte(step2), key)
	return hex.EncodeToString(ciphertext), err
}

func encryptAes128ECB(pt, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	mode := ecb.NewECBEncrypter(block)
	ct := make([]byte, len(pt))
	mode.CryptBlocks(ct, pt)
	return ct
}

type EncryptedSongReader struct {
	R      io.Reader
	S      SongData
	i      int
	chunk  []byte
	buffer *bytes.Buffer
}

func NewEncryptedSongReader(r io.Reader, s SongData) (*EncryptedSongReader, error) {
	reader := &EncryptedSongReader{R: r, S: s}
	return reader, nil
}

func (r *EncryptedSongReader) Read(p []byte) (int, error) {
	buf := bytes.Buffer{}
	for buf.Len() < len(p) {
		chunk, err := r.ReadChunk()
		buf.Write(chunk)
		if err != nil {
			n, _ := buf.Read(p)
			return n, err
		}
	}
	return buf.Read(p)
}

func (r *EncryptedSongReader) ReadChunk() ([]byte, error) {
	chunk := make([]byte, 2048)
	_, err := io.ReadFull(r.R, chunk)
	if err != nil {
		if err == io.ErrUnexpectedEOF {
			return chunk, io.EOF
		} else {
			return chunk, err
		}
	}
	if r.i%3 == 0 {
		fmt.Println("predecrypt", md5.Sum(chunk))
		chunk, _ := decryptChunk(chunk, getBlowfishKey(r.S.ID))
		//r.cryptmode.CryptBlocks(chunk, chunk)
		println("postdecrypt", chunk[0])
	}
	r.i++
	return chunk, nil
}

func getBlowfishKey(id string) []byte {
	idmd5 := md5.Sum([]byte(id))
	idmd5hex := hex.EncodeToString(idmd5[:])
	var key string
	for i := 0; i < 16; i++ {
		r := idmd5hex[i] ^ idmd5hex[i+16] ^ blowfishSecret[i]
		key += string(r)
	}
	return []byte(key)
}

func decryptChunk(chunk, key []byte) ([]byte, error) {
	ci, err := blowfish.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := []byte{0, 1, 2, 3, 4, 5, 6, 7}
	cbcDecrypter := cipher.NewCBCDecrypter(ci, iv)
	cbcDecrypter.CryptBlocks(chunk, chunk)
	return chunk, nil
}

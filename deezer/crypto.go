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

// SongDownloadURL returns a download URL which can be used to stream the song.
// The audio returned from the URL will be encrypted so you should use
// a EncryptedSongReader to read it.
func SongDownloadURL(s Song, quality Quality) string {
	key := songFilename(s, quality)
	cdn := string(s.MD5Origin[0])
	return "https://e-cdns-proxy-" + cdn + ".dzcdn.net/mobile/1/" + key
}

func songFilename(s Song, quality Quality) string {
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
	return hex.EncodeToString(ciphertext)
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
	r   io.Reader
	s   Song
	i   int
	buf bytes.Buffer
}

// NewEncryptedSongReader creates an EncryptedSongReader
// that reads from r and decrypts it using s.
func NewEncryptedSongReader(r io.Reader, s Song) (*EncryptedSongReader, error) {
	reader := &EncryptedSongReader{r: r, s: s}
	return reader, nil
}

// Read reads up to n(p) bytes into p, returning how many bytes
// were read and any error.
func (r *EncryptedSongReader) Read(p []byte) (int, error) {
	for r.buf.Len() < len(p) {
		chunk, err := r.ReadChunk()
		r.buf.Write(chunk)
		if err != nil {
			n, _ := r.buf.Read(p)
			return n, err
		}
	}
	return r.buf.Read(p)
}

// ReadChunk returns the next n<=2048 bytes of the song.
// It automatically decrypts chunks when it has to (every third chunk).
// You most likely would prefer to use Read instead of ReadChunk because it implements io.Reader
func (r *EncryptedSongReader) ReadChunk() ([]byte, error) {
	chunk := make([]byte, 2048)
	_, err := io.ReadFull(r.r, chunk)
	if err != nil {
		if err == io.ErrUnexpectedEOF {
			return chunk, io.EOF
		} else {
			return chunk, err
		}
	}
	if r.i%3 == 0 {
		chunk, err = decryptChunk(chunk, getBlowfishKey(r.s.ID))
		if err != nil {
			return chunk, err
		}
	}
	r.i++
	return chunk, nil
}

// getBlowfishKey returns the Blowfish key for a given song by its id.
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

// decryptChunk decrypts a 2048-byte chunk using the key.
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

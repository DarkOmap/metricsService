package hasher

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type Hasher struct {
	key []byte
}

func NewHasher(key []byte) Hasher {
	return Hasher{key}
}

func (h *Hasher) HashingRequest(req *resty.Request, body []byte) {
	if len(h.key) == 0 {
		return
	}

	hash := hmac.New(sha256.New, h.key)
	hash.Write(body)
	req.SetHeader("HashSHA256", hex.EncodeToString(hash.Sum(nil)))
}

type hashingResponseWriter struct {
	http.ResponseWriter
	bytes int
	key   []byte
}

func (r *hashingResponseWriter) Write(b []byte) (int, error) {
	h := hmac.New(sha256.New, []byte(r.key))
	h.Write(b)
	dst := h.Sum(nil)

	r.ResponseWriter.Header().Add("HashSHA256", hex.EncodeToString(dst))

	size, err := r.ResponseWriter.Write(b)
	r.bytes += size
	return size, err
}

func (h *Hasher) RequestHash(handler http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		hashHeader := r.Header.Get("HashSHA256")

		if len(h.key) == 0 || hashHeader == "" {
			handler.ServeHTTP(w, r)
			return
		}

		var buf bytes.Buffer
		_, err := buf.ReadFrom(r.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		r.Body = io.NopCloser(&buf)
		hash := hmac.New(sha256.New, h.key)
		hash.Write(buf.Bytes())
		dst := hash.Sum(nil)
		hh, err := hex.DecodeString(hashHeader)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !hmac.Equal(hh, dst) {
			http.Error(w, "hash not equal", http.StatusBadRequest)
			return
		}

		hw := hashingResponseWriter{
			ResponseWriter: w,
			key:            h.key,
		}

		handler.ServeHTTP(&hw, r)
	}

	return http.HandlerFunc(logFn)
}

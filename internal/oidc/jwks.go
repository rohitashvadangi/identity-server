package oidc

import (
    "encoding/json"
    "net/http"
	"crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "errors"
    "io/ioutil"
    "github.com/lestrrat-go/jwx/jwk"
)

var publicKey *rsa.PublicKey

func InitJWKS(pubKeyPath string) error {
    var err error
    publicKey, err = LoadRSAPublicKey(pubKeyPath)
    return err
}

func JWKSHandler(w http.ResponseWriter, r *http.Request) {
    set := jwk.NewSet()
    key, err := jwk.New(publicKey)
    if err != nil {
        http.Error(w, "failed to create JWK", http.StatusInternalServerError)
        return
    }

    // Optional: set key ID (kid)
    key.Set("kid", "my-key-id")
    set.Add(key)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(set)
}


// LoadRSAPublicKey reads a PEM encoded RSA public key
func LoadRSAPublicKey(path string) (*rsa.PublicKey, error) {
    keyBytes, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }

    block, _ := pem.Decode(keyBytes)
    if block == nil || block.Type != "PUBLIC KEY" {
        return nil, errors.New("failed to decode PEM block containing public key")
    }

    pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return nil, err
    }

    rsaPub, ok := pubKey.(*rsa.PublicKey)
    if !ok {
        return nil, errors.New("not an RSA public key")
    }

    return rsaPub, nil
}
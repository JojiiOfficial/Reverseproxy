package models

import (
	"crypto/tls"
)

//TLSKeyCertPair pair of tls key/cert file
type TLSKeyCertPair struct {
	Key  string
	Cert string
}

// GetCertificate gets certificate from TLSKeyCertPair
func (pair TLSKeyCertPair) GetCertificate() (tls.Certificate, error) {
	return tls.LoadX509KeyPair(pair.Cert, pair.Key)
}

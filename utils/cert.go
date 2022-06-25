package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	rd "math/rand"
	"net"
	"os"
	"time"
)

func CreateCert(masterIP string) {
	_, err := os.Stat("resources/tls.key")
	if err == nil {
		return
	} else if !os.IsNotExist(err) {
		fmt.Println(err)
		os.Exit(1)
	}

	//extSubjectAltName := pkix.Extension{}
	//extSubjectAltName.Id = asn1.ObjectIdentifier{2, 5, 29, 17}
	//extSubjectAltName.Critical = false
	//extSubjectAltName.Value = []byte(fmt.Sprintf("IP:%s", masterIP))

	CACert := &x509.Certificate{
		SerialNumber:          big.NewInt(rd.Int63()),
		BasicConstraintsValid: true,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		Subject: pkix.Name{
			CommonName: "ca-" + masterIP,
		},
		Issuer: pkix.Name{
			CommonName: "ca-" + masterIP,
		},
		IsCA: true,
		//ExtKeyUsage:    []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		//KeyUsage:       x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		EmailAddresses: []string{},
		//ExtraExtensions: []pkix.Extension{extSubjectAltName},
	}
	Cert := &x509.Certificate{
		SerialNumber:          big.NewInt(rd.Int63()),
		BasicConstraintsValid: true,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		Subject: pkix.Name{
			CommonName: "ca-" + masterIP,
		},
		Issuer: pkix.Name{
			CommonName: "ca-" + masterIP,
		},
		IsCA:           false,
		ExtKeyUsage:    []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:       x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		EmailAddresses: []string{},
		//ExtraExtensions: []pkix.Extension{extSubjectAltName},
		IPAddresses: []net.IP{net.ParseIP(masterIP)},
	}
	Key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var buf []byte
	buf, err = x509.CreateCertificate(rand.Reader, CACert, CACert, &Key.PublicKey, Key)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	caFile, err := os.Create("resources/ca.crt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = pem.Encode(caFile, &pem.Block{Type: "CERTIFICATE", Bytes: buf})
	//_, err = caFile.Write(buf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = caFile.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//CACert, _ = x509.ParseCertificate(buf)
	buf = nil
	buf, err = x509.CreateCertificate(rand.Reader, Cert, CACert, &Key.PublicKey, Key)
	certFile, err := os.Create("resources/tls.crt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: buf})
	//_, err = certFile.Write(buf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = certFile.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	buf = nil
	buf = x509.MarshalPKCS1PrivateKey(Key)
	keyFile, err := os.Create("resources/tls.key")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: buf})
	//_, err = keyFile.Write(buf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = keyFile.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

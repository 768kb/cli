//go:build integration
// +build integration

package integration

import (
	"fmt"
	"testing"
)

var testdata = "testdata"

type CertificateSignCmd struct {
	name      string
	command   CLICommand
	csr       string
	issuerCrt string
	issuerKey string
	pass      string
}

func (k CertificateSignCmd) setPass(pass string) CertificateSignCmd {
	return CertificateSignCmd{k.name, k.command, k.csr, k.issuerCrt, k.issuerKey, pass}
}

func (k CertificateSignCmd) fail(t *testing.T, expected string) {
	k.command.fail(t, k.name, expected, "")
}

func (k CertificateSignCmd) failNoPass(t *testing.T, expected string) {
	k.command.fail(t, k.name, expected, "")
}

func NewCertificateSignCmd(name, csr, crt, key string) CertificateSignCmd {
	csrFile := fmt.Sprintf("%s/%s", testdata, csr)
	crtFile := fmt.Sprintf("%s/%s", testdata, crt)
	keyFile := fmt.Sprintf("%s/%s", testdata, key)
	command := NewCLICommand().setCommand(fmt.Sprintf("step certificate sign %s %s %s",
		csrFile, crtFile, keyFile))
	return CertificateSignCmd{name, command, csrFile, crtFile, keyFile, "password"}
}

func TestCertificateSign(t *testing.T) {
	NewCertificateSignCmd("bad-sig", "certificate-create-bad-sig.csr", "intermediate_ca.crt", "intermediate_ca_key").failNoPass(t, "Certificate Request has invalid signature: crypto/rsa: verification error\n")
	//NewKeypairCmd("success", "foo.csr", "intermediate_ca.crt", "intermediate_ca_key").setPass("pass").test(t)
}

package utils

import (
	"reflect"
	"testing"
	"testing/fstest"
)

// To regenerate: openssl genrsa -out /tmp/test.pem
const (
	privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAzVaDKmmoG4ytBGFoiTdSAF4Miks01veNUgIyxGzX8QKF159w
Y61La3KIYfOHGBgfy/hCt08aL9cPiu8hojN7xN9fWTTv50Xl25bR5jx+SCzvwWGS
YLCVKhMYEHSCfrPGEhJoEfKhrQoTCilkQTgtEEq9hUoZW4Kgb0qKg7x38fQlZ6si
LAGs5pvJDFmjK5UAMwz+jkSdN6d1kjL/EkBX69FklBNYhngxJWTmXRMKkziktI+j
9ZdO1bxS3X+XsKifx/g6aPIQiRRmfoc9HMVfn/ep2w7DoztHza1PfUr015YBGqOh
VxX1JajVcN/IB4YuZOvHjBjskEpG7QWQIPLASQIDAQABAoIBAQCv/8ALYWZqvqgp
wggk3JrXn8Ul4BJZUvP5X9L908E/XYc06v8dIJMtdIz7UA3yE/NlE9SzZASxDqfO
0OrGKVSjyUXjo2EhnSLIlbwxmJYw7PtPiH87iv8/ggA1Unfre9GA+e/julDjjWgG
ZLX+xNSzSyyoi3uymQNEgOK6yZcRdUZ0yCz6UTlAdEK26GiN0BdnIJG9GvG5vHA8
bkTmMgOpIQOIJynCyX/q+iO8hPj6tmBIkTbdzsmrRShIi1KIF054het3gHhRBRtX
q86enmsU9lSlzEj9HrrwaHppDEZ8qTQsqWzuIqtjdXgcQNH52mBOU+KO/gL66tHs
NpoxtRKFAoGBAOXh6kLlDnx+4io36tF0E6D2phrPTpEJMafD+t+ggofCRU8rwDuf
PMi0Qz1ggqtPd0e4h7+/r+Cq2jDtj/XCeLq/IKg1E6cdbbTXop27PCwjK+6Yl+22
kCGymCQkscJgO86QCU8S5fxTw0/VG9IjzPAkBYCjSY+a1jCZS1MsFqUnAoGBAOSq
uWwIf+rdtT1zwLHHJKuED9JzyNbl/jXINdjiovK+KMfQNfk+AeLupABpabuOlwaE
pha4NZddndVnCPrsqT7GKpCmtIaprLw/r8deq1BynOI6m44gS6UPKvog/bdWBp3L
cSpxEescXm4r2K1bM7iiJLQ+zyT9hVi8kkzZ4jUPAoGBAKssreNh7IeHc6E8Qf31
ESiqgMU12KrmzbK+m/Ao9Qlh/3oUee/rgrdwgyEQ3Dvz0D33ih2d/risgAwu2SOG
y59C8m5OF3Q41ZfzeYM6CHRVPEFOHtNDPc/ZzLAdIsA6KE6HsmbPC7H4LVckuLKh
Ndka+X3wGLZ19Uf63bvw+GvBAoGAH/7LZxRhYamX/Hs/0SA+P0mBNT9CMN+JjFjx
P+GmTzTQW/UEOFW2ydv+UphtVPMEqsLQwokP5pgQx5VdKk8G92Oe/RJ2XAlNxCFd
JRZX/i+rR/RPY7mdHAFdUBZhqc99qYKX2QptKWqUw/GapdcHC6SUYiwPq+tVRy9L
gTlTb30CgYB4fC5Ejsx/gPoCTNu9uiVwpSH9W7GXA1PpQ35GN5TqSb8Ia0L/oqX3
vu6V83YC+zeHSWGJce0Eqwf3wwL8UrlvUQfwSRTDMmJhZsyNhWMqM9gqBLrs+BM2
DwQ6MDAzQulsdFC7fOHIJ0a4DhU4e9AtGw+gkVVolSvIcnK1d64cug==
-----END RSA PRIVATE KEY-----
`

	publicKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDL23od5lgLNeUCH5IgtcaPzdrj1IXVxnK0BnoEx5+4nrhaLAMqex46NcJrFp/Jevc0tWXUvHw0Rtz3Gss33DHw1bz1F642fJidzcpsUsdd3jqp3JFExVxZUSIGdT29CpFqLtJ6NIIXr0xef4dsqwgvg2jWbDP0FaWATLhoC2d5wf1fg2RyYxrCTdq/a7BkP0oeDIDSWR2l9uwyz47ly9NHtitDHaEPOFnBRJIkd7bQ/hZufwaatYcBvYAtjDdFRRLW5JpWt2PGmYgDt3tKyw/kcVXwwPe8+K5pvyOP+JmYyfr6SusMb0Dqegzhuwu4PunZeBWSVgHvQLvKKhsXv1mOF0ZG9dS6i65dDIrce4WyU3Bgq5eJ30ZnnB8gJzUFNWPfosqISIUTVSzNve0N/g/zWFMvRyaK58EYviQUXZfcCMqzLenS33Em5vXp20XEUaE2hN5WsLqGUfqQJFa1DEItKESWzZgYZ5X4Y9955azDrRemTwElO3aIrtDwEX/U7K8= test-user@localhost"
)

// TestDerefString tests that we can dereference string pointers to strings correctly
func TestDerefString(t *testing.T) {
	var input *string
	output := DerefString(input)

	got := reflect.TypeOf(output).Kind()
	expected := reflect.String

	if expected != got {
		t.Errorf("error DerefString(): expected: %s | got: %s", expected, got)
	}
}

// TestPublicSSHKey mocks a filesystem and tests that we can read data from a public ssh key
func TestPublicSSHKey(t *testing.T) {
	publicKeyFile := "~/.ssh/id_rsa.pub"

	fs := fstest.MapFS{
		publicKeyFile: {
			Data: []byte(publicKey),
		},
	}

	data, err := fs.ReadFile(publicKeyFile)
	if err != nil {
		t.Error(err)
	}

	expected := publicKey
	got := string(data)

	if expected != got {
		t.Errorf("expected: %s | got: %s", expected, got)
	}
}

// TestPrivateSSHKey mocks a filesystem and tests that we can read data from a private ssh key
func TestPrivateSSHKey(t *testing.T) {
	privateKeyFile := "~/.ssh/id_rsa"

	fs := fstest.MapFS{
		privateKeyFile: {
			Data: []byte(privateKey),
		},
	}

	data, err := fs.ReadFile(privateKeyFile)
	if err != nil {
		t.Error(err)
	}

	expected := privateKey
	got := string(data)

	if expected != got {
		t.Errorf("expected: %s | got: %s", expected, got)
	}
}

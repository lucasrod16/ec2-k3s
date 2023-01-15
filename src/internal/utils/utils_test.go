package utils

import (
	"reflect"
	"testing"
	"testing/fstest"
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
	publicSSHKey := "~/.ssh/id_rsa.pub"
	keyData := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDL23od5lgLNeUCH5IgtcaPzdrj1IXVxnK0BnoEx5+4nrhaLAMqex46NcJrFp/Jevc0tWXUvHw0Rtz3Gss33DHw1bz1F642fJidzcpsUsdd3jqp3JFExVxZUSIGdT29CpFqLtJ6NIIXr0xef4dsqwgvg2jWbDP0FaWATLhoC2d5wf1fg2RyYxrCTdq/a7BkP0oeDIDSWR2l9uwyz47ly9NHtitDHaEPOFnBRJIkd7bQ/hZufwaatYcBvYAtjDdFRRLW5JpWt2PGmYgDt3tKyw/kcVXwwPe8+K5pvyOP+JmYyfr6SusMb0Dqegzhuwu4PunZeBWSVgHvQLvKKhsXv1mOF0ZG9dS6i65dDIrce4WyU3Bgq5eJ30ZnnB8gJzUFNWPfosqISIUTVSzNve0N/g/zWFMvRyaK58EYviQUXZfcCMqzLenS33Em5vXp20XEUaE2hN5WsLqGUfqQJFa1DEItKESWzZgYZ5X4Y9955azDrRemTwElO3aIrtDwEX/U7K8= test-user@localhost"
	fs := fstest.MapFS{
		publicSSHKey: {
			Data: []byte(keyData),
		},
	}
	data, err := fs.ReadFile(publicSSHKey)
	if err != nil {
		t.Error(err)
	}

	expected := keyData
	got := string(data)

	if expected != got {
		t.Errorf("expected: %s | got: %s", expected, got)
	}

}

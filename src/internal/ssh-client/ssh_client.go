package ssh

import (
	"bytes"
	"io"
	"os"
	"sync"

	"github.com/lucasrod16/ec2-k3s/src/internal/utils"
	"golang.org/x/crypto/ssh"
)

// SSHClient initializes a ssh client connection
type SSHClient struct {
	conn *ssh.Client
}

// ExecuteCommand executes a command on a remote machine to install k3s
type ExecuteCommand interface {
	Execute(command string) (CommandOutput, error)
	ExecuteStdio(command string) (CommandOutput, error)
}

// CommandOutput contains the STDIO output from running a command
type CommandOutput struct {
	StdOut []byte
	StdErr []byte
}

// NewSSHClient creates a new ssh client connection
// with the provdided host and configuration
func NewSSHClient(host string, config *ssh.ClientConfig) (*SSHClient, error) {
	conn, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return nil, err
	}

	client := SSHClient{
		conn: conn,
	}

	return &client, nil
}

// ExecuteStdio pipes the remote command output to local stdio
func (s SSHClient) ExecuteOutput(command string, stream bool) (CommandOutput, error) {
	sess, err := s.conn.NewSession()
	if err != nil {
		return CommandOutput{}, err
	}

	defer sess.Close()

	sessStdOut, err := sess.StdoutPipe()
	if err != nil {
		return CommandOutput{}, err
	}

	output := bytes.Buffer{}
	wg := sync.WaitGroup{}

	var stdOutWriter io.Writer
	if stream {
		stdOutWriter = io.MultiWriter(os.Stdout, &output)
	} else {
		stdOutWriter = &output
	}

	wg.Add(1)
	go func() {
		io.Copy(stdOutWriter, sessStdOut)
		wg.Done()
	}()

	sessStderr, err := sess.StderrPipe()
	if err != nil {
		return CommandOutput{}, err
	}

	errorOutput := bytes.Buffer{}
	var stdErrWriter io.Writer
	if stream {
		stdErrWriter = io.MultiWriter(os.Stderr, &errorOutput)
	} else {
		stdErrWriter = &errorOutput
	}

	wg.Add(1)
	go func() {
		io.Copy(stdErrWriter, sessStderr)
		wg.Done()
	}()

	err = sess.Run(command)
	if err != nil {
		return CommandOutput{}, err
	}

	wg.Wait()

	return CommandOutput{
		StdErr: errorOutput.Bytes(),
		StdOut: output.Bytes(),
	}, nil
}

func (s SSHClient) Execute(command string) (CommandOutput, error) {
	return s.ExecuteOutput(command, true)
}

func (s SSHClient) Close() error {
	return s.conn.Close()
}

// ConfigureSSHClient configures a ssh client
// with a user, host, and ssh keys
func ConfigureSSHClient(region string) (*SSHClient, error) {
	user := "ubuntu"
	privateKey := utils.GetPrivateSSHKey()

	ip, err := utils.GetInstanceIp(region)
	if err != nil {
		return nil, err
	}

	port := "22"
	host := ip + ":" + port

	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshClient, err := NewSSHClient(host, config)
	if err != nil {
		return nil, err
	}

	return sshClient, nil
}

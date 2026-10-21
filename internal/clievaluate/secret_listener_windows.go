//go:build windows

package clievaluate

import (
	"net"

	"github.com/Microsoft/go-winio"
	"github.com/google/uuid"
	"golang.org/x/sys/windows"
)

func secretListener(string) (net.Listener, string, error) {
	user, err := windows.GetCurrentProcessToken().GetTokenUser()
	if err != nil {
		return nil, "", err
	}
	endpoint := `\\.\pipe\gosched-secret-` + uuid.NewString()
	listener, err := winio.ListenPipe(endpoint, &winio.PipeConfig{SecurityDescriptor: "D:P(A;;GA;;;" + user.User.Sid.String() + ")"})
	return listener, endpoint, err
}

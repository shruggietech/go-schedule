package notification

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/testsmtp"
)

func TestSMTPFinalDataRejection(t *testing.T) {
	for _, tc := range []struct {
		code  int
		state string
	}{{450, "temporary_failure"}, {550, "permanent_failure"}} {
		t.Run(fmt.Sprint(tc.code), func(t *testing.T) {
			receiver, err := testsmtp.Start(testsmtp.Options{DataCode: tc.code, ReplyText: "DO-NOT-EXPOSE-server-text"})
			if err != nil {
				t.Fatal(err)
			}
			defer receiver.Close()
			d := mailFixtureDelivery(t, receiver, "plain")
			result := NewSender(nil).SendSMTP(context.Background(), d, func([]domain.NotificationRecipient) error { return nil })
			for _, r := range result.Recipients {
				if r.State != tc.state || r.Code != tc.code || strings.Contains(r.Diagnostic, "DO-NOT-EXPOSE") {
					t.Fatalf("recipient %+v", r)
				}
			}
			if len(receiver.Messages()) != 0 {
				t.Fatal("rejected DATA recorded as accepted")
			}
		})
	}
}

// scriptedSMTP owns one bounded connection and controls the final acknowledgement.
func scriptedSMTP(t *testing.T, dropAck bool, stallData bool) (string, <-chan struct{}) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})
	t.Cleanup(func() { _ = listener.Close(); <-done })
	go func() {
		defer close(done)
		defer listener.Close()
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		_ = conn.SetDeadline(time.Now().Add(4 * time.Second))
		reader := bufio.NewReader(conn)
		reply := func(s string) bool { _, err := fmt.Fprint(conn, s+"\r\n"); return err == nil }
		if !reply("220 owned fixture") {
			return
		}
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			command := strings.ToUpper(strings.TrimSpace(line))
			switch {
			case strings.HasPrefix(command, "EHLO"), strings.HasPrefix(command, "HELO"):
				if !reply("250 owned fixture") {
					return
				}
			case strings.HasPrefix(command, "MAIL"), strings.HasPrefix(command, "RCPT"):
				if !reply("250 OK") {
					return
				}
			case command == "DATA":
				if stallData {
					var b [1]byte
					_, _ = conn.Read(b[:])
					return
				}
				if !reply("354 Continue") {
					return
				}
				for {
					line, err = reader.ReadString('\n')
					if err != nil {
						return
					}
					if line == ".\r\n" {
						break
					}
				}
				if dropAck {
					return
				}
				if !reply("250 Accepted") {
					return
				}
			case command == "QUIT":
				return // Deliberately omit the cleanup acknowledgement.
			default:
				return
			}
		}
	}()
	return listener.Addr().String(), done
}

func TestSMTPUncertainSubmissionAndQuitCleanup(t *testing.T) {
	for _, drop := range []bool{true, false} {
		t.Run(fmt.Sprint(drop), func(t *testing.T) {
			address, _ := scriptedSMTP(t, drop, false)
			cfg := smtpTestConfig()
			host, port, err := net.SplitHostPort(address)
			if err != nil {
				t.Fatal(err)
			}
			cfg.Host = host
			cfg.Port, err = strconv.Atoi(port)
			if err != nil {
				t.Fatal(err)
			}
			cfg.TLSMode = "plain"
			d := domain.NotificationDelivery{ID: "uncertain-message", SMTP: &cfg, Payload: []byte(`{"event":"test"}`)}
			checkpoints := 0
			result := NewSender(nil).SendSMTP(context.Background(), d, func(recipients []domain.NotificationRecipient) error {
				checkpoints++
				if recipients[0].State != "submitting" {
					t.Fatalf("checkpoint %+v", recipients)
				}
				return nil
			})
			expected := "accepted"
			if drop {
				expected = "uncertain"
			}
			if checkpoints != 1 || result.Recipients[0].State != expected {
				t.Fatalf("result %+v", result)
			}
			if drop && !strings.Contains(result.Diagnostic, "duplicate") {
				t.Fatalf("missing duplicate explanation %+v", result)
			}
		})
	}
}

func TestSMTPSharedDeadlineBoundsData(t *testing.T) {
	address, _ := scriptedSMTP(t, false, true)
	cfg := smtpTestConfig()
	cfg.Host, _, _ = net.SplitHostPort(address)
	_, port, _ := net.SplitHostPort(address)
	fmt.Sscan(port, &cfg.Port)
	cfg.TLSMode = "plain"
	cfg.TimeoutSeconds = 1
	start := time.Now()
	result := NewSender(nil).SendSMTP(context.Background(), domain.NotificationDelivery{ID: "deadline", SMTP: &cfg, Payload: []byte(`{"event":"test"}`)}, func([]domain.NotificationRecipient) error { return nil })
	if time.Since(start) > 3*time.Second || result.Recipients[0].State != "temporary_failure" || !strings.Contains(result.Diagnostic, "timed out") {
		t.Fatalf("deadline result %+v", result)
	}
}

func TestSMTPInvalidSnapshotFailsRecipients(t *testing.T) {
	cfg := smtpTestConfig()
	cfg.Port = -1
	d := domain.NotificationDelivery{SMTP: &cfg, Recipients: []domain.NotificationRecipient{{Address: "ops@example.test", State: "pending"}, {Address: "done@example.test", State: "accepted"}}}
	result := NewSender(nil).SendSMTP(context.Background(), d, nil)
	if result.Recipients[0].State != "permanent_failure" || result.Recipients[1].State != "accepted" {
		t.Fatalf("invalid snapshot left retryable recipients %+v", result)
	}
}

func TestSMTPRejectsExpiredAndWrongHostCertificates(t *testing.T) {
	for _, mode := range []string{"starttls", "tls"} {
		for _, expired := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/expired=%t", mode, expired), func(t *testing.T) {
				options := testsmtp.Options{TLSMode: mode, ExpiredCertificate: expired}
				if !expired {
					options.CertificateHost = "different.example.test"
				}
				receiver, err := testsmtp.Start(options)
				if err != nil {
					t.Fatal(err)
				}
				defer receiver.Close()
				d := mailFixtureDelivery(t, receiver, mode)
				result := NewSenderWithSMTPRoots(receiver.Roots()).SendSMTP(context.Background(), d, func([]domain.NotificationRecipient) error { return nil })
				for _, r := range result.Recipients {
					if r.State != "permanent_failure" {
						t.Fatalf("certificate accepted %+v", result)
					}
				}
				if len(receiver.Messages()) != 0 || !strings.Contains(result.Diagnostic, "trusted server certificate") {
					t.Fatalf("TLS evidence %+v", result)
				}
			})
		}
	}
}

func TestSMTPSocketTimeoutDiagnosticWithoutContextCancellation(t *testing.T) {
	if context.Background().Err() != nil {
		t.Fatal("background context unexpectedly canceled")
	}
	result := smtpFail(domain.SMTPAttemptResult{Recipients: []domain.NotificationRecipient{{Address: "ops@example.test", State: "pending"}}}, []int{0}, "message submission", &net.OpError{Op: "read", Net: "tcp", Err: os.ErrDeadlineExceeded}, true)
	if result.Recipients[0].State != "temporary_failure" || !strings.Contains(result.Diagnostic, "timed out") {
		t.Fatalf("socket timeout diagnostic %+v", result)
	}
}

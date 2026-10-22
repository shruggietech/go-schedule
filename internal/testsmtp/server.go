// Package testsmtp supplies an owned, bounded local SMTP receiver for evaluation.
package testsmtp

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"io"
	"log"
	"math/big"
	"net"
	"sync"
	"time"

	"github.com/emersion/go-sasl"
	smtp "github.com/emersion/go-smtp"
)

// Options configure deterministic local-only mail behavior.
type Options struct {
	TLSMode            string
	CertificateHost    string
	ExpiredCertificate bool
	Username, Password string
	Reject             map[string]int
	DataCode           int
	ReplyText          string
}

// Message is one message accepted by the owned receiver.
type Message struct {
	From       string
	Recipients []string
	Data       []byte
}

// Server owns a listener, all its SMTP sessions and bounded observed messages.
type Server struct {
	server   *smtp.Server
	listener net.Listener
	done     chan error
	mu       sync.Mutex
	options  Options
	messages []Message
	rcpts    []string
	rootPEM  []byte
}

// Start creates an owned loopback receiver with at most eight active connections.
func Start(options Options) (*Server, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	f := &Server{listener: listener, done: make(chan error, 1), options: options}
	f.options.Reject = map[string]int{}
	for a, c := range options.Reject {
		f.options.Reject[a] = c
	}
	server := smtp.NewServer(smtp.BackendFunc(func(_ *smtp.Conn) (smtp.Session, error) { return &session{owner: f}, nil }))
	server.Domain = "localhost"
	server.AllowInsecureAuth = true
	server.MaxRecipients = 100
	server.MaxMessageBytes = 256 << 10
	server.ReadTimeout = 5 * time.Second
	server.WriteTimeout = 5 * time.Second
	server.ErrorLog = log.New(io.Discard, "", 0)
	bounded := &limitedListener{Listener: listener, ready: make(chan struct{})}
	var serving net.Listener = bounded
	if options.TLSMode != "" && options.TLSMode != "plain" {
		certificate, root, err := certificate(options)
		if err != nil {
			_ = listener.Close()
			return nil, err
		}
		f.rootPEM = root
		server.TLSConfig = &tls.Config{Certificates: []tls.Certificate{certificate}, MinVersion: tls.VersionTLS12}
		if options.TLSMode == "tls" {
			serving = tls.NewListener(bounded, server.TLSConfig)
		}
	}
	f.server = server
	go func() { f.done <- server.Serve(serving) }()
	<-bounded.ready
	return f, nil
}

// Addr is the owned loopback endpoint.
func (s *Server) Addr() string { return s.listener.Addr().String() }

// Roots returns a trust pool containing only the owned fixture certificate.
func (s *Server) Roots() *x509.CertPool {
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(s.rootPEM)
	return pool
}

// RootPEM returns the public fixture trust certificate.
func (s *Server) RootPEM() []byte { return append([]byte(nil), s.rootPEM...) }

// Messages returns copies of accepted messages without credentials.
func (s *Server) Messages() []Message {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Message, len(s.messages))
	for i, m := range s.messages {
		out[i] = Message{From: m.From, Recipients: append([]string(nil), m.Recipients...), Data: append([]byte(nil), m.Data...)}
	}
	return out
}

// Recipients returns observed RCPT addresses for retry-exclusion assertions.
func (s *Server) Recipients() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.rcpts...)
}

// SetRecipientCode changes an owned address's deterministic rejection.
func (s *Server) SetRecipientCode(address string, code int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.options.Reject[address] = code
}

// Close terminates the listener and waits for every owned session to finish.
func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	err := s.server.Shutdown(ctx)
	if errors.Is(err, smtp.ErrServerClosed) {
		err = nil
	}
	<-s.done
	return err
}

type session struct {
	owner         *Server
	from          string
	recipients    []string
	authenticated bool
}

func (s *session) Reset()                   { s.from = ""; s.recipients = nil }
func (s *session) Logout() error            { return nil }
func (s *session) AuthMechanisms() []string { return []string{"PLAIN"} }
func (s *session) Auth(_ string) (sasl.Server, error) {
	return sasl.NewPlainServer(func(_, user, password string) error {
		s.owner.mu.Lock()
		defer s.owner.mu.Unlock()
		if user != s.owner.options.Username || password != s.owner.options.Password {
			return &smtp.SMTPError{Code: 535, Message: s.owner.options.ReplyText + " authentication rejected"}
		}
		s.authenticated = true
		return nil
	}), nil
}
func (s *session) Mail(from string, _ *smtp.MailOptions) error {
	s.owner.mu.Lock()
	required := s.owner.options.Username != ""
	s.owner.mu.Unlock()
	if required && !s.authenticated {
		return smtp.ErrAuthRequired
	}
	s.from = from
	return nil
}
func (s *session) Rcpt(to string, _ *smtp.RcptOptions) error {
	s.owner.mu.Lock()
	defer s.owner.mu.Unlock()
	if len(s.owner.rcpts) < 10000 {
		s.owner.rcpts = append(s.owner.rcpts, to)
	}
	code := s.owner.options.Reject[to]
	if code != 0 {
		return &smtp.SMTPError{Code: code, Message: s.owner.options.ReplyText + " recipient rejected"}
	}
	s.recipients = append(s.recipients, to)
	return nil
}
func (s *session) Data(reader io.Reader) error {
	data, err := io.ReadAll(io.LimitReader(reader, (256<<10)+1))
	if err != nil {
		return err
	}
	if len(data) > 256<<10 {
		return &smtp.SMTPError{Code: 552, Message: "message too large"}
	}
	s.owner.mu.Lock()
	defer s.owner.mu.Unlock()
	if s.owner.options.DataCode != 0 {
		return &smtp.SMTPError{Code: s.owner.options.DataCode, Message: s.owner.options.ReplyText + " submission rejected"}
	}
	if len(s.owner.messages) >= 1000 {
		return &smtp.SMTPError{Code: 452, Message: "owned message budget reached"}
	}
	s.owner.messages = append(s.owner.messages, Message{From: s.from, Recipients: append([]string(nil), s.recipients...), Data: data})
	return nil
}
func certificate(options Options) (tls.Certificate, []byte, error) {
	public, key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return tls.Certificate{}, nil, err
	}
	template := &x509.Certificate{SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "owned SMTP fixture"}, NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(24 * time.Hour), KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign, ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}, BasicConstraintsValid: true, IsCA: true, IPAddresses: []net.IP{net.ParseIP("127.0.0.1")}, DNSNames: []string{"localhost"}}
	if options.CertificateHost != "" {
		template.IPAddresses = nil
		template.DNSNames = []string{options.CertificateHost}
	}
	if options.ExpiredCertificate {
		template.NotBefore = time.Now().Add(-48 * time.Hour)
		template.NotAfter = time.Now().Add(-24 * time.Hour)
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, public, key)
	if err != nil {
		return tls.Certificate{}, nil, err
	}
	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return tls.Certificate{}, nil, err
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	cert, err := tls.X509KeyPair(certPEM, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER}))
	return cert, certPEM, err
}

type limitedListener struct {
	net.Listener
	mu     sync.Mutex
	active int
	ready  chan struct{}
	once   sync.Once
}

func (l *limitedListener) Accept() (net.Conn, error) {
	l.once.Do(func() { close(l.ready) })
	for {
		c, err := l.Listener.Accept()
		if err != nil {
			return nil, err
		}
		l.mu.Lock()
		if l.active >= 8 {
			l.mu.Unlock()
			_ = c.Close()
			continue
		}
		l.active++
		l.mu.Unlock()
		return &limitedConn{Conn: c, release: func() { l.mu.Lock(); l.active--; l.mu.Unlock() }}, nil
	}
}

type limitedConn struct {
	net.Conn
	once    sync.Once
	release func()
}

func (c *limitedConn) Close() error { err := c.Conn.Close(); c.once.Do(c.release); return err }

// SetCredentials replaces only the owned receiver authentication fixture.
func (s *Server) SetCredentials(username, password string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.options.Username, s.options.Password = username, password
}

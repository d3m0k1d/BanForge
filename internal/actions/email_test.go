package actions

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/d3m0k1d/BanForge/internal/config"
)

type simpleSMTPServer struct {
	listener net.Listener
	messages []string
	done     chan struct{}
}

func newSimpleSMTPServer(t *testing.T, useTLS bool) *simpleSMTPServer {
	t.Helper()

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}

	s := &simpleSMTPServer{
		listener: l,
		messages: make([]string, 0),
		done:     make(chan struct{}),
	}

	go s.serve(t, useTLS)

	time.Sleep(50 * time.Millisecond)

	return s
}

func (s *simpleSMTPServer) serve(t *testing.T, useTLS bool) {
	defer close(s.done)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}

		go func(c net.Conn) {
			defer c.Close()

			_, _ = c.Write([]byte("220 localhost ESMTP Test Server\r\n"))

			reader := bufio.NewReader(c)

			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					return
				}

				line = strings.TrimSpace(line)
				parts := strings.SplitN(line, " ", 2)
				cmd := strings.ToUpper(parts[0])

				switch cmd {
				case "EHLO", "HELO":
					if useTLS {
						_, _ = c.Write([]byte("250-localhost\r\n250 STARTTLS\r\n"))
					} else {
						_, _ = c.Write([]byte("250-localhost\r\n250 AUTH PLAIN\r\n"))
					}

				case "STARTTLS":
					if !useTLS {
						_, _ = c.Write([]byte("454 TLS not available\r\n"))
						return
					}
					_, _ = c.Write([]byte("220 Ready to start TLS\r\n"))

					tlsConfig := &tls.Config{
						InsecureSkipVerify: true,
						MinVersion:         tls.VersionTLS12,
					}
					tlsConn := tls.Server(c, tlsConfig)
					if err := tlsConn.Handshake(); err != nil {
						return
					}
					reader = bufio.NewReader(tlsConn)
					c = tlsConn

				case "AUTH":
					_, _ = c.Write([]byte("235 Authentication successful\r\n"))

				case "MAIL":
					_, _ = c.Write([]byte("250 OK\r\n"))

				case "RCPT":
					_, _ = c.Write([]byte("250 OK\r\n"))

				case "DATA":
					_, _ = c.Write([]byte("354 End data with <CR><LF>.<CR><LF>\r\n"))

					var msgBuilder strings.Builder
					for {
						msgLine, err := reader.ReadString('\n')
						if err != nil {
							return
						}
						if strings.TrimSpace(msgLine) == "." {
							break
						}
						msgBuilder.WriteString(msgLine)
					}
					s.messages = append(s.messages, msgBuilder.String())
					_, _ = c.Write([]byte("250 OK\r\n"))

				case "QUIT":
					_, _ = c.Write([]byte("221 Bye\r\n"))
					return

				default:
					_, _ = c.Write([]byte("502 Command not implemented\r\n"))
				}
			}
		}(conn)
	}
}

func (s *simpleSMTPServer) Addr() string {
	return s.listener.Addr().String()
}

func (s *simpleSMTPServer) Close() {
	_ = s.listener.Close()
	<-s.done
}

func (s *simpleSMTPServer) MessageCount() int {
	return len(s.messages)
}

func TestSendEmail_Validation(t *testing.T) {
	tests := []struct {
		name    string
		action  config.Action
		wantErr bool
		errMsg  string
	}{
		{
			name: "disabled action",
			action: config.Action{
				Type:        "email",
				Enabled:     false,
				Email:       "test@example.com",
				EmailSender: "sender@example.com",
				SMTPHost:    "smtp.example.com",
			},
			wantErr: false,
		},
		{
			name: "empty SMTP host",
			action: config.Action{
				Type:        "email",
				Enabled:     true,
				Email:       "test@example.com",
				EmailSender: "sender@example.com",
				SMTPHost:    "",
			},
			wantErr: true,
			errMsg:  "SMTP host is empty",
		},
		{
			name: "empty recipient email",
			action: config.Action{
				Type:        "email",
				Enabled:     true,
				Email:       "",
				EmailSender: "sender@example.com",
				SMTPHost:    "smtp.example.com",
			},
			wantErr: true,
			errMsg:  "recipient email is empty",
		},
		{
			name: "empty sender email",
			action: config.Action{
				Type:        "email",
				Enabled:     true,
				Email:       "test@example.com",
				EmailSender: "",
				SMTPHost:    "smtp.example.com",
			},
			wantErr: true,
			errMsg:  "sender email is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SendEmail(tt.action)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if err.Error() != tt.errMsg {
					t.Errorf("SendEmail() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestSendEmail_WithoutTLS(t *testing.T) {
	server := newSimpleSMTPServer(t, false)
	defer server.Close()

	host, portStr, _ := net.SplitHostPort(server.Addr())
	var port int
	fmt.Sscanf(portStr, "%d", &port)

	action := config.Action{
		Type:         "email",
		Enabled:      true,
		Email:        "recipient@example.com",
		EmailSender:  "sender@example.com",
		EmailSubject: "Test Subject",
		SMTPHost:     host,
		SMTPPort:     port,
		SMTPUser:     "user",
		SMTPPassword: "pass",
		SMTPTLS:      false,
		Body:         "Test message body",
	}

	err := SendEmail(action)
	if err != nil {
		t.Fatalf("SendEmail() unexpected error: %v", err)
	}

	if server.MessageCount() != 1 {
		t.Errorf("Expected 1 message, got %d", server.MessageCount())
	}
}

func TestSendEmail_WithTLS(t *testing.T) {
	t.Skip("TLS test requires proper TLS handshake handling")
	server := newSimpleSMTPServer(t, true)
	defer server.Close()

	host, portStr, _ := net.SplitHostPort(server.Addr())
	var port int
	fmt.Sscanf(portStr, "%d", &port)

	action := config.Action{
		Type:         "email",
		Enabled:      true,
		Email:        "recipient@example.com",
		EmailSender:  "sender@example.com",
		EmailSubject: "Test Subject TLS",
		SMTPHost:     host,
		SMTPPort:     port,
		SMTPUser:     "user",
		SMTPPassword: "pass",
		SMTPTLS:      true,
		Body:         "Test TLS message body",
	}

	err := SendEmail(action)
	if err != nil {
		t.Fatalf("SendEmail() unexpected error: %v", err)
	}

	if server.MessageCount() != 1 {
		t.Errorf("Expected 1 message, got %d", server.MessageCount())
	}
}

func TestSendEmail_DefaultSubject(t *testing.T) {
	server := newSimpleSMTPServer(t, false)
	defer server.Close()

	host, portStr, _ := net.SplitHostPort(server.Addr())
	var port int
	fmt.Sscanf(portStr, "%d", &port)

	action := config.Action{
		Type:        "email",
		Enabled:     true,
		Email:       "test@example.com",
		EmailSender: "sender@example.com",
		SMTPHost:    host,
		SMTPPort:    port,
		Body:        "Test body",
	}

	err := SendEmail(action)
	if err != nil {
		t.Fatalf("SendEmail() unexpected error: %v", err)
	}

	if server.MessageCount() != 1 {
		t.Errorf("Expected 1 message, got %d", server.MessageCount())
	}
}

func TestSendEmail_Integration(t *testing.T) {
	server := newSimpleSMTPServer(t, false)
	defer server.Close()

	host, portStr, _ := net.SplitHostPort(server.Addr())
	var port int
	fmt.Sscanf(portStr, "%d", &port)

	action := config.Action{
		Type:         "email",
		Enabled:      true,
		Email:        "to@example.com",
		EmailSender:  "from@example.com",
		EmailSubject: "Integration Test",
		SMTPHost:     host,
		SMTPPort:     port,
		SMTPUser:     "testuser",
		SMTPPassword: "testpass",
		SMTPTLS:      false,
		Body:         "Integration test message",
	}

	err := SendEmail(action)
	if err != nil {
		t.Fatalf("SendEmail() failed: %v", err)
	}

	t.Logf("Email sent successfully, server received %d message(s)", server.MessageCount())
}

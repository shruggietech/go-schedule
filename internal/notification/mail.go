package notification

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"html"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"net/textproto"
	"strings"
	"time"
	"unicode"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

const maxMailBytes = 256 << 10

// RenderMail creates a bounded UTF-8 MIME message from the shared notification event.
func RenderMail(d domain.NotificationDelivery, c domain.SMTPConfig) ([]byte, error) {
	if len(d.Payload) > 128<<10 {
		return nil, fmt.Errorf("notification event is too large to email")
	}
	var event domain.WebhookEvent
	if err := json.Unmarshal(d.Payload, &event); err != nil {
		return nil, fmt.Errorf("notification event cannot be read")
	}
	title := "go-schedule: test notification"
	body := "A test notification was requested.\n"
	if event.Event == string(domain.NotificationEventDaemonHealth) {
		title = "go-schedule: daemon healthy"
		body = "The scheduler daemon is running.\n"
		if event.Daemon.NextExpectedAt != nil {
			body += "Next expected heartbeat: " + event.Daemon.NextExpectedAt.Format(time.RFC3339) + "\n"
		}
	} else if event.Task != nil {
		title = "go-schedule: " + event.Task.Name
		body = "Task: " + event.Task.Name + "\n"
		if event.Task.GroupName != "" {
			body += "Group: " + event.Task.GroupName + "\n"
		}
		if event.Run != nil {
			title += " (" + string(event.Run.Outcome) + ")"
			body += "Outcome: " + string(event.Run.Outcome) + "\nRun: " + event.Run.ID + "\n"
		}
	}
	if event.Condition != nil {
		body += "Condition: " + event.Condition.Summary + "\n"
	}
	body += "Delivery: " + d.ID + "\nCreated: " + d.CreatedAt.UTC().Format(time.RFC3339) + "\n"
	title = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return ' '
		}
		return r
	}, title)
	if runes := []rune(title); len(runes) > 180 {
		title = string(runes[:180])
	}
	from, err := mail.ParseAddress(c.From)
	if err != nil {
		return nil, fmt.Errorf("email sender address is invalid")
	}
	recipients := make([]string, 0, len(c.Recipients))
	for _, address := range c.Recipients {
		parsed, err := mail.ParseAddress(address)
		if err != nil {
			return nil, fmt.Errorf("email recipient address is invalid")
		}
		recipients = append(recipients, parsed.String())
	}
	var out bytes.Buffer
	fmt.Fprintf(&out, "From: %s\r\n", from.String())
	out.WriteString("To: " + strings.Join(recipients, ",\r\n ") + "\r\n")
	subject := mime.QEncoding.Encode("utf-8", title)
	out.WriteString("Subject: " + strings.ReplaceAll(subject, "?= =?", "?=\r\n =?") + "\r\n")
	fmt.Fprintf(&out, "Date: %s\r\nMessage-ID: <%x@go-schedule.local>\r\nMIME-Version: 1.0\r\n", d.CreatedAt.UTC().Format(time.RFC1123Z), sha256.Sum256([]byte(d.ID)))
	if !c.HTML {
		out.WriteString("Content-Type: text/plain; charset=UTF-8\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\n")
		writer := quotedprintable.NewWriter(&out)
		if _, err := writer.Write([]byte(body)); err != nil {
			return nil, fmt.Errorf("encode email text")
		}
		if err := writer.Close(); err != nil {
			return nil, fmt.Errorf("finish email text")
		}
	} else {
		writer := multipart.NewWriter(&out)
		fmt.Fprintf(&out, "Content-Type: multipart/alternative; boundary=%q\r\n\r\n", writer.Boundary())
		for _, part := range []struct{ kind, body string }{{"text/plain", body}, {"text/html", "<!doctype html><html><body><pre>" + html.EscapeString(body) + "</pre></body></html>"}} {
			header := textproto.MIMEHeader{"Content-Type": {part.kind + "; charset=UTF-8"}, "Content-Transfer-Encoding": {"quoted-printable"}}
			destination, err := writer.CreatePart(header)
			if err != nil {
				return nil, fmt.Errorf("create email part")
			}
			encoded := quotedprintable.NewWriter(destination)
			if _, err := encoded.Write([]byte(part.body)); err != nil {
				return nil, fmt.Errorf("encode email part")
			}
			if err := encoded.Close(); err != nil {
				return nil, fmt.Errorf("finish email part")
			}
		}
		if err := writer.Close(); err != nil {
			return nil, fmt.Errorf("finish email message")
		}
	}
	if out.Len() > maxMailBytes {
		return nil, fmt.Errorf("email message exceeds 256 KiB")
	}
	return out.Bytes(), nil
}

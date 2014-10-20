// vim: tabstop=2 shiftwidth=2

package main

import (
	"os"
	"strings"
	"net/mail"
	"bytes"
	"fmt"
	"path"
	"encoding/hex"
	"net/smtp"
	"io/ioutil"
	"errors"
	"crypto/sha256"
)

func assemble(msg mail.Message) []byte {
	buf := new(bytes.Buffer)
	for h := range msg.Header {
		buf.WriteString(h + ": " + msg.Header.Get(h) + "\n")
	}
	buf.WriteString("\n")
	buf.ReadFrom(msg.Body)
	return buf.Bytes()
}

// testMail vets outbound messages to final recipients
func testMail(b []byte) (recipients []string, err error) {
	f := bytes.NewReader(b)
	msg, err := mail.ReadMessage(f)
	if err != nil {
		Trace.Printf("Outbound read failure: %s", err)
		return
	}
	var exists bool
	h := msg.Header
	_, exists = h["To"]
	if ! exists {
		err = errors.New("No recipient specified in final delivery")
		Trace.Println(err)
		return
	}
	addys, err := h.AddressList("To")
	if err != nil {
		return
	}
	for _, addy := range addys {
		recipients = append(recipients, addy.Address)
	}
	return
}

func mailFile(filename string) (err error) {
	var payload []byte
	payload, err = ioutil.ReadFile(filename)
	if err != nil {
		Error.Printf("Failed to read file for mailing: %s", err)
		return
	}
	sendTo := strings.TrimRight(string(payload[:80]), "\x00")
	payload = payload[80:]
	// Test if the message is destined for the local remailer
	Trace.Printf("Message recipient is: %s", sendTo)
	if cfg.Mail.Outfile {
		var f *os.File
		digest := sha256.New()
		digest.Write(payload)
		filename := "outfile-" + hex.EncodeToString(digest.Sum(nil)[:16])
		f, err = os.Create(path.Join(cfg.Files.Pooldir, filename))
		defer f.Close()
		_, err = f.WriteString(string(payload))
		if err != nil {
			Warn.Printf("Outfile write failed: %s\n", err)
			return
		}
	} else if cfg.Mail.Sendmail {
		err = sendmail(payload, sendTo)
		if err != nil {
			Warn.Println("Sendmail failed")
			return
		}
	} else {
		err = SMTPRelay(payload, sendTo)
		if err != nil {
			Warn.Println("SMTP relay failed")
			return
		}
	}
	return
}

func SMTPRelay(payload []byte, sendto string) (err error) {
	c, err := smtp.Dial(fmt.Sprintf("%s:%d", cfg.Mail.SMTPRelay, cfg.Mail.SMTPPort))
	if err != nil {
		Warn.Println(err)
		return
	}
	err = c.Mail(cfg.Mail.EnvelopeSender)
	if err != nil {
		Warn.Println(err)
		return
	}
	err = c.Rcpt(sendto)
	if err != nil {
		Warn.Println(err)
		return
	}
	wc, err := c.Data()
	if err != nil {
		Warn.Println(err)
		return
	}
	_, err = fmt.Fprintf(wc, string(payload))
	if err != nil {
		Warn.Println(err)
		return
	}
	err = wc.Close()
	if err != nil {
		Warn.Println(err)
		return
	}
	err = c.Quit()
	if err != nil {
		Warn.Println(err)
		return
	}
	return
}

// sendmail invokes go's sendmail method
func sendmail(payload []byte, sendto string) (err error) {
	auth := smtp.PlainAuth("", cfg.Mail.SMTPUsername, cfg.Mail.SMTPPassword, cfg.Mail.SMTPRelay)
	relay := fmt.Sprintf("%s:%d", cfg.Mail.SMTPRelay, cfg.Mail.SMTPPort)
	err = smtp.SendMail(relay, auth, cfg.Mail.EnvelopeSender, []string{sendto}, payload)
	if err != nil {
		Warn.Println(err)
		return
	}
	return
}

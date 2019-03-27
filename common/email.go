package common

import (
	"bytes"
	"encoding/base64"
	"errors"
	"mime/multipart"
	"net/textproto"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

func buildEmailInput(recipient string, png []byte) (*ses.SendRawEmailInput, error) {
	source := os.Getenv("SENDER_EMAIL_ADDRESS")
	message := "Scan this QR code, quickly"
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	header := make(textproto.MIMEHeader)
	header.Set("From", source)
	header.Set("To", recipient)
	header.Set("Return-Path", source)
	header.Set("Subject", "Your QR code")
	header.Set("Content-Language", "en-US")
	header.Set("Content-Type", "multipart/mixed; boundary=\""+writer.Boundary()+"\"")
	header.Set("MIME-Version", "1.0")

	_, err := writer.CreatePart(header)
	if err != nil {
		return nil, err
	}

	header = make(textproto.MIMEHeader)
	header.Set("Content-Transfer-Encoding", "7bit")
	header.Set("Content-Type", "text/plain; charset=us-ascii")

	part, err := writer.CreatePart(header)
	if err != nil {
		return nil, err
	}

	_, err = part.Write([]byte(message))
	if err != nil {
		return nil, err
	}

	filename := "qrcode.png"
	header = make(textproto.MIMEHeader)
	header.Set("Content-Disposition", "attachment; filename="+filename)
	header.Set("Content-Type", "image/png; x-unix-mode=0644; name=\""+filename+"\"")
	header.Set("Content-Transfer-Encoding", "base64")

	part, err = writer.CreatePart(header)
	if err != nil {
		return nil, err
	}

	encodedPng := base64.StdEncoding.EncodeToString(png)
	_, err = part.Write([]byte(encodedPng))
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	s := buf.String()
	if strings.Count(s, "\n") < 2 {
		return nil, errors.New("invalid e-mail content")
	}
	s = strings.SplitN(s, "\n", 2)[1]

	raw := ses.RawMessage{
		Data: []byte(s),
	}
	input := &ses.SendRawEmailInput{
		Destinations: []*string{aws.String(recipient)},
		Source:       aws.String(source),
		RawMessage:   &raw,
	}

	return input, nil
}

func SendEmail(email string, png []byte) error {
	svc := ses.New(session.New())

	var err error
	var emailInput *ses.SendRawEmailInput
	emailInput, err = buildEmailInput(email, png)
	if err != nil {
		return err
	}

	_, err = svc.SendRawEmail(emailInput)
	if err != nil {
		return err
	}

	return err
}

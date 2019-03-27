package common

import (
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
	"time"

	wa "github.com/Rhymen/go-whatsapp"
	qrcode "github.com/skip2/go-qrcode"
)

func SignIn(number string, email string) error {
	con, err := wa.NewConn(5 * time.Second)
	if err != nil {
		return err
	}

	qr := make(chan string)
	go func() {
		var png []byte
		png, err := qrcode.Encode(<-qr, qrcode.Medium, 256)
		if err != nil {
			log.Println(err)
		}

		err = SendEmail(email, png)
		if err != nil {
			log.Println(err)
		}
	}()

	session, err := con.Login(qr)
	if err != nil {
		return err
	}

	err = writeSession(number, session)
	if err != nil {
		return err
	}

	return nil
}

func sessionTempPath(sender string) string {
	return os.TempDir() + "/" + sender + ".gob"
}

func writeSession(number string, session wa.Session) error {
	sender := Hash(number)
	path := sessionTempPath(sender)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(session)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	prefix := os.Getenv("SESSION_SECRETS_MANAGER_PREFIX")
	key := prefix + sender
	err = UpdateSecret(key, data)

	return err
}

func readSession() (wa.Session, error) {
	sender := Hash("16464229653")
	prefix := os.Getenv("SESSION_SECRETS_MANAGER_PREFIX")
	key := prefix + sender
	var session wa.Session
	result, err := GetSecret(key)
	if err != nil {
		return session, err
	}

	path := sessionTempPath(sender)
	err = ioutil.WriteFile(path, result.SecretBinary, 0644)
	if err != nil {
		return session, err
	}

	file, err := os.Open(path)
	if err != nil {
		return session, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)

	err = decoder.Decode(&session)
	if err != nil {
		return session, err
	}

	return session, nil
}

func SendMessage(recipient string, message string) error {
	con, err := wa.NewConn(10 * time.Second)
	if err != nil {
		return err
	}

	session, err := readSession()

	if err != nil {
		return err
	}

	session, err = con.RestoreSession(session)
	if err != nil {
		return err
	}

	err = writeSession("16464229653", session)
	if err != nil {
		return err
	}

	<-time.After(3 * time.Second)

	msg := wa.TextMessage{
		Info: wa.MessageInfo{
			RemoteJid: recipient + "@s.whatsapp.net",
		},
		Text: message,
	}

	err = con.Send(msg)
	if err != nil {
		return err
	}

	err = writeSession("16464229653", session)
	if err != nil {
		return err
	}

	return err
}

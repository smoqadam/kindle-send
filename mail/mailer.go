package mail

import (
	"fmt"
	"os"
	"time"

	config "github.com/nikhil1raghav/kindle-send/config"
	"github.com/nikhil1raghav/kindle-send/util"
	gomail "gopkg.in/mail.v2"
)

func Send(files []string, timeout int) error {
	if len(files) == 0 {
		return fmt.Errorf("no files provided to send")
	}

	cfg := config.GetInstance()
	if cfg == nil {
		return fmt.Errorf("failed to get config instance")
	}

	if cfg.Sender == "" || cfg.Receiver == "" || cfg.Server == "" || cfg.Password == "" {
		return fmt.Errorf("incomplete email configuration")
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", cfg.Sender)
	msg.SetHeader("To", cfg.Receiver)
	msg.SetBody("text/plain", "")

	attachedFiles := make([]string, 0, len(files))
	for _, file := range files {
		if file == "" {
			util.Red.Println("Skipping empty file path")
			continue
		}

		fileInfo, err := os.Stat(file)
		if err != nil {
			util.Red.Printf("Couldn't find the file %s: %v\n", file, err)
			continue
		}

		if fileInfo.Size() == 0 {
			util.Red.Printf("Skipping empty file: %s\n", file)
			continue
		}

		msg.Attach(file)
		attachedFiles = append(attachedFiles, file)
	}

	if len(attachedFiles) == 0 {
		return fmt.Errorf("no valid files to send")
	}

	if timeout <= 0 {
		return fmt.Errorf("invalid timeout value: %d", timeout)
	}

	dialer := gomail.NewDialer(cfg.Server, cfg.Port, cfg.Sender, cfg.Password)
	dialer.Timeout = time.Duration(timeout) * time.Second

	util.CyanBold.Println("Sending mail")
	util.Cyan.Printf("Mail timeout: %s\n", dialer.Timeout.String())
	util.Cyan.Println("Following files will be sent:")

	for i, file := range attachedFiles {
		util.Cyan.Printf("%d. %s\n", i+1, file)
	}

	if err := dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send mail: %v", err)
	}

	util.GreenBold.Printf("Successfully mailed %d files to %s\n", len(attachedFiles), cfg.Receiver)
	return nil
}

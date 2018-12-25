package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/veqryn/go-email/email"
)

var opts struct {
	Verbose   bool   `short:"v" long:"verbose" description:"Show verbose debug information"`
	MessageId bool   `long:"message-id" description:"Message-ID"`
	From      bool   `long:"from" description:"From"`
	FromEmail bool   `long:"from-email" description:"From (email)"`
	To        bool   `long:"to" description:"To"`
	ToEmail   bool   `long:"to-email" description:"To (email)"`
	Subject   bool   `long:"subject" description:"Subject"`
	EmailFile string `short:"f" long:"email-file" description:"Path to email file" required:"true"`
}

func debug(message string, params ...interface{}) {
	if opts.Verbose {
		fmt.Printf(message+"\n", params...)
	}
}

func email_file_to_msg(file_path string) *email.Message {
	var reader io.Reader
	if file_path == "-" {
		reader = bufio.NewReader(os.Stdin)
	} else {
		file, _ := ioutil.ReadFile(file_path)
		reader = strings.NewReader(string(file))
	}
	msg, _ := email.ParseMessage(reader)

	return msg
}

// this is lazy way to decode subject,
// it will break at some emails, I'm sure
func string_decode(name string) string {
	re := regexp.MustCompile("^=\\?[a-zA-Z0-9_\\-]*\\?.\\?(.*)\\?=")
	newName := re.ReplaceAllString(name, "$1")
	if newName != name {
		re = regexp.MustCompile("=([A-F0-9][A-F0-9])")
		newName = re.ReplaceAllString(newName, "%$1")
		newName, _ = url.PathUnescape(newName)
	}
	return newName
}

func get_message_id(msg *email.Message) string {
	messageId := msg.Header.Get("Message-ID")
	re := regexp.MustCompile("<(.*)>")
	cleanId := re.ReplaceAllString(messageId, "$1")

	return cleanId
}

func get_from(msg *email.Message) string {
	return msg.Header.From()
}

func get_to(msg *email.Message) []string {
	return msg.Header.To()
}

func get_from_clean(msg *email.Message) string {
	from := get_from(msg)
	re := regexp.MustCompile(".*<(.*)>")
	clean := re.ReplaceAllString(from, "$1")

	return clean
}

func get_to_clean(msg *email.Message) []string {
	to := get_to(msg)
	re := regexp.MustCompile(".*<(.*)>")
	var clean []string
	for _, to := range to {
		clean = append(clean, re.ReplaceAllString(to, "$1"))
	}

	return clean
}
func get_subject(msg *email.Message) string {
	return string_decode(msg.Header.Get("Subject"))
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		return
	}

	msg := email_file_to_msg(opts.EmailFile)

	if opts.MessageId {
		fmt.Println(get_message_id(msg))
	}
	if opts.From {
		fmt.Println(get_from(msg))
	}
	if opts.FromEmail {
		fmt.Println(get_from_clean(msg))
	}
	if opts.To {
		fmt.Println(get_to(msg))
	}
	if opts.ToEmail {
		fmt.Println(get_to_clean(msg))
	}
	if opts.Subject {
		fmt.Println(get_subject(msg))
	}
}

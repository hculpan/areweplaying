package email

import (
	"fmt"
	"io"
	"log"
	"net/mail"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func FetchEmails() ([]Email, error) {
	config, err := NewConfigFromEnv()
	if err != nil {
		return nil, err
	}

	imapClient, err := config.ConnectToServer(IMAP_SERVICE)
	if err != nil {
		log.Fatal(err)
	}
	defer imapClient.Logout()

	return fetchEmails(imapClient)
}

func concatAddresses(addresses []*imap.Address) string {
	result := []string{}

	for _, address := range addresses {
		if len(address.PersonalName) > 0 {
			result = append(result, fmt.Sprintf("%s <%s>", address.PersonalName, address.Address()))
		} else {
			result = append(result, address.Address())
		}

	}

	return strings.Join(result, "; ")
}

func fetchEmails(imapClient *client.Client) ([]Email, error) {
	result := []Email{}

	// Select the mailbox you want to read
	_, err := imapClient.Select("INBOX", false)
	if err != nil {
		return nil, err
	}

	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}
	uids, err := imapClient.Search(criteria)
	if err != nil {
		log.Printf("Search error: %s\n", err)
	}

	if len(uids) > 0 {
		// Define the range of emails to fetch
		seqSet := new(imap.SeqSet)
		seqSet.AddNum(uids...)

		// Fetch the required message attributes
		messages := make(chan *imap.Message, 10)
		section := &imap.BodySectionName{}
		items := []imap.FetchItem{section.FetchItem(), imap.FetchEnvelope}

		go func() {
			if err := imapClient.Fetch(seqSet, items, messages); err != nil {
				log.Fatal("Fetch error: " + err.Error())
			}
		}()

		for msg := range messages {
			toAddress := concatAddresses(msg.Envelope.To)
			fromAddress := msg.Envelope.From[0].Address()
			r := msg.GetBody(section)
			if r == nil {
				return result, fmt.Errorf("server didn't returned message body")
			}
			m, err := mail.ReadMessage(r)
			if err != nil {
				return result, err
			}
			body, err := io.ReadAll(m.Body)

			email := NewEmail(
				[]string{toAddress},
				fromAddress,
				msg.Envelope.Subject,
				string(body),
			)
			result = append(result, *email)
		}
	}

	return result, nil
}

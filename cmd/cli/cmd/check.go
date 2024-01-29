/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/hculpan/areweplaying/pkg/data"
	"github.com/hculpan/areweplaying/pkg/email"
	"github.com/spf13/cobra"
)

func isAttending(body string) string {
	body = strings.ToLower(body)
	lines := strings.Fields(body)

	for _, l := range lines {
		line := strings.Trim(l, " \t\r\n")
		if strings.HasPrefix(line, "yes") {
			return "yes"
		} else if strings.HasPrefix(line, "no") {
			return "no"
		}
	}

	return "unknown"
}

func getGmEmail(session data.Session) *data.Player {
	for i, p := range session.Players {
		if p.Gm {
			return &session.Players[i]
		}
	}

	return nil
}

func sendEventNotice(session data.Session, subject string, body string) error {
	gm := getGmEmail(session)
	if gm == nil {
		return fmt.Errorf("unable to find GM address")
	}

	e := email.NewEmail([]string{gm.Email}, "culpanvtt@gmail.com", subject, body)
	if e == nil {
		return fmt.Errorf("unable to construct Email object")
	}

	return email.SendEmail(e)
}

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Checks for new emails on player status",
	Long: `Checks at culpanvtt@gmail.com for any new
emails sent by players to indicate if they
will attend the next session.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Default().Println("Checking for emails")
		emails, err := email.FetchEmails()
		if err != nil {
			return err
		}

		if len(emails) == 0 {
			log.Default().Println("No new emails")
		} else {
			log.Default().Println("Reading new emails")
			session, err := data.ReadSessionData()
			if err != nil {
				log.Default().Println(err)
				return err
			}

			for _, email := range emails {
				if len(email.Body) > 0 {
					player := session.GetPlayerByEmail(email.FromAddress)
					if player != nil {
						attending := isAttending(email.Body)
						if attending == "yes" || attending == "no" {
							player.Attending = attending
							log.Default().Printf("Received player response for %s: %s", player.Name, player.Attending)
							if !player.Gm {
								if err := sendEventNotice(
									session,
									"Received player response",
									fmt.Sprintf(`Player %s has responded with %q`, player.Name, player.Attending),
								); err != nil {
									log.Default().Printf("ERROR: %s", err)
								}
							}
						} else {
							if !player.Gm {
								if err := sendEventNotice(
									session,
									"Received invalid player response",
									fmt.Sprintf(`Player %s has responded with an invalid response:

%s`, player.Name, email.Body),
								); err != nil {
									log.Default().Printf("ERROR: %s", err)
								}
							}
						}
					}
				}
			}

			data.PersistSession(session)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

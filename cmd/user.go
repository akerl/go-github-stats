package cmd

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/akerl/go-github-stats/stats"
	"github.com/spf13/cobra"
)

var userTemplate = template.Must(template.New("name").Parse(
	`Contribution data for {{.Name}}:
Today's score: {{.Today}}
Current streak: {{len .Streak}}
Longest streak: {{len .LongestStreak}}
High score: {{.Max.Score}} on {{.Max.Date}}
Quartile boundaries: {{.QuartileBoundaries}}`,
))

var userCmd = &cobra.Command{
	Use:   "user [NAME]",
	Short: "Show stats for user",
	RunE:  userRunner,
}

func init() {
	rootCmd.AddCommand(userCmd)
}

func userRunner(_ *cobra.Command, args []string) error {
	var name string
	if len(args) > 0 {
		name = args[0]
	}
	user, err := stats.LookupUser(name)
	if err != nil {
		return err
	}
	var output bytes.Buffer
	err = userTemplate.Execute(&output, user)
	if err != nil {
		return err
	}
	fmt.Println(output.String())
	return nil
}

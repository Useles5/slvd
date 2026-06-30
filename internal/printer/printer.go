package printer

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"
	"unicode/utf8"

	"github.com/Useles5/slvd/internal/models"
	"github.com/atotto/clipboard"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func PrintTable(solved []models.Submission, modeStr string, totalFetched, processed int, asMarkdown bool, copyToClipboard bool) {
	if len(solved) == 0 {
		fmt.Println("No problems found")
		return
	}

	caser := cases.Title(language.English)
	zoneName, _ := time.Now().Local().Zone()
	platformsMap := make(map[string]struct{})

	var finalOutput strings.Builder

	if asMarkdown {
		finalOutput.WriteString(fmt.Sprintf("| PLATFORM | ID | PROBLEM NAME | TIME (%s) |\n", zoneName))
		finalOutput.WriteString(fmt.Sprintf("|---|---|---|---|\n"))

		for _, sub := range solved {

			platformsMap[sub.Platform] = struct{}{}
			platform := caser.String(strings.ToLower(sub.Platform))
			timeStr := sub.SubmittedAt.Local().Format("02 Jan 06 15:04")

			finalOutput.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", platform, sub.ProblemKey, sub.ProblemName, timeStr))
		}
		finalOutput.WriteString(fmt.Sprintf("\n**SUMMARY:** %d solved | %d platform/s | %d fetched | Mode: %s\n\n", len(solved), len(platformsMap), totalFetched, modeStr))

	} else {
		var buf bytes.Buffer
		w := tabwriter.NewWriter(&buf, 0, 15, 3, ' ', 0)

		fmt.Fprintf(w, "PLATFORM\tID\tPROBLEM NAME\tTIME (%s)\n", zoneName)

		for _, sub := range solved {

			platformsMap[sub.Platform] = struct{}{}

			platform := caser.String(strings.ToLower(sub.Platform))
			name := sub.ProblemName
			if len(name) > 35 {
				name = name[:32] + "..."
			}
			timeStr := sub.SubmittedAt.Local().Format("02 Jan 06 15:04")

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", platform, sub.ProblemKey, name, timeStr)
		}

		w.Flush()

		lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
		maxLen := 0
		for _, line := range lines {
			l := utf8.RuneCountInString(line)
			maxLen = max(maxLen, l)
		}

		separator := strings.Repeat("-", maxLen)

		fmt.Println()
		if len(lines) > 0 {
			finalOutput.WriteString(lines[0] + "\n")
			finalOutput.WriteString(separator + "\n")

			for i := 1; i < len(lines); i++ {
				finalOutput.WriteString(lines[i] + "\n")
			}

			finalOutput.WriteString(separator + "\n")
		}

		finalOutput.WriteString(fmt.Sprintf("\nSUMMARY: %d solved | %d platform/s | %d fetched | Mode: %s\n", len(solved), len(platformsMap), totalFetched, modeStr))
	}

	fmt.Println(finalOutput.String())

	if copyToClipboard {
		err := clipboard.WriteAll(finalOutput.String())
		if err != nil {
			fmt.Fprintln(os.Stderr, "\n Failed to copy to clipboard: %v\n", err)
		} else {
			fmt.Fprintln(os.Stderr, "\n Successfully copied to clipboard\n")
		}
	}
}

package printer

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"
	"time"
	"unicode/utf8"

	"github.com/Useles5/slvd/internal/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func PrintTable(solved []models.Submission, modeStr string, totalFetched, processed int) {
	if len(solved) == 0 {
		fmt.Println("No problems found")
		return
	}

	var buf bytes.Buffer

	w := tabwriter.NewWriter(&buf, 0, 15, 3, ' ', 0)
	caser := cases.Title(language.English)
	zoneName, _ := time.Now().Local().Zone()

	fmt.Fprintf(w, "PLATFORM\tID\tPROBLEM NAME\tTIME (%s)\n", zoneName)

	platformsMap := make(map[string]struct{})

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
		fmt.Println(lines[0])
		fmt.Println(separator)

		for i := 1; i < len(lines); i++ {
			fmt.Println(lines[i])
		}

		fmt.Println(separator)
	}

	fmt.Printf("SUMMARY: %d solved | %d platform/s | %d fetched | Mode: %s\n\n",
		len(solved), len(platformsMap), totalFetched, modeStr)
}

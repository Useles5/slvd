package models

import "time"

type Submission struct {
	Platform    string
	ProblemKey  string
	ProblemName string
	IsAccepted  bool
	SubmittedAt time.Time
}

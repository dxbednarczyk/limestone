package web

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type Track struct {
	Name      string `json:"title"`
	Performer struct {
		Name string `json:"name"`
	} `json:"performer"`
	Duration        int  `json:"duration"`
	ParentalWarning bool `json:"parental_warning"`
	ID              int  `json:"id"`
}

func (t Track) Title() string {
	var builder strings.Builder

	if t.ParentalWarning {
		builder.WriteString("[E] ")
	}

	builder.WriteString(t.Name)

	return builder.String()
}

func (t Track) Description() string { return fmt.Sprintf("%s | %s", t.Performer.Name, t.FormatTime()) }
func (t Track) FilterValue() string { return t.Name }

func (t Track) FormatTime() string {
	duration := time.Duration(t.Duration) * time.Second

	minutes := math.Floor(duration.Minutes())
	seconds := math.Floor(duration.Seconds()) - (minutes * 60)

	if seconds == 0 {
		return fmt.Sprintf("%.0f minutes", minutes)
	}

	return fmt.Sprintf("%.0f minutes, %.0f seconds", minutes, seconds)
}

package template

import (
	"os"
	"testing"
	"time"
)

func TestCreateTemplate(t *testing.T) {
	templateName := "post.tmpl"
	tags := []string{"concert", "project"}
	categories := []string{"concert", "project"}
	timeLayout := "2006-01-02T15:04:05-07:00"
	date1, err := time.Parse(timeLayout, "2023-11-11T15:00:00.000+01:00")
	detail1 := &Details{
		Date: date1,
		Venue: Venue{
			Name:    "Venue 1",
			Address: "fake address 78, 2347 XC",
			City:    "New York City",
		},
	}
	date2, err := time.Parse(timeLayout, "2023-11-12T15:00:00.000+01:00")
	detail2 := &Details{
		Date: date2,
		Venue: Venue{
			Name:    "Venue 2",
			Address: "fake address 78, 2347 XC",
			City:    "Bolzano",
		},
	}
	details := []*Details{detail1, detail2}
	post := Post{
		ID:           "eeb172bd-b61a-4469-8cca-b09db1aa16e6",
		Name:         "Concert tite",
		Description:  "some Sort of description",
		Author:       "Author",
		Ensemble:     "Tallis Scholars",
		Tags:         tags,
		Categories:   categories,
		ImageName:    "cover.jpg",
		BuyTicketUrl: "https://ticket.example.com",
		FirstDate:    date1,
		Details:      details,
	}
	tmpl, err := createTemplate(templateName)
	if err != nil {
		t.Error(err)
	}
	err = tmpl.ExecuteTemplate(os.Stdout, "post", post)
	if err != nil {
		t.Error(err)
	}
}

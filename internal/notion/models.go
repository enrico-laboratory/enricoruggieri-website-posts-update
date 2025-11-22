package notion

import "github.com/enrico-laboratory/notion-api-personal-client/client/models/parsedmodels"

type Concert struct {
	Date     *parsedmodels.Task
	Location *parsedmodels.Location
}
type Project struct {
	Project        *parsedmodels.MusicProject
	Dates          []*Concert
	GDriveImageUrl string
}

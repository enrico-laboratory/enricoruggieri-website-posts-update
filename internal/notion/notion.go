package notion

import (
	notion "github.com/enrico-laboratory/notion-api-personal-client/client"
	"github.com/enrico-laboratory/notion-api-personal-client/client/models/parsedmodels"
)

type NotiontClient struct {
	client *notion.NotionApiClient
}

func NewNotionClient(token string) (*NotiontClient, error) {

	client, err := notion.NewClient(token)
	if err != nil {
		return nil, err
	}
	return &NotiontClient{
		client: client,
	}, nil
}

func (notion *NotiontClient) GetProjects() ([]parsedmodels.MusicProject, error) {
	projects, err := notion.client.MusicProjects.GetAll()
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (notion *NotiontClient) GetConcerts(projectID string) ([]parsedmodels.Task, error) {
	concerts, err := notion.client.Schedule.GetByProjectIdAndType(projectID, "Concert")
	if err != nil {
		return nil, err
	}
	return concerts, nil
}

func (notion *NotiontClient) GetLocations() ([]parsedmodels.Location, error) {
	locations, err := notion.client.Locations.Query("")
	if err != nil {
		return nil, err
	}

	return locations, nil
}

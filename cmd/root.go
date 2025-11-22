package cmd

import (
	"errors"
	"fmt"
	"github.com/enrico-laboratory/notion-api-personal-client/client/models/parsedmodels"
	"github.com/enrico-laboratory/website-update/internal/config"
	"github.com/enrico-laboratory/website-update/internal/helpers"
	"github.com/enrico-laboratory/website-update/internal/notion"
	"github.com/enrico-laboratory/website-update/internal/template"
	"log"
	"os"
	"time"

	//"github.com/go-git/go-git/v5/plumbing/object"
	//"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/spf13/cobra"
)

const (
	imageName    = "poster.png"
	postFileName = "index.md"
)

var (
	repoURL          string
	clonePath        string
	branchName       string
	postsPath        string
	imagePath        string
	templateFileName string
)

var rootCmd = &cobra.Command{
	Use:   "update",
	Short: "Update website with new concerts",
	Run: func(cmd *cobra.Command, args []string) {
		pipeline()
	},
}

func init() {
	rootCmd.Flags().StringVar(&repoURL, "repo", "https://github.com/enrico-laboratory/enricoruggieri.com", "Git repository URL (required)")
	rootCmd.Flags().StringVar(&clonePath, "clone-path", "/tmp/repo", "Local clone path")
	rootCmd.Flags().StringVar(&branchName, "branch", "main", "Branch to use for testing before deploying")
	rootCmd.Flags().StringVar(&postsPath, "posts-path", "content/posts", "Posts path")
	rootCmd.Flags().StringVar(&imagePath, "images-path", "static/images", "Images path")
	rootCmd.Flags().StringVar(&templateFileName, "template-file-name", "post.tmpl", "Name of the template file")

	//err := rootCmd.MarkFlagRequired("repo")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err = rootCmd.MarkFlagRequired("add")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err = rootCmd.MarkFlagRequired("content")
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setConcertsForProject(nConcert *parsedmodels.Task, nLocations []parsedmodels.Location, concerts *[]*notion.Concert) {
	concert := &notion.Concert{
		Date: nConcert,
	}

	// Build a lookup table
	locationMap := make(map[string]*parsedmodels.Location)
	for i := range nLocations {
		location := &nLocations[i]
		locationMap[location.Id] = location
	}

	// Find matching location
	for _, id := range nConcert.LocationId {
		if loc := locationMap[id]; loc != nil {
			concert.Location = loc
		}
	}

	// Append using the pointer
	*concerts = append(*concerts, concert)
}

func setTemplateDataDetails(dates []*notion.Concert) ([]*template.Details, error) {
	var postTemplateDetails []*template.Details
	for _, concert := range dates {
		var postTemplateDateVenue template.Venue
		if concert.Location != nil {
			postTemplateDateVenue = template.Venue{
				Name:    concert.Location.Location,
				Address: concert.Location.Address,
				City:    concert.Location.City,
			}
		}
		timeLayout := "2006-01-02T15:04:05-07:00"
		date, err := time.Parse(timeLayout, concert.Date.StartDateAndTime)
		if err != nil {
			timeLayout := "2006-01-02"
			date, err = time.Parse(timeLayout, concert.Date.StartDateAndTime)
			if err != nil {
				return nil, err
			}
		}
		postTemplateDataDetail := &template.Details{
			Date:  date,
			Venue: postTemplateDateVenue,
		}
		postTemplateDetails = append(postTemplateDetails, postTemplateDataDetail)
	}
	return postTemplateDetails, nil
}

func setPostTemplateData(project *notion.Project, postTemplateDetails []*template.Details) template.Post {
	ImageName := ""
	if project.GDriveImageUrl != "" {
		ImageName = imageName
	}
	buyTicketUrl := ""
	if project.Project.Ticket != "" {
		buyTicketUrl = project.Project.Ticket
	}
	return template.Post{
		ID:           project.Project.Id,
		Name:         project.Project.Title,
		Description:  project.Project.Description,
		Author:       "Author",
		Ensemble:     project.Project.ChoirRollup,
		Tags:         []string{"concerts"},
		Categories:   []string{"concerts"},
		ImageName:    ImageName,
		BuyTicketUrl: buyTicketUrl,
		FirstDate:    postTemplateDetails[0].Date,
		Details:      postTemplateDetails,
	}
}

func setProjects(nProjects []parsedmodels.MusicProject, nc *notion.NotiontClient, nLocations []parsedmodels.Location) ([]*notion.Project, error) {
	var projects []*notion.Project
	for _, nProject := range nProjects {
		log.Printf("Processing project: %v with ID: %v", nProject.Title, nProject.Id)

		//if nProject.Status == "Cancelled" {
		//	log.Printf("Project %v has been cancelled. Not processing", nProject.Title)
		//	continue
		//}

		nConcerts, err := nc.GetConcerts(nProject.Id)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to get notion locations: %v", err))
		}
		if nConcerts == nil {
			log.Println("No notion project concerts found. Skipping project...")
			continue
		}

		gDriveImageUrl := ""
		if nProject.Poster != "" {
			gDriveImageUrl, err = helpers.BuildGDriveImageUrl(nProject.Poster)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("failed to convert GDrive image share url to download url: %v", err))
			}
		}
		log.Printf("Download image URL generated: %v", gDriveImageUrl)

		var concerts []*notion.Concert
		for _, nConcert := range nConcerts {
			log.Printf("Adding concert with ID: %v starting date: %v", nConcert.Id, nConcert.StartDateAndTime)
			setConcertsForProject(&nConcert, nLocations, &concerts)
		}
		project := &notion.Project{
			Project:        &nProject,
			Dates:          concerts,
			GDriveImageUrl: gDriveImageUrl,
		}

		projects = append(projects, project)
		log.Printf("Project %v added", project.Project.Title)
	}

	return projects, nil
}

func pipeline() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Initializing notion client...")
	nc, err := notion.NewNotionClient(cfg.NotionAPIKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Notion client initialized")

	log.Println("Initializing Git client...")
	gitClient, err := helpers.NewMyGit(cfg.GitPAT)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Cloning repository...")
	err = gitClient.CloneRepository(clonePath, repoURL)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Cloned repository.")

	log.Println("Checking out branch...")
	branchRefName, err := gitClient.Reference(branchName)
	if err != nil {
		err = gitClient.CheckoutBranch(branchName, true)
		if err != nil {
			log.Fatalf("could not create a new branch: %v", err)
		}
	}
	err = gitClient.CheckoutBranch(branchName, false)
	if err != nil {
		log.Fatalf("could not checkout branch %v: %v", branchRefName.String(), err)
	}

	log.Println("Pulling Notion Projects...")
	nProjects, err := nc.GetProjects()
	if err != nil {
		log.Fatalf("failed to get notion projects: %v", err)
	}
	log.Println("Pulling Notion Locations...")
	nLocations, err := nc.GetLocations()
	if err != nil {
		log.Fatalf("failed to get notion locations: %v", err)
	}

	log.Println("Setting up Projects list...")
	projects, err := setProjects(nProjects, nc, nLocations)
	if err != nil {
		log.Fatalf("failed to set Project list: %v", err)
	}
	log.Printf("%v projects where found.", len(projects))

	log.Println("Emptying post folder...")
	err = helpers.NukeFolder(fmt.Sprintf("%v/%v", clonePath, postsPath))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Generating posts from projects...")
	for _, project := range projects {
		log.Printf("Processing project: %v", project.Project.Title)
		postFolderName := fmt.Sprintf("%v/%v/%v-%v", clonePath, postsPath, project.Project.Year, project.Project.Id)
		log.Printf("Creating post folder: %v", postFolderName)
		err = helpers.CreateFolder(postFolderName)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Setting up template...")
		postTemplateDetails, err := setTemplateDataDetails(project.Dates)
		if err != nil {
			log.Fatal(err)
		}
		postTemplateData := setPostTemplateData(project, postTemplateDetails)

		log.Printf("Generating post as %v/%v", postFolderName, postFileName)
		err = template.GeneratePost(templateFileName, postFolderName, postFileName, postTemplateData)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Adding images to post...")
		if project.GDriveImageUrl != "" {
			imagePath := fmt.Sprintf("%v/%v", postFolderName, imageName)
			err = helpers.DownloadImage(project.GDriveImageUrl, imagePath)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	log.Println("Add, commit and push...")
	commitHash, err := gitClient.AddAndCommit()
	if err != nil {
		log.Fatal(err)
	}

	if commitHash != nil {
		log.Println("Commit:", commitHash.String())
		err = gitClient.Push(branchName)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Pushed: %v", commitHash.String())
	} else {
		log.Println("No changes detected skipping commit.")
	}

	log.Println("Cleaning up...")
	err = os.RemoveAll(clonePath)
	if err != nil {
		log.Fatalf("failed to clean up: %v", err)
	}
	log.Println("Repository deleted.")

	return
}

package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var apiVersion string = "v3"
var markedForDeletionTag string = "cma-markedfordeletion"

func GetMoviesData() ([]byte, error) {
	apiUrl := Config.Radarr.URL + "/api/" + apiVersion + "/movie"
	apiKey, err := Base64Decode(Config.Radarr.B64APIKey)
	// fmt.Println(apiKey)
	if err != nil {
		return nil, err
	}

	// Create request
	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		log.Println("Failed to get movies library", err.Error())
		return nil, err
	}
	req.Header.Set("Authorization", apiKey)

	// Create client
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request
	res, err := client.Do(req)
	if err != nil {
		log.Println("Failed to get movies library", err.Error())
		return nil, err
	}
	if res.StatusCode != 200 {
		log.Println("Failed to get movies library", res.Status)
		return nil, errors.New("Failed to get movies library")
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Failed to get movies library", err.Error())
		return nil, err
	}

	return data, nil
}

func GetMovieImportEvents(movieId int) ([]MovieImportEvent, error) {
	apiUrl := Config.Radarr.URL + "/api/" + apiVersion + "/history/movie?eventType=downloadFolderImported&movieId=" + fmt.Sprintf("%d", movieId)
	apiKey, err := Base64Decode(Config.Radarr.B64APIKey)
	// fmt.Println(apiKey)
	if err != nil {
		return nil, err
	}

	// Create request
	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		log.Println("Failed to get movie's import events", err.Error())
		return nil, err
	}
	req.Header.Set("Authorization", apiKey)

	// Create client
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request
	res, err := client.Do(req)
	if err != nil {
		log.Println("Failed to get movie's import events", err.Error())
		return nil, err
	}
	if res.StatusCode != 200 {
		log.Println("Failed to get movie's import events", res.Status)
		return nil, errors.New("Failed to get movie's import events")
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Failed to get movie's import events", err.Error())
		return nil, err
	}

	var movieImportEvent []MovieImportEvent

	err = json.Unmarshal(data, &movieImportEvent)
	if err != nil {
		log.Println("Failed to get movie's import events", err.Error())
		return nil, err
	}

	return movieImportEvent, nil
}

func MarkMoviesForDeletion(moviesdata []byte, moviesIgnored []string, nextMaintenanceCycle time.Time, isDryRun bool) ([]string, error) {
	// fmt.Println("Entered function MarkMoviesForDeletion")
	apiUrl := Config.Radarr.URL + "/api/" + apiVersion + "/movie/editor"
	apiKey, err := Base64Decode(Config.Radarr.B64APIKey)
	if err != nil {
		return nil, err
	}

	durationToNextMaintenance := nextMaintenanceCycle.Sub(Now).Hours() / 24

	var movies []Movie

	err = json.Unmarshal(moviesdata, &movies)
	if err != nil {
		log.Println("Failed to mark movies for deletion", err.Error())
		return nil, err
	}

	tagId, err := GetTagIdFromRadarr(markedForDeletionTag)
	if err != nil {
		return nil, err
	}
	if tagId == nil {
		tagId, err = CreateTagInRadarr(markedForDeletionTag)
		if err != nil {
			return nil, err
		}
	}

	var movieIdsMarkedForDeletionStringified string = "["
	var tagIdStringified string = "[" + fmt.Sprintf("%d", *tagId) + "]"
	var emptyList bool = true
	var movieNamesMarkedForDeletion []string

	for _, movie := range movies {
		if !movie.HasFile {
			continue
		}

		var checkIfIgnored bool = false
		for _, ignoredMovie := range moviesIgnored {
			if movie.Title == ignoredMovie {
				checkIfIgnored = true
				break
			}
		}

		if checkIfIgnored {
			continue
		}

		durationInDays, err := GetMovieAge(movie)
		if err != nil {
			return nil, err
		}

		ageOfMovieAtMaintenance := *durationInDays + durationToNextMaintenance

		if ageOfMovieAtMaintenance >= float64(Config.DeleteAfterDays) {
			emptyList = false
			movieIdsMarkedForDeletionStringified = movieIdsMarkedForDeletionStringified + fmt.Sprintf("%d,", movie.ID)
			movieNamesMarkedForDeletion = append(movieNamesMarkedForDeletion, movie.Title)
		}
	}

	movieIdsMarkedForDeletionStringified = strings.TrimSuffix(movieIdsMarkedForDeletionStringified, ",")
	movieIdsMarkedForDeletionStringified = movieIdsMarkedForDeletionStringified + "]"

	if emptyList {
		log.Println("No movies to mark for deletion")
		return nil, nil
	}

	// Create request
	reqBodyValue := `{"movieIds": ` + movieIdsMarkedForDeletionStringified + `, "tags": ` + tagIdStringified + `, "applyTags": "add"}`
	requestBody := bytes.NewReader([]byte(reqBodyValue))

	// fmt.Println(reqBodyValue)

	req, err := http.NewRequest(http.MethodPut, apiUrl, requestBody)
	if err != nil {
		log.Println("Failed to mark movies for deletion", err.Error())
		return nil, err
	}
	req.Header.Set("Authorization", apiKey)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	// Create client
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	if !isDryRun {
		// Make request
		res, err := client.Do(req)
		if err != nil {
			log.Println("Failed to mark movies for deletion", err.Error())
			return nil, err
		}
		if res.StatusCode/100 != 2 {
			log.Println("Failed to mark movies for deletion", res.Status)
			return nil, errors.New("Failed mark movies for deletion")
		}
	}

	log.Println("Movies marked for deletion:", movieNamesMarkedForDeletion)
	return movieNamesMarkedForDeletion, nil

}

func GetMovieAge(movie Movie) (*float64, error) {
	// dateAdded := movie.MovieFile.DateAdded
	movieImportEvents, err := GetMovieImportEvents(movie.MovieFile.MovieId)
	if err != nil {
		log.Println("Failed get age of movie", movie.Title, err.Error())
		return nil, err
	}

	sizeOfEvents := len(movieImportEvents)
	dateAdded := movieImportEvents[sizeOfEvents-1].Date

	// fmt.Println(dateAdded)
	parsedDateAdded, err := time.Parse("2006-01-02T15:04:05Z", dateAdded)
	if err != nil {
		log.Println("Failed get age of movie", movie.Title, err.Error())
		return nil, err
	}
	durationInDays := Now.Sub(parsedDateAdded).Hours() / 24
	// fmt.Println(dateAdded, parsedDateAdded, durationInDays, movie.Tags)

	// fmt.Println(movie.Title, dateAdded)
	// fmt.Println(durationInDays)

	movieAgeDetail, _ := json.Marshal(map[string]any{"Title": movie.Title, "AddedOn": dateAdded, "Age": fmt.Sprintf("%.1f", durationInDays) + " days"})
	log.Println(string(movieAgeDetail))

	return &durationInDays, nil
}

func GetTagIdFromRadarr(tagLabel string) (*int, error) {
	apiUrl := Config.Radarr.URL + "/api/" + apiVersion + "/tag"
	apiKey, err := Base64Decode(Config.Radarr.B64APIKey)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		log.Println("Failed to get", tagLabel, "tags", err.Error())
		return nil, err
	}
	req.Header.Set("Authorization", apiKey)

	// Create client
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request
	res, err := client.Do(req)
	if err != nil {
		log.Println("Failed to get", tagLabel, "tags", err.Error())
		return nil, err
	}
	if res.StatusCode != 200 {
		log.Println("Failed to get", tagLabel, "tags", res.Status)
		return nil, errors.New("Failed to get " + tagLabel + " tags")
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Failed to get", tagLabel, "tags", err.Error())
		return nil, err
	}

	var tags []Tag

	err = json.Unmarshal(data, &tags)
	if err != nil {
		log.Println("Failed to get", tagLabel, "tags", err.Error())
		return nil, err
	}

	for _, tag := range tags {
		if tag.Label == tagLabel {
			return &tag.Id, nil
		}
	}

	return nil, nil
}

func CreateTagInRadarr(tagLabel string) (*int, error) {
	apiUrl := Config.Radarr.URL + "/api/" + apiVersion + "/tag"
	apiKey, err := Base64Decode(Config.Radarr.B64APIKey)
	// fmt.Println(apiKey)
	if err != nil {
		return nil, err
	}

	// Create request
	reqBodyValue := `{"label": "` + tagLabel + `"}`
	// reqBodyValue := []byte(`{"label": cma-markedfordeletion}`)
	requestBody := bytes.NewReader([]byte(reqBodyValue))

	if err != nil {
		log.Println("Failed to create "+tagLabel+" tag", err.Error())
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, apiUrl, requestBody)
	if err != nil {
		log.Println("Failed to create "+tagLabel+" tag", err.Error())
		return nil, err
	}
	req.Header.Set("Authorization", apiKey)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	// Create client
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request
	res, err := client.Do(req)
	if err != nil {
		log.Println("Failed to create "+tagLabel+" tag", err.Error())
		return nil, err
	}
	if res.StatusCode/100 != 2 {
		log.Println("Failed to create "+tagLabel+" tag", res.Status)
		return nil, errors.New("Failed to create " + tagLabel + " tag")
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Failed to create "+tagLabel+" tag", err.Error())
		return nil, err
	}

	// fmt.Println(string(data))

	var tag Tag

	err = json.Unmarshal(data, &tag)
	if err != nil {
		log.Println("Failed to create "+tagLabel+" tag", err.Error())
		return nil, err
	}
	// fmt.Println(tag.Id, tag.Label)

	return &tag.Id, nil
}

func DeleteExpiredMovies(moviesdata []byte, moviesIgnored []string, nextMaintenanceCycle time.Time, isDryRun bool) ([]string, error) {
	apiUrl := Config.Radarr.URL + "/api/" + apiVersion + "/movie"
	apiKey, err := Base64Decode(Config.Radarr.B64APIKey)
	if err != nil {
		return nil, err
	}

	var movies []Movie

	err = json.Unmarshal(moviesdata, &movies)
	if err != nil {
		log.Println("Failed to delete expired movies", err.Error())
		return nil, err
	}

	var emptyList bool = true
	var moviesDeleted []string
	var moviesFailedToDelete []string

	for _, movie := range movies {
		if !movie.HasFile {
			continue
		}

		var checkIfIgnored bool = false
		for _, ignoredMovie := range moviesIgnored {
			if movie.Title == ignoredMovie {
				checkIfIgnored = true
				break
			}
		}

		if checkIfIgnored {
			continue
		}

		durationInDays, err := GetMovieAge(movie)
		if err != nil {
			log.Println("Failed to delete movie", movie.Title)
			moviesFailedToDelete = append(moviesFailedToDelete, movie.Title)
			continue
		}

		durationToNextMaintenance := nextMaintenanceCycle.Sub(Now).Hours() / 24

		ageOfMovieAtMaintenance := *durationInDays + durationToNextMaintenance

		if ageOfMovieAtMaintenance >= float64(Config.DeleteAfterDays) {
			emptyList = false
			deleteApiURL := apiUrl + fmt.Sprintf("/%d", movie.ID) + "?deleteFiles=true&addImportExclusion=true"

			// Create movie delete request
			req, err := http.NewRequest(http.MethodDelete, deleteApiURL, nil)
			if err != nil {
				log.Println("Failed to delete movie", movie.Title, err.Error())
				moviesFailedToDelete = append(moviesFailedToDelete, movie.Title)
				continue
			}
			req.Header.Set("Authorization", apiKey)

			// Create client
			client := http.Client{
				Timeout: 30 * time.Second,
			}

			// Make requests
			if !isDryRun {

				res, err := client.Do(req)

				if err != nil {
					log.Println("Failed to delete movie", movie.Title, err.Error())
					moviesFailedToDelete = append(moviesFailedToDelete, movie.Title)
					continue
				}
				if res.StatusCode/100 != 2 {
					log.Println("Failed to delete movie", movie.Title, err.Error())
					moviesFailedToDelete = append(moviesFailedToDelete, movie.Title)
					continue
				}
			}

			moviesDeleted = append(moviesDeleted, movie.Title)
		}
	}

	log.Println("Movies ignored:", moviesIgnored)

	if emptyList {
		log.Println("No movies to delete")
		return nil, nil
	}

	log.Println("Movies deleted:", moviesDeleted)
	log.Println("Movies failed to delete:", moviesFailedToDelete)

	return moviesDeleted, nil
}

func GetMoviesIgnored(ignoreTagId int, moviesdata []byte) ([]string, error) {
	var movies []Movie
	var moviesIgnored []string

	err := json.Unmarshal(moviesdata, &movies)
	if err != nil {
		log.Println("Failed to get ignored movies", err.Error())
		return nil, err
	}

	for _, movie := range movies {
		if !movie.HasFile {
			continue
		}

		var checkIgnoreTag bool = false
		for _, tag := range movie.Tags {
			if tag == ignoreTagId {
				checkIgnoreTag = true
				break
			}
		}

		if checkIgnoreTag {
			moviesIgnored = append(moviesIgnored, movie.Title)
		}
	}

	return moviesIgnored, nil
}

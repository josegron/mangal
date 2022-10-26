package anilist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/metafates/mangal/log"
	"github.com/metafates/mangal/network"
	"net/http"
	"strconv"
)

type searchByNameResponse struct {
	Data struct {
		Page struct {
			Media []*Manga `json:"media"`
		} `json:"page"`
	} `json:"data"`
}

type searchByIDResponse struct {
	Data struct {
		Media *Manga `json:"media"`
	} `json:"data"`
}

func GetByID(id int) (*Manga, error) {
	if manga, ok := idCacher.Get(id); ok {
		return manga, nil
	}

	// prepare body
	log.Infof("Searching anilist for manga with id: %d", id)
	body := map[string]interface{}{
		"query": searchByIDQuery,
		"variables": map[string]interface{}{
			"id": id,
		},
	}

	// parse body to json
	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// send request
	log.Info("Sending request to Anilist")
	req, err := http.NewRequest(http.MethodPost, "https://graphql.anilist.co", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := network.Client.Do(req)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("Anilist returned status code " + strconv.Itoa(resp.StatusCode))
		return nil, fmt.Errorf("invalid response code %d", resp.StatusCode)
	}

	// decode response
	var response searchByIDResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Error(err)
		return nil, err
	}

	manga := response.Data.Media
	log.Infof("Got response from Anilist, found manga with id %d", manga.ID)
	_ = idCacher.Set(id, manga)
	return manga, nil
}

func SearchByName(name string) ([]*Manga, error) {
	if mangas, ok := searchCacher.Get(name); ok {
		return mangas, nil
	}

	// prepare body
	log.Info("Searching anilist for manga: " + name)
	body := map[string]interface{}{
		"query": searchByNameQuery,
		"variables": map[string]interface{}{
			"query": name,
		},
	}

	// parse body to json
	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// send request
	log.Info("Sending request to Anilist")
	req, err := http.NewRequest(http.MethodPost, "https://graphql.anilist.co", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := network.Client.Do(req)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("Anilist returned status code " + strconv.Itoa(resp.StatusCode))
		return nil, fmt.Errorf("invalid response code %d", resp.StatusCode)
	}

	// decode response
	var response searchByNameResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Error(err)
		return nil, err
	}

	mangas := response.Data.Page.Media
	log.Infof("Got response from Anilist, found %d results", strconv.Itoa(len(mangas)))
	_ = searchCacher.Set(name, mangas)
	return mangas, nil
}

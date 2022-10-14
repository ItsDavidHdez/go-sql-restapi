package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ItsDavidHdez/go-sql-restapi/db"
	"github.com/ItsDavidHdez/go-sql-restapi/models"
)

func GetAllMusicHander(w http.ResponseWriter, r *http.Request) {
	var songs []models.Music
	db.DB.Find(&songs)
	json.NewEncoder(w).Encode(&songs)
}

func GetMusicHander(w http.ResponseWriter, r *http.Request) {
	var responseObject models.Response
	var music models.Music
	term := r.URL.Query().Get("term")
	newTerm := strings.Replace(term, " ", "+", 10)

	response, err := http.Get("https://itunes.apple.com/search?term=" + newTerm + "&limit=1")

    if err != nil {
        fmt.Print(err.Error())
        os.Exit(1)
    }

    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }
	
	json.Unmarshal(responseData, &responseObject)

	for i := 0; i < len(responseObject.Music); i++{
		dataJson := models.Music{
			Name: responseObject.Music[i].Name,
			Artist: responseObject.Music[i].Artist,
			Duration: responseObject.Music[i].Duration,
			Album: responseObject.Music[i].Album,
			Artwork: responseObject.Music[i].Artwork,
			Price: responseObject.Music[i].Price,
			Origin: responseObject.Music[i].Origin,
		}
		createJson, _ := json.Marshal(dataJson)

		json.NewDecoder(strings.NewReader(string(createJson))).Decode(&music)
		createdPlaylist := db.DB.Create(&music)
		err := createdPlaylist.Error
		
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		} 

		json.NewEncoder(w).Encode(&music)
	}
}

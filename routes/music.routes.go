package routes

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/ItsDavidHdez/go-sql-restapi/db"
	"github.com/ItsDavidHdez/go-sql-restapi/models"
	jwt "github.com/dgrijalva/jwt-go"
)

func validToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenStr := cookie.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET")), nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Printf("You are logged")
}


func GetAllMusicHander(w http.ResponseWriter, r *http.Request) {

	validToken(w, r)

	var songs []models.Music
	var songsSoap []models.MusicSoap

	db.DB.Find(&songs)
	sort.Slice(songs, func(i, j int) bool {
		return songs[i].Name < songs[j].Name
	})

	db.DB.Find(&songsSoap)
	sort.Slice(songsSoap, func(i, j int) bool {
		return songsSoap[i].Name < songsSoap[j].Name
	})

	json.NewEncoder(w).Encode(&songs)
	json.NewEncoder(w).Encode(&songsSoap)
}

func GetMusicHander(w http.ResponseWriter, r *http.Request) {
	// validToken(w, r)
	
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
		// createdPlaylist := db.DB.Create(&music)
		// err := createdPlaylist.Error
		
		// if err != nil {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	w.Write([]byte(err.Error()))
		// } 
	}
			json.NewEncoder(w).Encode(&music)
}

func GetMusicSoapHandler(w http.ResponseWriter, r *http.Request) {

	var music models.MusicSoap
	var MusicSoapResponse models.MusicSoapResponse

	artistParam := r.URL.Query().Get("artist")
	songParam := r.URL.Query().Get("song")

	newArtistParam := strings.Replace(artistParam, " ", "+", 10)
	newSongParam := strings.Replace(songParam, " ", "+", 10)


	responseSoapData, errSoap := http.Get("http://api.chartlyrics.com/apiv1.asmx/SearchLyricDirect?artist=" + newArtistParam +"&song=" + newSongParam)

	if errSoap != nil {
        fmt.Print(errSoap.Error())
        os.Exit(1)
    }

	responseDataSoap, errResponseDataSoap := ioutil.ReadAll(responseSoapData.Body)

	if errResponseDataSoap != nil {
		log.Fatal(errResponseDataSoap)
	}

	xml.Unmarshal([]byte(responseDataSoap), &music)

	fmt.Println("responseObject", music)

	for i := 0; i < len(MusicSoapResponse.MusicSoapResponse); i++{
		dataJson := models.MusicSoap{
			Name: MusicSoapResponse.MusicSoapResponse[i].Name,
			Artist: MusicSoapResponse.MusicSoapResponse[i].Artist,
			Album: MusicSoapResponse.MusicSoapResponse[i].Album,
			Artwork: MusicSoapResponse.MusicSoapResponse[i].Artwork,
			Origin: MusicSoapResponse.MusicSoapResponse[i].Origin,
		}

		fmt.Println("hola")

		createJson, _ := json.Marshal(dataJson)

		json.NewDecoder(strings.NewReader(string(createJson))).Decode(&music)
		createdPlaylist := db.DB.Create(&music)
		err := createdPlaylist.Error
		
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		} 
	}

	json.NewEncoder(w).Encode(&music)

}

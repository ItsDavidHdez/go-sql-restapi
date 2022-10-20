package models

import "gorm.io/gorm"

type Response struct {
	Music []Music `json:"results"`
}

type MusicSoapResponse struct {
	MusicSoapResponse []MusicSoap `json:"results"`
}

type Music struct {
	gorm.Model

	Name string `json:"trackName" xml:"LyricSong"`
	Artist string `json:"artistName" xml:"LyricArtist"`
	Duration string `json:"trackTimeMillis"`
	Album string `json:"collectionName" xml:"LyricCovertArtUrl"`
	Artwork string `json:"artworkUrl100" xml:"LyricCorrectUrl"`
	Price float64 `json:"trackPrice"`
	Origin string `json:"trackViewUrl" xml:"LyricUrl"`
}

type MusicSoap struct {
	gorm.Model 

	Name string `gorm:"not null" xml:"LyricSong"`
	Artist string `gorm:"not null" xml:"LyricArtist"`
	Album string `gorm:"not null" xml:"LyricCovertArtUrl"`
	Artwork string `gorm:"not null" xml:"LyricCorrectUrl"`
	Origin string `gorm:"not null" xml:"LyricUrl"`
	Music []Music `json:"results"`
}
package models

import "gorm.io/gorm"

type Response struct {
	Music []Music `json:"results"`
}

type Music struct {
	gorm.Model

	Name string `gorm:"not null" json:"trackName"`
	Artist string `gorm:"not null" json:"artistName"`
	Duration string `gorm:"not null" json:"trackTimeMillis"`
	Album string `gorm:"not null" json:"collectionName"`
	Artwork string `gorm:"not null" json:"artworkUrl100"`
	Price float64 `gorm:"not null" json:"trackPrice"`
	Origin string `gorm:"not null" json:"trackViewUrl"`
}
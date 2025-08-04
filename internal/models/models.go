package models

import "encoding/xml"

type Users struct {
	XMLName xml.Name `xml:"users"`
	Users   []User   `xml:"user"`
}

type User struct {
	ID    string `xml:"id"`
	Name  string `xml:"name"`
	Email string `xml:"email"`
	Age   int    `xml:"age"`
}

type UserJSON struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	AgeGroup string `json:"age_group"`
}

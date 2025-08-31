package models

type Actor struct {
	ID        int     `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	ImageURL  *string `json:"image_url,omitempty"` //If URL not avilable it will point to nill value!
}

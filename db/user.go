package database

type User struct {
	Id       string      `json:"id"`
	Email    string      `json:"string"`
	Password string      `json:"-"`
	Profile  UserProfile `json:"profile"`
}

type UserProfile struct {
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	ProfilePictureURL string    `json:"pfp_url"`
	Bio               string    `json:"bio"`
	SocialLinks       []Socials `json:"socials"`
}

type Socials struct {
	Url      string `json:"url"`
	LinkName string `json:"link_name"`
}



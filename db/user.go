package db

import "fmt"

type User struct {
	Id       string       `pg:"id,pk"`
	Email    string       `pg:"email,unique"`
	Password string       `pg:"password"`
	Profile  *UserProfile `pg:"rel:has-one"`
}

type UserProfile struct {
	Id                int64     `pg:"id,pk"`
	UserId            string    `pg:"user_id"`
	Name              string    `pg:"name"`
	Email             string    `pg:"email"`
	ProfilePictureURL string    `pg:"profile_picture_url"`
	Bio               string    `pg:"bio"`
	SocialLinks       []Socials `pg:"rel:has-many"`
}

type Socials struct {
	Id            int64  `pg:"id,pk"`
	UserProfileId int64  `pg:"user_profile_id"`
	Url           string `pg:"url"`
	LinkName      string `pg:"link_name"`
}

func (u User) String() string {
	return fmt.Sprintf("User<%s, %s>", u.Id, u.Email)
}

func (p UserProfile) String() string {
	return fmt.Sprintf("Profile<%s, %s, %s, %s>", p.Name, p.Email, p.ProfilePictureURL, p.Bio)
}

func (s Socials) String() string {
	return fmt.Sprintf("Socials<%s %s>", s.Url, s.LinkName)
}

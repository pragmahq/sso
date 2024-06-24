package database

import (
	"fmt"

	"github.com/go-pg/pg/v10"
)

const (
	PermissionUser = 1 << iota
	PermissionEditor
	PermissionAdmin
)

type User struct {
	Id               string        `pg:"id,pk"`
	Email            string        `pg:"email,unique"`
	Password         string        `pg:"password"`
	Permissions      int           `pg:"permissions"`
	Profile          *UserProfile  `pg:"rel:has-one"`
	GeneratedInvites []*InviteCode `pg:"rel:has-many,fk:generated_by"`
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

func (u *User) Create(db *DB) error {
	_, err := db.Model(u).Insert()
	return err
}

func (u *User) Read(db *DB) error {
	return db.Model(u).WherePK().Select()
}

func (u *User) Update(db *DB) error {
	_, err := db.Model(u).WherePK().Update()
	return err
}

func (u *User) Delete(db *DB) error {
	_, err := db.Model(u).WherePK().Delete()
	return err
}

func (u *User) IsUser() bool {
	return u.Permissions&PermissionUser != 0
}

func (u *User) IsEditor() bool {
	return u.Permissions&PermissionEditor != 0
}

func (u *User) IsAdmin() bool {
	return u.Permissions&PermissionAdmin != 0
}

func (u *User) SetUser() {
	u.Permissions |= PermissionUser
}

func (u *User) SetEditor() {
	u.Permissions |= PermissionEditor
}

func (u *User) SetAdmin() {
	u.Permissions |= PermissionAdmin
}

func (u *User) RemoveUser() {
	u.Permissions &^= PermissionUser
}

func (u *User) RemoveEditor() {
	u.Permissions &^= PermissionEditor
}

func (u *User) RemoveAdmin() {
	u.Permissions &^= PermissionAdmin
}

func (p *UserProfile) Create(db *DB) error {
	_, err := db.Model(p).Insert()
	return err
}

func (p *UserProfile) Read(db *DB) error {
	return db.Model(p).WherePK().Select()
}

func (p *UserProfile) Update(db *DB) error {
	_, err := db.Model(p).WherePK().Update()
	return err
}

func (p *UserProfile) Delete(db *DB) error {
	_, err := db.Model(p).WherePK().Delete()
	return err
}

func (s *Socials) Create(db *DB) error {
	_, err := db.Model(s).Insert()
	return err
}

func (s *Socials) Read(db *DB) error {
	return db.Model(s).WherePK().Select()
}

func (s *Socials) Update(db *DB) error {
	_, err := db.Model(s).WherePK().Update()
	return err
}

func (s *Socials) Delete(db *DB) error {
	_, err := db.Model(s).WherePK().Delete()
	return err
}

func (u *User) UpdatePermissions(db *pg.DB) error {
	_, err := db.Model(u).
		Set("permissions = ?", u.Permissions).
		Where("id = ?", u.Id).
		Update()
	return err
}

func (u *User) GetHighestPermission() string {
	if u.IsAdmin() {
		return "Admin"
	} else if u.IsEditor() {
		return "Editor"
	} else if u.IsUser() {
		return "User"
	}
	return "None"
}

func (u *User) GetUserWithProfile(db *DB) error {
	return db.Model(u).Relation("Profile").WherePK().Select()
}

func (p *UserProfile) GetUserProfileWithSocials(db *DB) error {
	return db.Model(p).Relation("SocialLinks").WherePK().Select()
}

func GetUserByEmail(db *DB, email string) (*User, error) {
	user := &User{}
	err := db.Model(user).Where("email = ?", email).Select()
	if err != nil {
		return nil, err
	}
	return user, nil
}

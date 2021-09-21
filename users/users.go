package users

import (
	"encoding/json"
	"io/ioutil"

	"gopkg.in/Iwark/spreadsheet.v2"
)

const (
	userDataJsonFile = "user_data.json"
)

type User struct {
	Email            string
	FullName         string
	Name             string
	Phone            string
	SocialLink       string
	City             string
	ExperienceAnswer string
	Age              string
}

type Users = map[string]User

type SheetConfig struct {
	MailIdx             int
	NameIdx             int
	PhoneIdx            int
	SocialLinkIdx       int
	CityIdx             int
	ExperienceAnswerIdx int
	AgeIdx              int
	Skip                int
}

type users struct {
	users Users
}

type UsersInt interface {
	AddUsers(sheet *spreadsheet.Sheet, config *SheetConfig) (err error)
	GetUsers() Users
	DumpUsers() error
}

func (u *users) AddUsers(sheet *spreadsheet.Sheet, config *SheetConfig) (err error) {
	for i := range sheet.Rows {
		if i < config.Skip {
			// skip
			continue
		}
		user := User{}

		var mail string
		if len(sheet.Rows[i]) > config.MailIdx {
			mail = sheet.Rows[i][config.MailIdx].Value
		}
		user.Email = mail

		var fullName string
		if len(sheet.Rows[i]) > config.NameIdx {
			fullName = sheet.Rows[i][config.NameIdx].Value
		}
		user.FullName = fullName
		user.Name = fullName

		var age string
		if len(sheet.Rows[i]) > config.AgeIdx {
			age = sheet.Rows[i][config.AgeIdx].Value
		}
		user.Age = age

		var phone string
		if len(sheet.Rows[i]) > config.PhoneIdx {
			phone = sheet.Rows[i][config.PhoneIdx].Value
		}
		user.Phone = phone

		var socialLink string
		if len(sheet.Rows[i]) > config.SocialLinkIdx {
			socialLink = sheet.Rows[i][config.SocialLinkIdx].Value
		}
		user.SocialLink = socialLink

		var city string
		if len(sheet.Rows[i]) > config.CityIdx {
			city = sheet.Rows[i][config.CityIdx].Value
		}
		user.City = city

		var experienceAnswer string
		if len(sheet.Rows[i]) > config.ExperienceAnswerIdx {
			experienceAnswer = sheet.Rows[i][config.ExperienceAnswerIdx].Value
		}
		user.ExperienceAnswer = experienceAnswer

		u.users[mail] = user
	}
	return err
}

func (u users) GetUsers() Users {
	return u.users
}

func (u users) DumpUsers() error {
	file, err := json.MarshalIndent(u.users, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(userDataJsonFile, file, 0644) //nolint:gosec
	return err
}
func NewUsers() UsersInt {
	makeUsers := make(Users)
	return &users{
		users: makeUsers,
	}
}

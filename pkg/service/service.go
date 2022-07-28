package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

type Data struct {
	Emails []string `json:"emails"`
}

func (s *Service) GetEmails() ([]string, error) {
	if _, err := os.Stat("emails.json"); errors.Is(err, os.ErrNotExist) {
		_, err := os.Create("emails.json")
		if err != nil {
			return nil, err
		}
	}
	byteValue, err := ioutil.ReadFile("emails.json")
	if err != nil {
		return nil, err
	}
	var result Data
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		err = s.WriteEmails(nil)
		return nil, err
	}

	return result.Emails, nil
}

func (s *Service) WriteEmails(emails []string) error {
	writeData := Data{Emails: emails}

	file, err := json.MarshalIndent(writeData, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("emails.json", file, 0644)
	return err
}

func (s *Service) EmailIsExist(email string, emails []string) bool {
	for _, v := range emails {
		if v == email {
			return true
		}
	}
	return false
}

func (s *Service) AddEmail(email string) error {
	emails, err := s.GetEmails()
	if err != nil {
		return err
	}
	if s.EmailIsExist(email, emails) {
		return errors.New("email already subscribed")
	}
	emails = append(emails, email)

	return s.WriteEmails(emails)
}

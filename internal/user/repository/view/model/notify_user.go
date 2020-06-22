package model

import (
	"encoding/json"
	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	"time"
)

const (
	NotifyUserKeyUserID = "id"
)

type NotifyUser struct {
	ID                string    `json:"-" gorm:"column:id;primary_key"`
	CreationDate      time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate        time.Time `json:"-" gorm:"column:change_date"`
	ResourceOwner     string    `json:"-" gorm:"column:resource_owner"`
	UserName          string    `json:"userName" gorm:"column:user_name"`
	FirstName         string    `json:"firstName" gorm:"column:first_name"`
	LastName          string    `json:"lastName" gorm:"column:last_name"`
	NickName          string    `json:"nickName" gorm:"column:nick_name"`
	DisplayName       string    `json:"displayName" gorm:"column:display_name"`
	PreferredLanguage string    `json:"preferredLanguage" gorm:"column:preferred_language"`
	Gender            int32     `json:"gender" gorm:"column:gender"`
	LastEmail         string    `json:"email" gorm:"column:last_email"`
	VerifiedEmail     string    `json:"-" gorm:"column:verified_email"`
	LastPhone         string    `json:"phone" gorm:"column:last_phone"`
	VerifiedPhone     string    `json:"-" gorm:"column:verified_phone"`
	PasswordSet       bool      `json:"-" gorm:"column:password_set"`
	Sequence          uint64    `json:"-" gorm:"column:sequence"`
}

func NotifyUserFromModel(user *model.NotifyUser) *NotifyUser {
	return &NotifyUser{
		ID:                user.ID,
		ChangeDate:        user.ChangeDate,
		CreationDate:      user.CreationDate,
		ResourceOwner:     user.ResourceOwner,
		UserName:          user.UserName,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		NickName:          user.NickName,
		DisplayName:       user.DisplayName,
		PreferredLanguage: user.PreferredLanguage,
		Gender:            int32(user.Gender),
		LastEmail:         user.LastEmail,
		VerifiedEmail:     user.VerifiedEmail,
		LastPhone:         user.LastPhone,
		VerifiedPhone:     user.VerifiedPhone,
		PasswordSet:       user.PasswordSet,
		Sequence:          user.Sequence,
	}
}

func NotifyUserToModel(user *NotifyUser) *model.NotifyUser {
	return &model.NotifyUser{
		ID:                user.ID,
		ChangeDate:        user.ChangeDate,
		CreationDate:      user.CreationDate,
		ResourceOwner:     user.ResourceOwner,
		UserName:          user.UserName,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		NickName:          user.NickName,
		DisplayName:       user.DisplayName,
		PreferredLanguage: user.PreferredLanguage,
		Gender:            model.Gender(user.Gender),
		LastEmail:         user.LastEmail,
		VerifiedEmail:     user.VerifiedEmail,
		LastPhone:         user.LastPhone,
		VerifiedPhone:     user.VerifiedPhone,
		PasswordSet:       user.PasswordSet,
		Sequence:          user.Sequence,
	}
}

func (u *NotifyUser) AppendEvent(event *models.Event) (err error) {
	u.ChangeDate = event.CreationDate
	u.Sequence = event.Sequence
	switch event.Type {
	case es_model.UserAdded,
		es_model.UserRegistered:
		u.CreationDate = event.CreationDate
		u.setRootData(event)
		err = u.setData(event)
		if err != nil {
			return err
		}
		err = u.setPasswordData(event)
	case es_model.UserProfileChanged:
		err = u.setData(event)
	case es_model.UserEmailChanged:
		err = u.setData(event)
	case es_model.UserEmailVerified:
		u.VerifiedEmail = u.LastEmail
	case es_model.UserPhoneChanged:
		err = u.setData(event)
	case es_model.UserPhoneVerified:
		u.VerifiedPhone = u.LastPhone
	case es_model.UserPasswordChanged:
		err = u.setPasswordData(event)
	}
	return err
}

func (u *NotifyUser) setRootData(event *models.Event) {
	u.ID = event.AggregateID
	u.ResourceOwner = event.ResourceOwner
}

func (u *NotifyUser) setData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, u); err != nil {
		logging.Log("MODEL-lso9e").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(nil, "MODEL-8iows", "could not unmarshal data")
	}
	return nil
}

func (u *NotifyUser) setPasswordData(event *models.Event) error {
	password := new(es_model.Password)
	if err := json.Unmarshal(event.Data, password); err != nil {
		logging.Log("MODEL-dfhw6").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(nil, "MODEL-BHFD2", "could not unmarshal data")
	}
	u.PasswordSet = password.Secret != nil
	return nil
}
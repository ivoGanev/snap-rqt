package model

import "snap-rq/internal/data"

type UserSessionEventListener interface {
	OnUserSessionLoadedLoaded(collections *[]data.Collection)
}

type UserSessionModel struct {
	userSessionData data.UserSession
	store           data.Store
}

func (m *UserSessionModel) Load() {
	m.userSessionData, _ = m.store.LoadSessionData()
}

func NewUserSessionModel(store data.Store) *UserSessionModel {
	return &UserSessionModel{
		store: store,
	}
}

func (m *UserSessionModel) GetUserSession() data.UserSession {
	return m.userSessionData
}
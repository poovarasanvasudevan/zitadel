package view

import (
	"github.com/caos/zitadel/internal/user/repository/view"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view"
)

const (
	userSessionTable = "auth.user_sessions"
)

func (v *View) UserSessionByID(sessionID string) (*model.UserSessionView, error) {
	return view.UserSessionByID(v.Db, userSessionTable, sessionID)
}

func (v *View) UserSessionByIDs(agentID, userID string) (*model.UserSessionView, error) {
	return view.UserSessionByIDs(v.Db, userSessionTable, agentID, userID)
}

func (v *View) UserSessionsByAgentID(agentID string) ([]*model.UserSessionView, error) {
	return view.UserSessionsByAgentID(v.Db, userSessionTable, agentID)
}

func (v *View) PutUserSession(userSession *model.UserSessionView) error {
	err := view.PutUserSession(v.Db, userSessionTable, userSession)
	if err != nil {
		return err
	}
	return v.ProcessedUserSessionSequence(userSession.Sequence)
}

func (v *View) DeleteUserSession(sessionID string, eventSequence uint64) error {
	err := view.DeleteUserSession(v.Db, userSessionTable, sessionID)
	if err != nil {
		return nil
	}
	return v.ProcessedUserSessionSequence(eventSequence)
}

func (v *View) GetLatestUserSessionSequence() (uint64, error) {
	return v.latestSequence(userSessionTable)
}

func (v *View) ProcessedUserSessionSequence(eventSequence uint64) error {
	return v.saveCurrentSequence(userSessionTable, eventSequence)
}

func (v *View) GetLatestUserSessionFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(userSessionTable, sequence)
}

func (v *View) ProcessedUserSessionFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
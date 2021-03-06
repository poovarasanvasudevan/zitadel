package view

import (
	"github.com/caos/zitadel/internal/token/repository/view"
	"github.com/caos/zitadel/internal/token/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	tokenTable = "auth.tokens"
)

func (v *View) TokenByID(tokenID string) (*model.Token, error) {
	return view.TokenByID(v.Db, tokenTable, tokenID)
}

func (v *View) IsTokenValid(tokenID string) (bool, error) {
	return view.IsTokenValid(v.Db, tokenTable, tokenID)
}

func (v *View) PutToken(token *model.Token) error {
	err := view.PutToken(v.Db, tokenTable, token)
	if err != nil {
		return err
	}
	return v.ProcessedTokenSequence(token.Sequence)
}

func (v *View) DeleteToken(tokenID string, eventSequence uint64) error {
	err := view.DeleteToken(v.Db, tokenTable, tokenID)
	if err != nil {
		return nil
	}
	return v.ProcessedTokenSequence(eventSequence)
}

func (v *View) DeleteSessionTokens(agentID, userID string, eventSequence uint64) error {
	err := view.DeleteSessionTokens(v.Db, tokenTable, agentID, userID)
	if err != nil {
		return nil
	}
	return v.ProcessedTokenSequence(eventSequence)
}

func (v *View) GetLatestTokenSequence() (*repository.CurrentSequence, error) {
	return v.latestSequence(tokenTable)
}

func (v *View) ProcessedTokenSequence(eventSequence uint64) error {
	return v.saveCurrentSequence(tokenTable, eventSequence)
}

func (v *View) GetLatestTokenFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(tokenTable, sequence)
}

func (v *View) ProcessedTokenFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}

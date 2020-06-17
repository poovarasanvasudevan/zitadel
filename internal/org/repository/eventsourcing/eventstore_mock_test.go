package eventsourcing

import (
	"context"
	"encoding/json"
	"time"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/mock"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	repo_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/golang/mock/gomock"
)

func GetMockedEventstore(ctrl *gomock.Controller, mockEs *mock.MockEventstore) *OrgEventstore {
	return &OrgEventstore{
		Eventstore: mockEs,
	}
}

func GetMockedOrgByIDOk(ctrl *gomock.Controller) *OrgEventstore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockedOrgByIDEventsOk(ctrl *gomock.Controller) *OrgEventstore {
	events := []*es_models.Event{
		&es_models.Event{Sequence: 6},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockedOrgByIDFilterFailedOk(ctrl *gomock.Controller) *OrgEventstore {
	events := []*es_models.Event{
		&es_models.Event{},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, errors.ThrowInternal(nil, "EVENT-SAa1O", "message"))
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockedDeactivateOrgPushFailed(ctrl *gomock.Controller) *OrgEventstore {
	startSequence := uint64(0)
	err := errors.ThrowInternal(nil, "EVENT-S8WzW", "test")

	events := []*es_models.Event{orgCreatedEventX()}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("test"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, aggregates ...*es_models.Aggregate) error {
			for _, aggregate := range aggregates {
				for _, event := range aggregate.Events {
					event.Sequence = startSequence
					startSequence++
				}
			}
			return err
		})
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockedDeactivateOrgPushCorrect(ctrl *gomock.Controller) *OrgEventstore {
	startSequence := uint64(6)

	events := []*es_models.Event{orgCreatedEventX()}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("test"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, aggregates ...*es_models.Aggregate) error {
			for _, aggregate := range aggregates {
				for _, event := range aggregate.Events {
					event.Sequence = startSequence
					startSequence++
				}
			}
			return nil
		})
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockedDeactivateOrgAlreadyInactive(ctrl *gomock.Controller) *OrgEventstore {
	startSequence := uint64(6)

	events := []*es_models.Event{orgCreatedEventX(), orgInactiveEventX()}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("test"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, aggregates ...*es_models.Aggregate) error {
			for _, aggregate := range aggregates {
				for _, event := range aggregate.Events {
					event.Sequence = startSequence
					startSequence++
				}
			}
			return nil
		})
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockedReactivateOrgPushFailed(ctrl *gomock.Controller) *OrgEventstore {
	startSequence := uint64(0)
	err := errors.ThrowInternal(nil, "EVENT-S8WzW", "test")

	events := []*es_models.Event{orgCreatedEvent(), orgInactiveEvent()}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("test"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, aggregates ...*es_models.Aggregate) error {
			for _, aggregate := range aggregates {
				for _, event := range aggregate.Events {
					event.Sequence = startSequence
					startSequence++
				}
			}
			return err
		})
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockedReactivateOrgPushCorrect(ctrl *gomock.Controller) *OrgEventstore {
	startSequence := uint64(6)

	events := []*es_models.Event{orgCreatedEventX(), orgInactiveEventX()}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("test"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, aggregates ...*es_models.Aggregate) error {
			for _, aggregate := range aggregates {
				for _, event := range aggregate.Events {
					event.Sequence = startSequence
					startSequence++
				}
			}
			return nil
		})
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockedReactivateOrgAlreadyInactive(ctrl *gomock.Controller) *OrgEventstore {
	startSequence := uint64(6)

	events := []*es_models.Event{orgCreatedEventX()}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("test"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, aggregates ...*es_models.Aggregate) error {
			for _, aggregate := range aggregates {
				for _, event := range aggregate.Events {
					event.Sequence = startSequence
					startSequence++
				}
			}
			return nil
		})
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockedOrgMemberByIDsNewEvents(ctrl *gomock.Controller) *OrgEventstore {
	events := []*es_models.Event{
		{Sequence: 6, Data: []byte("{\"userId\": \"banana\", \"roles\": [\"bananaa\"]}"), Type: model.OrgMemberChanged},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockedOrgMemberByIDsNoMember(ctrl *gomock.Controller) *OrgEventstore {
	events := []*es_models.Event{
		{Sequence: 6, Data: []byte("{\"userId\": \"banana\", \"roles\": [\"bananaa\"]}"), Type: model.OrgMemberAdded},
		{Sequence: 7, Data: []byte("{\"userId\": \"apple\"}"), Type: model.OrgMemberRemoved},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func AddOrgMemberMemberAlreadyExists(ctrl *gomock.Controller) *OrgEventstore {
	startSequence := uint64(0)
	err := errors.ThrowAlreadyExists(nil, "EVENT-yLTI6", "weiss nöd wie teste")

	events := []*es_models.Event{{
		Type:     model.OrgMemberAdded,
		Data:     []byte(`{"userId": "hodor", "roles": ["master"]}`),
		Sequence: 6,
	}}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("test"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, aggregates ...*es_models.Aggregate) error {
			for _, aggregate := range aggregates {
				for _, event := range aggregate.Events {
					event.Sequence = startSequence
					startSequence++
				}
			}
			return err
		})
	return GetMockedEventstore(ctrl, mockEs)
}
func AddOrgMemberMemberDeletedSuccess(ctrl *gomock.Controller) *OrgEventstore {
	startSequence := uint64(10)

	events := []*es_models.Event{{
		Type:     model.OrgMemberAdded,
		Data:     []byte(`{"userId": "hodor", "roles": ["master"]}`),
		Sequence: 6,
	},
		{
			Type:     model.OrgMemberRemoved,
			Data:     []byte(`{"userId": "hodor"}`),
			Sequence: 10,
		}}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("test"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, aggregates ...*es_models.Aggregate) error {
			for _, aggregate := range aggregates {
				for _, event := range aggregate.Events {
					event.Sequence = startSequence
					startSequence++
				}
			}
			return nil
		})
	return GetMockedEventstore(ctrl, mockEs)
}

func AddOrgMemberOrgNotExistsError(ctrl *gomock.Controller) *OrgEventstore {
	startSequence := uint64(10)
	err := errors.ThrowAlreadyExists(nil, "EVENT-yLTI6", "weiss nöd wie teste")

	events := []*es_models.Event{nil}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("test"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, aggregates ...*es_models.Aggregate) error {
			for _, aggregate := range aggregates {
				for _, event := range aggregate.Events {
					event.Sequence = startSequence
					startSequence++
				}
			}
			return err
		})
	return GetMockedEventstore(ctrl, mockEs)
}

func ChangeOrgMemberMemberNotFoundError(ctrl *gomock.Controller) *OrgEventstore {
	startSequence := uint64(0)

	events := []*es_models.Event{
		{
			AggregateID: "hodor-org",
			Type:        model.OrgAdded,
			Sequence:    4,
			Data:        []byte("{}"),
		},
		{
			AggregateID: "hodor-org",
			Type:        model.OrgMemberAdded,
			Data:        []byte(`{"userId": "brudi", "roles": ["master of desaster"]}`),
			Sequence:    6,
		},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("test"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, aggregates ...*es_models.Aggregate) error {
			for _, aggregate := range aggregates {
				for _, event := range aggregate.Events {
					event.Sequence = startSequence
					startSequence++
				}
			}
			return nil
		})
	return GetMockedEventstore(ctrl, mockEs)
}

func ChangeOrgMemberMemberFoundNoChangesError(ctrl *gomock.Controller) *OrgEventstore {
	startSequence := uint64(0)

	events := []*es_models.Event{
		{
			AggregateID: "hodor-org",
			Type:        model.OrgAdded,
			Sequence:    4,
			Data:        []byte("{}"),
		},
		{
			AggregateID: "hodor-org",
			Type:        model.OrgMemberAdded,
			Data:        []byte(`{"userId": "hodor", "roles": ["master"]}`),
			Sequence:    6,
		},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("test"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, aggregates ...*es_models.Aggregate) error {
			for _, aggregate := range aggregates {
				for _, event := range aggregate.Events {
					event.Sequence = startSequence
					startSequence++
				}
			}
			return nil
		})
	return GetMockedEventstore(ctrl, mockEs)
}

func ChangeOrgMemberPushError(ctrl *gomock.Controller) *OrgEventstore {
	startSequence := uint64(10)
	err := errors.ThrowInternal(nil, "PEVENT-3wqa2", "test")

	events := []*es_models.Event{
		{
			AggregateID: "hodor-org",
			Type:        model.OrgAdded,
			Sequence:    4,
			Data:        []byte("{}"),
		},
		{
			AggregateID: "hodor-org",
			Type:        model.OrgMemberAdded,
			Data:        []byte(`{"userId": "hodor", "roles": ["master"]}`),
			Sequence:    6,
		},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("test"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, aggregates ...*es_models.Aggregate) error {
			for _, aggregate := range aggregates {
				for _, event := range aggregate.Events {
					event.Sequence = startSequence
					startSequence++
				}
			}
			return err
		})
	return GetMockedEventstore(ctrl, mockEs)
}

func ChangeOrgMemberChangeSuccess(ctrl *gomock.Controller) *OrgEventstore {
	startSequence := uint64(7)

	events := []*es_models.Event{
		{
			AggregateID: "hodor-org",
			Type:        model.OrgAdded,
			Sequence:    4,
			Data:        []byte("{}"),
		},
		{
			AggregateID: "hodor-org",
			Type:        model.OrgMemberAdded,
			Data:        []byte(`{"userId": "hodor", "roles": ["master"]}`),
			Sequence:    6,
		},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("test"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, aggregates ...*es_models.Aggregate) error {
			for _, aggregate := range aggregates {
				for _, event := range aggregate.Events {
					event.Sequence = startSequence
					startSequence++
				}
			}
			return nil
		})
	return GetMockedEventstore(ctrl, mockEs)
}

func RemoveOrgMemberMemberNotFoundError(ctrl *gomock.Controller) *OrgEventstore {
	startSequence := uint64(0)

	events := []*es_models.Event{
		{
			AggregateID: "hodor-org",
			Type:        model.OrgAdded,
			Sequence:    4,
			Data:        []byte("{}"),
		},
		{
			AggregateID: "hodor-org",
			Type:        model.OrgMemberAdded,
			Data:        []byte(`{"userId": "brudi", "roles": ["master of desaster"]}`),
			Sequence:    6,
		},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	mockEs.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("test"))
	mockEs.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, aggregates ...*es_models.Aggregate) error {
			for _, aggregate := range aggregates {
				for _, event := range aggregate.Events {
					event.Sequence = startSequence
					startSequence++
				}
			}
			return nil
		})
	return GetMockedEventstore(ctrl, mockEs)
}

func orgCreatedEventX() *es_models.Event {
	return &es_models.Event{
		AggregateID:      "hodor-org",
		AggregateType:    model.OrgAggregate,
		AggregateVersion: "v1",
		CreationDate:     time.Now().Add(-1 * time.Minute),
		Data:             []byte(`{"name": "hodor-org", "domain":"hodor.org"}`),
		EditorService:    "testsvc",
		EditorUser:       "testuser",
		ID:               "sdlfö4t23kj",
		ResourceOwner:    "hodor-org",
		Sequence:         32,
		Type:             model.OrgAdded,
	}
}

func orgInactiveEventX() *es_models.Event {
	return &es_models.Event{
		AggregateID:      "hodor-org",
		AggregateType:    model.OrgAggregate,
		AggregateVersion: "v1",
		CreationDate:     time.Now().Add(-1 * time.Minute),
		Data:             nil,
		EditorService:    "testsvc",
		EditorUser:       "testuser",
		ID:               "sdlfö4t23kj",
		ResourceOwner:    "hodor-org",
		Sequence:         52,
		Type:             model.OrgDeactivated,
	}
}

func GetMockChangesOrgOK(ctrl *gomock.Controller) *OrgEventstore {
	org := model.Org{
		Name:   "MusterOrg",
		Domain: "myDomain",
	}
	data, err := json.Marshal(org)
	if err != nil {

	}
	events := []*es_models.Event{
		&es_models.Event{AggregateID: "AggregateIDApp", Sequence: 1, AggregateType: repo_model.OrgAggregate, Data: data},
	}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstore(ctrl, mockEs)
}

func GetMockChangesOrgNoEvents(ctrl *gomock.Controller) *OrgEventstore {
	events := []*es_models.Event{}
	mockEs := mock.NewMockEventstore(ctrl)
	mockEs.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, nil)
	return GetMockedEventstore(ctrl, mockEs)
}

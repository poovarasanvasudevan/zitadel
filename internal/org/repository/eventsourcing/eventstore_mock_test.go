package eventsourcing

import (
	"context"
	"encoding/json"
	"time"

	"github.com/caos/zitadel/internal/eventstore/mock"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/golang/mock/gomock"
)

type testOrgEventstore struct {
	OrgEventstore
	mockEventstore *mock.MockEventstore
}

func newTestEventstore(ctrl *gomock.Controller) *testOrgEventstore {
	mockEs := mock.NewMockEventstore(ctrl)
	return &testOrgEventstore{OrgEventstore: OrgEventstore{Eventstore: mockEs}, mockEventstore: mockEs}
}

func (es *testOrgEventstore) expectFilterEvents(events []*es_models.Event, err error) *testOrgEventstore {
	es.mockEventstore.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, err)
	return es
}

func (es *testOrgEventstore) expectPushEvents(startSequence uint64, err error) *testOrgEventstore {
	es.mockEventstore.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, aggregates ...*es_models.Aggregate) error {
			for _, aggregate := range aggregates {
				for _, event := range aggregate.Events {
					event.Sequence = startSequence
					startSequence++
				}
			}
			return err
		})
	return es
}

func (es *testOrgEventstore) expectAggregateCreator() *testOrgEventstore {
	es.mockEventstore.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("test"))
	return es
}

func GetMockedEventstore(ctrl *gomock.Controller, mockEs *mock.MockEventstore) *OrgEventstore {
	return &OrgEventstore{
		Eventstore: mockEs,
	}
}

func orgCreatedEvent() *es_models.Event {
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

func orgInactiveEvent() *es_models.Event {
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

func orgMemberAddedEvent() *es_models.Event {
	return &es_models.Event{
		AggregateID: "hodor-org",
		Type:        model.OrgMemberAdded,
		Data:        []byte(`{"userId": "hodor", "roles": ["master"]}`),
		Sequence:    6,
	}
}

func orgMemberRemovedEvent() *es_models.Event {
	return &es_models.Event{
		Sequence: 7,
		Data:     []byte("{\"userId\": \"apple\"}"),
		Type:     model.OrgMemberRemoved,
	}
}

func orgChangesEvent() *es_models.Event {
	org := model.Org{
		Name:   "MusterOrg",
		Domain: "myDomain",
	}
	data, err := json.Marshal(org)
	if err != nil {

	}

	return &es_models.Event{AggregateID: "AggregateIDApp", Sequence: 1, AggregateType: model.OrgAggregate, Data: data}
}

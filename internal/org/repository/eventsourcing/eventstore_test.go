package eventsourcing

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"

	"github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_mock "github.com/caos/zitadel/internal/eventstore/mock"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/golang/mock/gomock"
)

type testOrgEventstore struct {
	OrgEventstore
	mockEventstore *es_mock.MockEventstore
}

func newTestEventstore(t *testing.T) *testOrgEventstore {
	mock := mockEventstore(t)
	return &testOrgEventstore{OrgEventstore: OrgEventstore{Eventstore: mock}, mockEventstore: mock}
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

func mockEventstore(t *testing.T) *es_mock.MockEventstore {
	ctrl := gomock.NewController(t)
	e := es_mock.NewMockEventstore(ctrl)

	return e
}

func TestOrgEventstore_OrgByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	type res struct {
		expectedSequence uint64
		isErr            func(error) bool
	}
	type args struct {
		es  *OrgEventstore
		ctx context.Context
		org *org_model.Org
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "no input org",
			args: args{
				es:  GetMockedOrgByIDOk(ctrl),
				ctx: auth.NewMockContext("user", "org"),
				org: nil,
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "no aggregate id in input org",
			args: args{
				es:  GetMockedOrgByIDOk(ctrl),
				ctx: auth.NewMockContext("user", "org"),
				org: &org_model.Org{ObjectRoot: es_models.ObjectRoot{Sequence: 4}},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsPreconditionFailed,
			},
		},
		{
			name: "no events found success",
			args: args{
				es:  GetMockedOrgByIDOk(ctrl),
				ctx: auth.NewMockContext("user", "org"),
				org: &org_model.Org{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "hodor"}},
			},
			res: res{
				expectedSequence: 4,
				isErr:            nil,
			},
		},
		{
			name: "filter fail",
			args: args{
				es:  GetMockedOrgByIDFilterFailedOk(ctrl),
				ctx: auth.NewMockContext("user", "org"),
				org: &org_model.Org{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "hodor"}},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsInternal,
			},
		},
		{
			name: "new events found and added success",
			args: args{
				es:  GetMockedOrgByIDEventsOk(ctrl),
				ctx: auth.NewMockContext("user", "org"),
				org: &org_model.Org{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "hodor-org", ChangeDate: time.Now(), CreationDate: time.Now()}},
			},
			res: res{
				expectedSequence: 6,
				isErr:            nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.es.OrgByID(tt.args.ctx, tt.args.org)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got:%T %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if got == nil && tt.res.expectedSequence != 0 {
				t.Errorf("org should be nil but was %v", got)
				t.FailNow()
			}
			if tt.res.expectedSequence != 0 && tt.res.expectedSequence != got.Sequence {
				t.Errorf("org should have sequence %d but had %d", tt.res.expectedSequence, got.Sequence)
			}
		})
	}
}

func TestOrgEventstore_DeactivateOrg(t *testing.T) {
	ctrl := gomock.NewController(t)
	type res struct {
		expectedSequence uint64
		isErr            func(error) bool
	}
	type args struct {
		es    *OrgEventstore
		ctx   context.Context
		orgID string
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "no input org",
			args: args{
				es:    GetMockedOrgByIDOk(ctrl),
				ctx:   auth.NewMockContext("user", "org"),
				orgID: "",
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "push failed",

			args: args{
				es:    GetMockedDeactivateOrgPushFailed(ctrl),
				ctx:   auth.NewMockContext("user", "org"),
				orgID: "hodor",
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsInternal,
			},
		},
		{
			name: "push correct",
			args: args{
				es:    GetMockedDeactivateOrgPushCorrect(ctrl),
				ctx:   auth.NewMockContext("user", "org"),
				orgID: "hodor",
			},
			res: res{
				expectedSequence: 6,
				isErr:            nil,
			},
		},
		{
			name: "org already inactive error",
			args: args{
				es:    GetMockedDeactivateOrgAlreadyInactive(ctrl),
				ctx:   auth.NewMockContext("user", "org"),
				orgID: "hodor",
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorInvalidArgument,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.es.DeactivateOrg(tt.args.ctx, tt.args.orgID)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got:%T %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if got == nil && tt.res.expectedSequence != 0 {
				t.Errorf("org should be nil but was %v", got)
				t.FailNow()
			}
			if tt.res.expectedSequence != 0 && tt.res.expectedSequence != got.Sequence {
				t.Errorf("org should have sequence %d but had %d", tt.res.expectedSequence, got.Sequence)
			}
		})
	}
}

func TestOrgEventstore_ReactivateOrg(t *testing.T) {
	ctrl := gomock.NewController(t)
	type res struct {
		expectedSequence uint64
		isErr            func(error) bool
	}
	type args struct {
		es    *OrgEventstore
		ctx   context.Context
		orgID string
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "no input org",
			args: args{
				es:    GetMockedOrgByIDOk(ctrl),
				ctx:   auth.NewMockContext("user", "org"),
				orgID: "",
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "push failed",
			args: args{
				es:    GetMockedReactivateOrgPushFailed(ctrl),
				ctx:   auth.NewMockContext("user", "org"),
				orgID: "hodor",
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsInternal,
			},
		},
		{
			name: "push correct",
			args: args{
				es:    GetMockedReactivateOrgPushCorrect(ctrl),
				ctx:   auth.NewMockContext("user", "org"),
				orgID: "hodor",
			},
			res: res{
				expectedSequence: 6,
				isErr:            nil,
			},
		},
		{
			name: "org already active error",
			args: args{
				es:    GetMockedReactivateOrgAlreadyInactive(ctrl),
				ctx:   auth.NewMockContext("user", "org"),
				orgID: "hodor",
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorInvalidArgument,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.es.ReactivateOrg(tt.args.ctx, tt.args.orgID)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got:%T %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if got == nil && tt.res.expectedSequence != 0 {
				t.Errorf("org should be nil but was %v", got)
				t.FailNow()
			}
			if tt.res.expectedSequence != 0 && tt.res.expectedSequence != got.Sequence {
				t.Errorf("org should have sequence %d but had %d", tt.res.expectedSequence, got.Sequence)
			}
		})
	}
}

func TestOrgEventstore_OrgMemberByIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	type res struct {
		expectedSequence uint64
		isErr            func(error) bool
	}
	type args struct {
		es     *OrgEventstore
		ctx    context.Context
		member *org_model.OrgMember
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "no input member",
			args: args{
				es:     GetMockedOrgByIDOk(ctrl),
				ctx:    auth.NewMockContext("user", "org"),
				member: nil,
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsPreconditionFailed,
			},
		},
		{
			name: "no aggregate id in input member",
			args: args{
				es:     GetMockedOrgByIDOk(ctrl),
				ctx:    auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{ObjectRoot: es_models.ObjectRoot{Sequence: 4}, UserID: "asdf"},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsPreconditionFailed,
			},
		},
		{
			name: "no aggregate id in input member",
			args: args{
				es:     GetMockedOrgByIDOk(ctrl),
				ctx:    auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "asdf"}},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsPreconditionFailed,
			},
		},
		{
			name: "no events found success",
			args: args{
				es:     GetMockedOrgByIDOk(ctrl),
				ctx:    auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "plants"}, UserID: "banana"},
			},
			res: res{
				expectedSequence: 4,
				isErr:            nil,
			},
		},
		{
			name: "filter fail",
			args: args{
				es:     GetMockedOrgByIDFilterFailedOk(ctrl),
				ctx:    auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "plants"}, UserID: "banana"},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsInternal,
			},
		},
		{
			name: "new events found and added success",
			args: args{
				es:     GetMockedOrgMemberByIDsNewEvents(ctrl),
				ctx:    auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "plants", ChangeDate: time.Now(), CreationDate: time.Now()}, UserID: "banana"},
			},
			res: res{
				expectedSequence: 6,
				isErr:            nil,
			},
		},
		{
			name: "not member of org error",
			args: args{
				es:     GetMockedOrgMemberByIDsNoMember(ctrl),
				ctx:    auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "plants", ChangeDate: time.Now(), CreationDate: time.Now()}, UserID: "apple"},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.es.OrgMemberByIDs(tt.args.ctx, tt.args.member)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got:%T %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if got == nil && tt.res.expectedSequence != 0 {
				t.Errorf("org should be nil but was %v", got)
				t.FailNow()
			}
			if tt.res.expectedSequence != 0 && tt.res.expectedSequence != got.Sequence {
				t.Errorf("org should have sequence %d but had %d", tt.res.expectedSequence, got.Sequence)
			}
		})
	}
}

func TestOrgEventstore_AddOrgMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	type res struct {
		expectedSequence uint64
		isErr            func(error) bool
	}
	type args struct {
		es     *OrgEventstore
		ctx    context.Context
		member *org_model.OrgMember
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "no input member",
			args: args{
				es:     GetMockedOrgByIDOk(ctrl),
				ctx:    auth.NewMockContext("user", "org"),
				member: nil,
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsPreconditionFailed,
			},
		},
		{
			name: "push failed",
			args: args{
				es:  GetMockedDeactivateOrgPushFailed(ctrl),
				ctx: auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{
						Sequence:    4,
						AggregateID: "hodor-org",
					},
					UserID: "hodor",
					Roles:  []string{"nix"},
				},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsInternal,
			},
		},
		{
			name: "push correct",
			args: args{
				es:  GetMockedDeactivateOrgPushCorrect(ctrl),
				ctx: auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{
						Sequence:    4,
						AggregateID: "hodor-org",
					},
					UserID: "hodor",
					Roles:  []string{"nix"},
				},
			},
			res: res{
				expectedSequence: 6,
				isErr:            nil,
			},
		},
		{
			name: "member already exists error",
			args: args{
				es:  AddOrgMemberMemberAlreadyExists(ctrl),
				ctx: auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{
						Sequence:    4,
						AggregateID: "hodor-org",
					},
					UserID: "hodor",
					Roles:  []string{"nix"},
				},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorAlreadyExists,
			},
		},
		{
			name: "member deleted success",
			args: args{
				es:  AddOrgMemberMemberDeletedSuccess(ctrl),
				ctx: auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{
						Sequence:    4,
						AggregateID: "hodor-org",
					},
					UserID: "hodor",
					Roles:  []string{"nix"},
				},
			},
			res: res{
				expectedSequence: 10,
				isErr:            nil,
			},
		},
		{
			name: "org not exists error",
			args: args{
				es:  AddOrgMemberOrgNotExistsError(ctrl),
				ctx: auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{
						Sequence:    4,
						AggregateID: "hodor-org",
					},
					UserID: "hodor",
					Roles:  []string{"nix"},
				},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorAlreadyExists,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.es.AddOrgMember(tt.args.ctx, tt.args.member)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got:%T %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if got == nil && tt.res.expectedSequence != 0 {
				t.Errorf("org should not be nil but was %v", got)
				t.FailNow()
			}
			if tt.res.expectedSequence != 0 && tt.res.expectedSequence != got.Sequence {
				t.Errorf("org should have sequence %d but had %d", tt.res.expectedSequence, got.Sequence)
			}
		})
	}
}

func TestOrgEventstore_ChangeOrgMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	type res struct {
		isErr            func(error) bool
		expectedSequence uint64
	}
	type args struct {
		es     *OrgEventstore
		ctx    context.Context
		member *org_model.OrgMember
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "no input member",
			args: args{
				es:     GetMockedOrgByIDOk(ctrl),
				ctx:    auth.NewMockContext("user", "org"),
				member: nil,
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsPreconditionFailed,
			},
		},
		{
			name: "member not found error",
			args: args{
				es:  ChangeOrgMemberMemberNotFoundError(ctrl),
				ctx: auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org", Sequence: 5},
					UserID:     "hodor",
					Roles:      []string{"master"},
				},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsNotFound,
			},
		},
		{
			name: "member found no changes error",

			args: args{
				es:  ChangeOrgMemberMemberFoundNoChangesError(ctrl),
				ctx: auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org", Sequence: 5},
					UserID:     "hodor",
					Roles:      []string{"master"},
				},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "push error",
			args: args{
				es:  ChangeOrgMemberPushError(ctrl),
				ctx: auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org", Sequence: 5},
					UserID:     "hodor",
					Roles:      []string{"master of desaster"},
				},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsInternal,
			},
		},
		{
			name: "change success",
			args: args{
				es:  ChangeOrgMemberChangeSuccess(ctrl),
				ctx: auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org", Sequence: 5},
					UserID:     "hodor",
					Roles:      []string{"master of desaster"},
				},
			},
			res: res{
				expectedSequence: 7,
				isErr:            nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.es.ChangeOrgMember(tt.args.ctx, tt.args.member)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got:%T %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if got == nil && tt.res.expectedSequence != 0 {
				t.Errorf("org should not be nil but was %v", got)
				t.FailNow()
			}
			if tt.res.expectedSequence != 0 && tt.res.expectedSequence != got.Sequence {
				t.Errorf("org should have sequence %d but had %d", tt.res.expectedSequence, got.Sequence)
			}
		})
	}
}

func TestOrgEventstore_RemoveOrgMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	type res struct {
		isErr func(error) bool
	}
	type args struct {
		es     *OrgEventstore
		ctx    context.Context
		member *org_model.OrgMember
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "no input member",
			args: args{
				es:     GetMockedOrgByIDOk(ctrl),
				ctx:    auth.NewMockContext("user", "org"),
				member: nil,
			},
			res: res{
				isErr: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "member not found error",

			args: args{
				es:  RemoveOrgMemberMemberNotFoundError(ctrl),
				ctx: auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org", Sequence: 5},
					UserID:     "hodor",
					Roles:      []string{"master"},
				},
			},
			res: res{
				isErr: nil,
			},
		},
		{
			name: "push error",
			args: args{
				es:  ChangeOrgMemberPushError(ctrl),
				ctx: auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org", Sequence: 5},
					UserID:     "hodor",
				},
			},
			res: res{
				isErr: errors.IsInternal,
			},
		},
		{
			name: "remove success",
			args: args{
				es:  ChangeOrgMemberChangeSuccess(ctrl),
				ctx: auth.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org", Sequence: 5},
					UserID:     "hodor",
				},
			},
			res: res{
				isErr: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RemoveOrgMember(tt.args.ctx, tt.args.member)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got:%T %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
		})
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

func TestChangesOrg(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es           *OrgEventstore
		id           string
		lastSequence uint64
		limit        uint64
	}
	type res struct {
		changes *org_model.OrgChanges
		org     *model.Org
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "changes from events, ok",
			args: args{
				es:           GetMockChangesOrgOK(ctrl),
				id:           "1",
				lastSequence: 0,
				limit:        0,
			},
			res: res{
				changes: &org_model.OrgChanges{Changes: []*org_model.OrgChange{&org_model.OrgChange{EventType: "", Sequence: 1, Modifier: ""}}, LastSequence: 1},
				org:     &model.Org{Name: "MusterOrg", Domain: "myDomain"},
			},
		},
		{
			name: "changes from events, no events",
			args: args{
				es:           GetMockChangesOrgNoEvents(ctrl),
				id:           "2",
				lastSequence: 0,
				limit:        0,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.OrgChanges(nil, tt.args.id, tt.args.lastSequence, tt.args.limit)

			org := &model.Org{}
			if result != nil && len(result.Changes) > 0 {
				b, err := json.Marshal(result.Changes[0].Data)
				json.Unmarshal(b, org)
				if err != nil {
				}
			}
			if !tt.res.wantErr && result.LastSequence != tt.res.changes.LastSequence && org.Name != tt.res.org.Name {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.changes.LastSequence, result.LastSequence)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

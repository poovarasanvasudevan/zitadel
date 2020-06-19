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
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/golang/mock/gomock"
)

func TestOrgEventstore_OrgByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	type res struct {
		expectedSequence uint64
		isErr            func(error) bool
	}
	type args struct {
		es  *testOrgEventstore
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
				es:  newTestEventstore(ctrl),
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
				es:  newTestEventstore(ctrl),
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
				es:  newTestEventstore(ctrl).expectFilterEvents([]*es_models.Event{}, nil),
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
				es:  newTestEventstore(ctrl).expectFilterEvents([]*es_models.Event{}, errors.ThrowInternal(nil, "EVENT-SAa1O", "message")),
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
				es:  newTestEventstore(ctrl).expectFilterEvents([]*es_models.Event{orgCreatedSimpleEvent()}, nil),
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
		es    *testOrgEventstore
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
				es:    newTestEventstore(ctrl),
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
				es: newTestEventstore(ctrl).expectFilterEvents([]*es_models.Event{orgCreatedEvent()}, nil).
					expectAggregateCreator().expectAggregateCreator().
					expectPushEvents(0, errors.ThrowInternal(nil, "EVENT-S8WzW", "test")),
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
				es: newTestEventstore(ctrl).expectFilterEvents([]*es_models.Event{orgCreatedEvent()}, nil).
					expectAggregateCreator().
					expectPushEvents(6, nil),
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
				es: newTestEventstore(ctrl).
					expectFilterEvents([]*es_models.Event{orgCreatedEvent(), orgInactiveEvent()}, nil).
					expectAggregateCreator().
					expectPushEvents(6, nil),
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
		es    *testOrgEventstore
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
				es:    newTestEventstore(ctrl),
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
				es: newTestEventstore(ctrl).
					expectFilterEvents([]*es_models.Event{orgCreatedEvent(), orgInactiveEvent()}, nil).
					expectAggregateCreator().
					expectPushEvents(0, errors.ThrowInternal(nil, "EVENT-S8WzW", "test")),
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
				es: newTestEventstore(ctrl).
					expectFilterEvents([]*es_models.Event{orgCreatedEvent(), orgInactiveEvent()}, nil).
					expectAggregateCreator().
					expectPushEvents(6, nil),
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
				es: newTestEventstore(ctrl).
					expectFilterEvents([]*es_models.Event{orgCreatedEvent()}, nil).
					expectAggregateCreator().
					expectPushEvents(6, nil),
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
		es     *testOrgEventstore
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
				es:     newTestEventstore(ctrl),
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
				es:     newTestEventstore(ctrl),
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
				es:     newTestEventstore(ctrl),
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
				es: newTestEventstore(ctrl).expectAggregateCreator().
					expectFilterEvents([]*es_models.Event{}, nil),
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
				es: newTestEventstore(ctrl).
					expectFilterEvents([]*es_models.Event{}, errors.ThrowInternal(nil, "EVENT-SAa1O", "message")),
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
				es: newTestEventstore(ctrl).
					expectFilterEvents([]*es_models.Event{orgMemberChangedEvent()}, nil),
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
				es:     newTestEventstore(ctrl).expectFilterEvents([]*es_models.Event{orgMemberAddedEvent(), orgMemberRemovedEvent()}, nil),
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
		es     *testOrgEventstore
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
				es:     newTestEventstore(ctrl),
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
				es: newTestEventstore(ctrl).
					expectFilterEvents([]*es_models.Event{orgCreatedEvent()}, nil).
					expectAggregateCreator().
					expectPushEvents(0, errors.ThrowInternal(nil, "EVENT-S8WzW", "test")),
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
				es: newTestEventstore(ctrl).
					expectAggregateCreator().
					expectFilterEvents([]*es_models.Event{orgCreatedEvent()}, nil).
					expectPushEvents(6, nil),
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
				es: newTestEventstore(ctrl).
					expectAggregateCreator().
					expectFilterEvents([]*es_models.Event{orgMemberAddedEvent()}, nil).
					expectPushEvents(0, errors.ThrowAlreadyExists(nil, "EVENT-yLTI6", "weiss nöd wie teste")),
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
				es: newTestEventstore(ctrl).
					expectAggregateCreator().
					expectPushEvents(10, nil).
					expectFilterEvents([]*es_models.Event{orgMemberAddedEvent(), orgMemberRemovedEvent()}, nil),
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
				es: newTestEventstore(ctrl).
					expectAggregateCreator().
					expectFilterEvents(nil, nil).
					expectPushEvents(0, errors.ThrowAlreadyExists(nil, "EVENT-yLTI6", "weiss nöd wie teste")),
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
		es     *testOrgEventstore
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
				es:     newTestEventstore(ctrl),
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
				es: newTestEventstore(ctrl).
					expectAggregateCreator().
					expectFilterEvents([]*es_models.Event{orgMemberOrgAddedEvent(), orgMemberAddedEventNoMemberFound()}, nil),
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
				es: newTestEventstore(ctrl).
					expectAggregateCreator().
					expectFilterEvents([]*es_models.Event{orgCreatedEvent(), orgMemberAddedEvent()}, nil),
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
				es: newTestEventstore(ctrl).
					expectAggregateCreator().
					expectFilterEvents([]*es_models.Event{orgCreatedEvent(), orgMemberAddedEvent()}, nil).
					expectPushEvents(0, errors.ThrowInternal(nil, "PEVENT-3wqa2", "test")),
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
				es: newTestEventstore(ctrl).
					expectAggregateCreator().
					expectFilterEvents([]*es_models.Event{orgCreatedEvent(), orgMemberAddedEvent()}, nil).
					expectPushEvents(7, nil),
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
		es     *testOrgEventstore
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
				es:     newTestEventstore(ctrl),
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
				es: newTestEventstore(ctrl).
					expectAggregateCreator().
					expectFilterEvents([]*es_models.Event{orgMemberOrgAddedEvent(), orgMemberAddedEventNoMemberFound()}, nil),
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
				es: newTestEventstore(ctrl).
					expectAggregateCreator().
					expectFilterEvents([]*es_models.Event{orgMemberOrgAddedEvent(), orgMemberAddedEvent()}, nil).
					expectPushEvents(0, errors.ThrowInternal(nil, "PEVENT-3wqa2", "test")),
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
				es: newTestEventstore(ctrl).
					expectAggregateCreator().
					expectFilterEvents([]*es_models.Event{orgMemberOrgAddedEvent(), orgMemberAddedEvent()}, nil).
					expectPushEvents(7, nil),
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

func TestChangesOrg(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es           *testOrgEventstore
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
				es:           newTestEventstore(ctrl).expectFilterEvents([]*es_models.Event{orgChangesEvent()}, nil),
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
				es:           newTestEventstore(ctrl).expectFilterEvents([]*es_models.Event{}, nil),
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

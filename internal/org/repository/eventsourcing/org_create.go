package eventsourcing

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/org/model"
)

func CreateOrg(org *model.Org, creator *es_models.AggregateCreator) []es_models.AggregateStruct {
	return []es_models.AggregateStruct{
		orgCreateFromModel(org, creator),
		domainReservedFromModel(org, creator),
		nameReservedFromModel(org, creator),
	}
}

func AddMember(member *model.OrgMember, creator *es_models.AggregateCreator) []es_models.AggregateStruct {
	return []es_models.AggregateStruct{
		NewOrgMemberAddFromModel(member, creator),
	}
}

func orgCreateFromModel(org *model.Org, creator *es_models.AggregateCreator) *OrgCreateEvent {
	id, err := idGenerator.NextID()
	logging.Log("EVENT-tYVI5").OnError(err).Fatal("id get failed")

	org.ObjectRoot.AggregateID = strconv.FormatUint(id, 10)
	return &OrgCreateEvent{
		creator: creator,
		Domain:  org.Domain,
		Name:    org.Name,
		ObjectRoot: es_models.ObjectRoot{
			AggregateID: strconv.FormatUint(id, 10),
		},
	}
}

type OrgCreateEvent struct {
	creator *es_models.AggregateCreator

	es_models.ObjectRoot
	Name   string `json:"name,omitempty"`
	Domain string `json:"domain,omitempty"`
}

func (org *OrgCreateEvent) ValidationQuery() *es_models.SearchQuery {
	return es_models.NewSearchQuery().AggregateIDFilter(org.ObjectRoot.AggregateID).AggregateTypeFilter(model.OrgAggregate).OrderDesc()
}

func (org *OrgCreateEvent) Validate(events ...*es_models.Event) error {
	// if org.Name == "" || org.Domain == "" {
	// 	return errors.ThrowInvalidArgument(nil, "EVENT-3Zs9L", "name or domain invalid")
	// }
	if len(events) == 0 {
		return nil
	}
	return errors.ThrowAlreadyExists(nil, "EVENT-cDU5v", "org already exists")
}

func (org *OrgCreateEvent) ToEvents(ctx context.Context) ([]*es_models.Event, error) {
	orgAgg, err := org.creator.NewAggregate(ctx, org.ObjectRoot.AggregateID, model.OrgAggregate, orgVersion, 0)
	if err != nil {
		return nil, err
	}

	orgAgg, err = orgAgg.AppendEvent(model.OrgAdded, org)
	if err != nil {
		return nil, err
	}
	return orgAgg.Events, nil
}

func (org *OrgCreateEvent) PreviousSequence() uint64 {
	return 0
}

func domainReservedFromModel(org *model.Org, creator *es_models.AggregateCreator) *domainReservedEvent {
	return &domainReservedEvent{
		creator: creator,
		domain:  org.Domain,
	}
}

type domainReservedEvent struct {
	creator *es_models.AggregateCreator

	domain           string
	previousSequence uint64
}

func (domain *domainReservedEvent) ValidationQuery() *es_models.SearchQuery {
	return es_models.NewSearchQuery().AggregateIDFilter(domain.domain).AggregateTypeFilter(model.OrgDomainAggregate).OrderDesc().SetLimit(1)
}

func (domain *domainReservedEvent) Validate(events ...*es_models.Event) error {
	if domain.domain == "" {
		return errors.ThrowInvalidArgument(nil, "EVENT-3Zs9L", "domain invalid")
	}
	if len(events) == 0 {
		return nil
	}
	if events[0].Type == model.OrgDomainReserved {
		return errors.ThrowAlreadyExists(nil, "EVENT-cDU5v", "domain reserved")
	}
	domain.previousSequence = events[0].Sequence
	return nil
}

func (domain *domainReservedEvent) ToEvents(ctx context.Context) ([]*es_models.Event, error) {
	//TODO: previous sequence is not needed in this case
	agg, err := domain.creator.NewAggregate(ctx, domain.domain, model.OrgDomainAggregate, orgVersion, domain.previousSequence)
	if err != nil {
		return nil, err
	}
	agg, err = agg.AppendEvent(model.OrgDomainReserved, nil)
	if err != nil {
		return nil, err
	}
	return agg.Events, nil
}

func (domain *domainReservedEvent) PreviousSequence() uint64 {
	return domain.previousSequence
}

func nameReservedFromModel(org *model.Org, creator *es_models.AggregateCreator) *nameReservedEvent {
	return &nameReservedEvent{
		creator: creator,
		name:    org.Name,
	}
}

type nameReservedEvent struct {
	creator *es_models.AggregateCreator

	name             string
	previousSequence uint64
}

func (name *nameReservedEvent) ValidationQuery() *es_models.SearchQuery {
	return es_models.NewSearchQuery().AggregateIDFilter(name.name).AggregateTypeFilter(model.OrgNameAggregate).OrderDesc()
}

func (name *nameReservedEvent) Validate(events ...*es_models.Event) error {
	if name.name == "" {
		return errors.ThrowInvalidArgument(nil, "EVENT-ggapz", "domain invalid")
	}
	if len(events) == 0 {
		return nil
	}
	if events[0].Type == model.OrgNameReserved {
		return errors.ThrowAlreadyExists(nil, "EVENT-tOP3w", "domain reserved")
	}
	name.previousSequence = events[0].Sequence
	return nil
}

func (name *nameReservedEvent) ToEvents(ctx context.Context) ([]*es_models.Event, error) {
	agg, err := name.creator.NewAggregate(ctx, name.name, model.OrgNameAggregate, orgVersion, name.previousSequence)
	if err != nil {
		return nil, err
	}
	agg, err = agg.AppendEvent(model.OrgNameReserved, nil)
	if err != nil {
		return nil, err
	}

	return agg.Events, nil
}

func (name *nameReservedEvent) PreviousSequence() uint64 {
	return name.previousSequence
}

func OrgMemberAddFromModel(member *model.OrgMember, creator *es_models.AggregateCreator) *orgMemberAddEvent {
	return &orgMemberAddEvent{
		creator:    creator,
		ObjectRoot: member.ObjectRoot,
		Roles:      member.Roles,
		UserID:     member.UserID,
	}
}

func NewOrgMemberAddFromModel(member *model.OrgMember, creator *es_models.AggregateCreator) *orgMemberAddEvent {
	return &orgMemberAddEvent{
		creator:    creator,
		ObjectRoot: member.ObjectRoot,
		Roles:      member.Roles,
		UserID:     member.UserID,
	}
}

type orgMemberAddEvent struct {
	creator *es_models.AggregateCreator

	//ObjectRoot must be the root of org
	es_models.ObjectRoot `json:"-"`

	UserID string   `json:"userId"`
	Roles  []string `json:"roles"`

	previousSequence uint64
}

func (member *orgMemberAddEvent) ValidationQuery() *es_models.SearchQuery {
	return es_models.NewSearchQuery().AggregateIDFilter(member.AggregateID).AggregateTypeFilter(model.OrgAggregate)
}

func (member *orgMemberAddEvent) Validate(events ...*es_models.Event) error {
	if member.AggregateID == "" || member.UserID == "" {
		return errors.ThrowInvalidArgument(nil, "EVENT-3b9oL", "member must contain orgID and userID")
	}

	if len(events) == 0 {
		return errors.ThrowNotFound(nil, "EVENT-cfsbr", "org not found")
	}

	isMember := false
	existsUser := false
	_ = existsUser
	for _, event := range events {
		member.ObjectRoot.AppendEvent(event)

		if event.Type == model.OrgMemberAdded {
			userID := userIDfromEvent(event)
			if userID.UserID == member.UserID {
				isMember = true
			}
		}
		if event.Type == model.OrgMemberRemoved {
			userID := userIDfromEvent(event)
			if userID.UserID == member.UserID {
				isMember = false
			}
		}
	}

	if isMember {
		return errors.ThrowAlreadyExists(nil, "EVENT-cDU5v", "user is already member of org")
	}
	member.previousSequence = events[len(events)-1].Sequence
	return nil
}

func (member *orgMemberAddEvent) ToEvents(ctx context.Context) ([]*es_models.Event, error) {
	orgAgg, err := member.creator.NewAggregate(ctx, member.AggregateID, model.OrgAggregate, orgVersion, member.Sequence)
	if err != nil {
		return nil, err
	}
	orgAgg, err = orgAgg.AppendEvent(model.OrgMemberAdded, member)
	if err != nil {
		return nil, err
	}

	return orgAgg.Events, nil
}

func (member *orgMemberAddEvent) PreviousSequence() uint64 {
	return member.previousSequence
}

type orgMemberRemoveEvent struct {
	creator *es_models.AggregateCreator

	//ObjectRoot must be the root of org
	es_models.ObjectRoot `json:"-"`

	UserID string `json:"userId"`
}

func (member *orgMemberRemoveEvent) ValidationQuery() *es_models.SearchQuery {
	return es_models.NewSearchQuery().AggregateIDFilter(member.AggregateID).AggregateTypeFilter(model.OrgAggregate).OrderDesc()
}

func (member *orgMemberRemoveEvent) Validate(events ...*es_models.Event) error {
	if len(events) == 0 {
		return errors.ThrowNotFound(nil, "EVENT-hJkKk", "org not found")
	}

	member.ObjectRoot.AppendEvent(events[len(events)-1])

	isMember := false
	for _, event := range events {
		if event.Type == model.OrgMemberAdded {
			userID := userIDfromEvent(event)
			if userID.UserID == member.UserID {
				isMember = true
			}
		}
		if event.Type == model.OrgMemberRemoved {
			userID := userIDfromEvent(event)
			if userID.UserID == member.UserID {
				isMember = false
			}
		}
	}

	if !isMember {
		return errors.ThrowAlreadyExists(nil, "EVENT-cDU5v", "user is not member of org")
	}
	return nil
}

func (member *orgMemberRemoveEvent) ToEvents(ctx context.Context) ([]*es_models.Event, error) {
	orgAgg, err := member.creator.NewAggregate(ctx, member.AggregateID, model.OrgAggregate, orgVersion, member.Sequence)
	if err != nil {
		return nil, err
	}
	orgAgg, err = orgAgg.AppendEvent(model.OrgMemberAdded, member)
	if err != nil {
		return nil, err
	}

	return orgAgg.Events, nil
}

type memberUserID struct {
	UserID string `json:"userId"`
}

func userIDfromEvent(event *es_models.Event) *memberUserID {
	userID := new(memberUserID)
	err := json.Unmarshal(event.Data, userID)
	logging.Log("EVENT-bIKRN").OnError(err).Warn("error unmarshal add member")
	return userID
}

// type OrgUpdate struct {
// 	creator *es_models.AggregateCreator

// 	es_models.ObjectRoot `json:"-"`

// 	changes map[string]interface{}
// }

// func (c *OrgUpdate) ValidationQuery() *es_models.SearchQuery {
// 	return es_models.NewSearchQuery().AggregateIDFilter(c.AggregateID).AggregateTypeFilter(orgType).OrderDesc()
// }

// func (c *OrgUpdate) Validate(events ...*es_models.Event) error {
// 	if len(events) == 0 {
// 		return errors.ThrowNotFound(nil, "EVENT-tCUd4", "org not found")
// 	}
// 	for _, event := range events {
// 		if event.Type == model.OrgDeactivated {
// 			return errors.ThrowPreconditionFailed(nil, "EVENT-Q5tCh", "org inactive")
// 		}
// 		if event.Type == model.OrgReactivated {
// 			return nil
// 		}
// 	}
// 	return nil
// }

// func (c *OrgUpdate) ToAggregate(ctx context.Context) (*es_models.Aggregate, error) {
// 	id, err := idGenerator.NextID()
// 	if err != nil {
// 		return nil, errors.ThrowInternal(err, "EVENT-QFCqb", "id gen failed")
// 	}
// 	orgAgg, err := c.creator.NewAggregate(ctx, strconv.FormatUint(id, 10), orgType, orgVersion, 0)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return []*es_models.Aggregate{
// 		orgAgg,
// 	}, nil
// }

// type OrgDeactivate struct {
// 	es_models.ObjectRoot `json:"-"`
// }

// type OrgReactivate struct {
// 	es_models.ObjectRoot `json:"-"`
// }

type UserCreateEvent struct {
	UserName string
}

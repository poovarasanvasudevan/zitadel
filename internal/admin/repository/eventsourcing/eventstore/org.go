package eventstore

import (
	"context"

	"github.com/caos/logging"
	admin_model "github.com/caos/zitadel/internal/admin/model"
	admin_view "github.com/caos/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/eventstore"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/sdk"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_es "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	"github.com/caos/zitadel/internal/org/repository/view/model"
	policy_model "github.com/caos/zitadel/internal/policy/model"
	policy_es "github.com/caos/zitadel/internal/policy/repository/eventsourcing"
	usr_es "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

const (
	orgOwnerRole = "ORG_OWNER"
)

type OrgRepo struct {
	Eventstore       eventstore.Eventstore
	OrgEventstore    *org_es.OrgEventstore
	UserEventstore   *usr_es.UserEventstore
	PolicyEventstore *policy_es.PolicyEventstore

	View *admin_view.View

	SearchLimit    uint64
	SystemDefaults systemdefaults.SystemDefaults
}

func (repo *OrgRepo) SetUpOrg(ctx context.Context, setUp *admin_model.SetupOrg) (*admin_model.SetupOrg, error) {
	pwPolicy, err := repo.PolicyEventstore.GetPasswordComplexityPolicy(ctx, policy_model.DefaultPolicy)
	if err != nil {
		return nil, err
	}
	orgPolicy, err := repo.OrgEventstore.GetOrgIAMPolicy(ctx, policy_model.DefaultPolicy)
	if err != nil {
		return nil, err
	}
	users := func(ctx context.Context, domain string) ([]*es_models.Aggregate, error) {
		userIDs, err := repo.View.UserIDsByDomain(domain)
		if err != nil {
			return nil, err
		}
		return repo.UserEventstore.PrepareDomainClaimed(ctx, userIDs)
	}
	org, aggregates, err := repo.OrgEventstore.PrepareCreateOrg(ctx, setUp.Org, users)
	if err != nil {
		return nil, err
	}
	user, userAggregates, err := repo.UserEventstore.PrepareCreateUser(ctx, setUp.User, pwPolicy, orgPolicy, org.AggregateID)
	if err != nil {
		return nil, err
	}

	aggregates = append(aggregates, userAggregates...)
	setupModel := &Setup{Org: org, User: user}

	member := org_model.NewOrgMemberWithRoles(org.AggregateID, user.AggregateID, orgOwnerRole)
	_, memberAggregate, err := repo.OrgEventstore.PrepareAddOrgMember(ctx, member, org.AggregateID)
	if err != nil {
		return nil, err
	}
	aggregates = append(aggregates, memberAggregate)

	err = sdk.PushAggregates(ctx, repo.Eventstore.PushAggregates, setupModel.AppendEvents, aggregates...)
	if err != nil {
		return nil, err
	}

	return SetupToModel(setupModel), nil
}

func (repo *OrgRepo) OrgByID(ctx context.Context, id string) (*org_model.Org, error) {
	return repo.OrgEventstore.OrgByID(ctx, org_model.NewOrg(id))
}

func (repo *OrgRepo) SearchOrgs(ctx context.Context, query *org_model.OrgSearchRequest) (*org_model.OrgSearchResult, error) {
	query.EnsureLimit(repo.SearchLimit)
	sequence, err := repo.View.GetLatestOrgSequence()
	logging.Log("EVENT-LXo9w").OnError(err).Warn("could not read latest iam sequence")
	orgs, count, err := repo.View.SearchOrgs(query)
	if err != nil {
		return nil, err
	}
	result := &org_model.OrgSearchResult{
		Offset:      query.Offset,
		Limit:       query.Limit,
		TotalResult: uint64(count),
		Result:      model.OrgsToModel(orgs),
	}
	if err == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.CurrentTimestamp
	}
	return result, nil
}

func (repo *OrgRepo) IsOrgUnique(ctx context.Context, name, domain string) (isUnique bool, err error) {
	return repo.OrgEventstore.IsOrgUnique(ctx, name, domain)
}

func (repo *OrgRepo) GetOrgIamPolicyByID(ctx context.Context, id string) (*org_model.OrgIAMPolicy, error) {
	return repo.OrgEventstore.GetOrgIAMPolicy(ctx, id)
}

func (repo *OrgRepo) CreateOrgIamPolicy(ctx context.Context, policy *org_model.OrgIAMPolicy) (*org_model.OrgIAMPolicy, error) {
	return repo.OrgEventstore.AddOrgIAMPolicy(ctx, policy)
}

func (repo *OrgRepo) ChangeOrgIamPolicy(ctx context.Context, policy *org_model.OrgIAMPolicy) (*org_model.OrgIAMPolicy, error) {
	return repo.OrgEventstore.ChangeOrgIAMPolicy(ctx, policy)
}

func (repo *OrgRepo) RemoveOrgIamPolicy(ctx context.Context, id string) error {
	return repo.OrgEventstore.RemoveOrgIAMPolicy(ctx, id)
}

package api

import (
	"github.com/HewlettPackard/galadriel/pkg/common"
	"github.com/HewlettPackard/galadriel/pkg/common/entity"
	"github.com/google/uuid"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

func (td TrustDomain) ToEntity() (*entity.TrustDomain, error) {
	harvesterSpiffeID, err := spiffeid.FromString(*td.HarvesterSpiffeId)
	if err != nil {
		return nil, common.ErrWrongSPIFFEID{Cause: err}
	}

	tdName, err := spiffeid.TrustDomainFromString(td.Name)
	if err != nil {
		return nil, common.ErrWrongTrustDomain{Cause: err}
	}

	description := ""
	if td.Description != nil {
		description = *td.Description
	}

	onboardingBundle := []byte{}
	if td.OnboardingBundle != nil {
		onboardingBundle = []byte(*td.OnboardingBundle)
	}

	uuid := uuid.NullUUID{
		UUID:  td.Id,
		Valid: true,
	}

	return &entity.TrustDomain{
		ID:                uuid,
		Name:              tdName,
		CreatedAt:         td.CreatedAt,
		UpdatedAt:         td.UpdatedAt,
		Description:       description,
		OnboardingBundle:  onboardingBundle,
		HarvesterSpiffeID: harvesterSpiffeID,
	}, nil
}

func TrustDomainFromEntity(entity *entity.TrustDomain) *TrustDomain {
	onboardingBundle := string(entity.OnboardingBundle)
	harvesterSpiffeID := entity.HarvesterSpiffeID.String()

	return &TrustDomain{
		Id:                entity.ID.UUID,
		Name:              entity.Name.String(),
		Description:       &entity.Description,
		UpdatedAt:         entity.UpdatedAt,
		CreatedAt:         entity.CreatedAt,
		OnboardingBundle:  &onboardingBundle,
		HarvesterSpiffeId: &harvesterSpiffeID,
	}
}

func RelationshipFromEntity(entity *entity.Relationship) *Relationship {
	return &Relationship{
		Id:                  entity.ID.UUID,
		CreatedAt:           entity.CreatedAt,
		UpdatedAt:           entity.UpdatedAt,
		TrustDomainAId:      entity.TrustDomainAID,
		TrustDomainBId:      entity.TrustDomainBID,
		TrustDomainBConsent: entity.TrustDomainBConsent,
		TrustDomainAConsent: entity.TrustDomainAConsent,
	}
}

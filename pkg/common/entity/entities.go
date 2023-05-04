package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

// TODO: these entities will be defined by the OpenAPI specs and autogenerated.

type TrustDomain struct {
	ID                uuid.NullUUID
	Name              spiffeid.TrustDomain
	Description       string
	HarvesterSpiffeID spiffeid.ID
	OnboardingBundle  []byte
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Relationship struct {
	ID                  uuid.NullUUID
	TrustDomainAID      uuid.UUID
	TrustDomainBID      uuid.UUID
	TrustDomainAName    spiffeid.TrustDomain
	TrustDomainBName    spiffeid.TrustDomain
	TrustDomainAConsent bool
	TrustDomainBConsent bool
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type JoinToken struct {
	ID              uuid.NullUUID
	Token           string
	Used            bool
	TrustDomainID   uuid.UUID
	TrustDomainName spiffeid.TrustDomain
	ExpiresAt       time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Bundle represents a SPIFFE Trust bundle along with its digest.
type Bundle struct {
	ID                 uuid.NullUUID
	Data               []byte
	Signature          []byte
	SignatureAlgorithm string
	SigningCertificate []byte
	TrustDomainID      uuid.UUID
	TrustDomainName    spiffeid.TrustDomain
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

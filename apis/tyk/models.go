package tyk

// AllowedURL defines an allowed URL in an API.
type AllowedURL struct {
	URL     string   `json:"url"`
	Methods []string `json:"methods"`
}

// APIAccessRules is used while defining what an auth token has access to.
type APIAccessRules struct {
	// Best to leave it empty and inherit from policies.
	APIName string `json:"api_name"`

	// Must match key inside map where this object is value.
	APIID string `json:"api_id,omitempty"`

	// What versions of the API
	Versions []string `json:"versions,omitempty"`

	// Allowed URLs in this API. Best to not define it and inherit from policies.
	AllowedURLs []*AllowedURL `json:"allowed_urls,omitempty"`
}

// JWTData contains JWT secret.
type JWTData struct {
	Secret string `json:"secret"`
}

// Session is incomplete definition of Tyk session object, to cover our needs while managing Tyk sessions.
type Session struct {
	// Deprecated but expected. Needs to be same as Rate.
	Allowance int `json:"allowance"`

	// The number of requests that are allowed in the specified rate limiting window.
	Rate int `json:"rate"`
	// The number of seconds that the rate window should encompass.
	Per int `json:"per"`

	// An epoch that defines when the key should expire.
	Expires int64 `json:"expires"`

	// The maximum number of requests allowed during the quota period.
	QuotaMax int `json:"quota_max"`
	// An epoch that defines when the quota renews.
	QuotaRenews int64 `json:"quota_renews,omitempty"`
	// The number of requests remaining for this user’s quota (unrelated to rate limit).
	QuotaRemaining int64 `json:"quota_remaining"`
	// The time, in seconds. during which the quota is valid. So for 1000 requests per hour,
	// this value would be 3600 while quota_max and quota_remaining would be 1000.
	QuotaRenewalRate int `json:"quota_renewal_rate,omitempty"`

	// Defines what APIs and versions this token has access to. API IDs mapped to API access rules.
	AccessRights map[string]*APIAccessRules `json:"access_rights"`

	// ID of organization which this token belongs to.
	OrgID string `json:"org_id,omitempty"`

	// JWTData contains JWT secret.
	JWTData JWTData `json:"jwt_data"`

	// List of policy IDs to apply to this session.
	ApplyPolicies []string `json:"apply_policies,omitempty"`

	// Meta data to be included as part of the session, this is a key/value string map
	// that can be used in other middleware such as transforms and header injection to embed user-specific data
	// into a request, or alternatively to query the providence of a key.
	MetaData map[string]interface{} `json:"meta_data,omitempty"`

	// Tags are embedded into analytics data when the request completes. If a policy has tags,
	// those tags will supersede the ones carried by the token (they will be overwritten).
	Tags []string `json:"tags,omitempty"`

	// As of v2.1, an alias offers a way to identify a token in a more human-readable manner,
	// add an alias to a token in order to have the data transferred into Analytics later on
	// so you can track both hashed and un-hashed tokens to a meaningful identifier
	// that doesn’t expose the security of the underlying token.
	alias string `json:"tags,omitempty"`
}

// NewSession creates a new Tyk session object with our default settings.
func NewSession() *Session {
	s := &Session{
		Per:            60,   // Rate frame == every 60 seconds
		Rate:           1000, // How many requests per frame == 1000
		Expires:        -1,   // Never expires
		QuotaMax:       -1,   // Infinite
		QuotaRemaining: -1,   // Infinite
		AccessRights:   map[string]*APIAccessRules{},
	}
	s.Allowance = s.Rate // Allowance must always be equal to rate

	return s
}

// SetJWTSecret sets JWT secret.
func (session *Session) SetJWTSecret(secret string) *Session {
	session.JWTData.Secret = secret
	return session
}

// AddAccess adds access to API with versions.
func (session *Session) AddAccess(apiID string, versions ...string) *Session {
	ar := &APIAccessRules{
		APIID: apiID,
	}
	if versions != nil {
		ar.Versions = versions
	}
	if session.AccessRights == nil {
		session.AccessRights = map[string]*APIAccessRules{}
	}
	session.AccessRights[apiID] = ar

	return session
}

// WithPolicies adds policies to session.
func (session *Session) WithPolicies(policyIDs ...string) *Session {
	session.ApplyPolicies = append(session.ApplyPolicies, policyIDs...)
	return session
}

// SetAlias sets an alias.
func (session *Session) SetAlias(alias string) *Session {
	session.alias = alias
	return session
}

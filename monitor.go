package bettergoapi

import "time"

const bettergoapiURL = "https://bettergoapiuptime.com/api/v2/"

type Monitor struct {
	ID                  string     `json:"id,omitempty"`
	Type                string     `json:"type,omitempty"`
	URL                 string     `json:"url,omitempty"`
	PronounceableName   string     `json:"pronounceable_name,omitempty"`
	AuthUsername        string     `json:"auth_username,omitempty"`
	AuthPassword        string     `json:"auth_password,omitempty"`
	MonitorType         string     `json:"monitor_type,omitempty"`
	MonitorGroupID      any        `json:"monitor_group_id,omitempty"`
	LastCheckedAt       *time.Time `json:"last_checked_at,omitempty"`
	Status              string     `json:"status,omitempty"`
	PolicyID            any        `json:"policy_id,omitempty"`
	RequiredKeyword     string     `json:"required_keyword,omitempty"`
	VerifySsl           bool       `json:"verify_ssl,omitempty"`
	CheckFrequency      int        `json:"check_frequency,omitempty"`
	Call                bool       `json:"call,omitempty"`
	Sms                 bool       `json:"sms,omitempty"`
	Email               bool       `json:"email,omitempty"`
	Push                bool       `json:"push,omitempty"`
	TeamWait            any        `json:"team_wait,omitempty"`
	HTTPMethod          string     `json:"http_method,omitempty"`
	RequestTimeout      int        `json:"request_timeout,omitempty"`
	RecoveryPeriod      int        `json:"recovery_period,omitempty"`
	RequestHeaders      []any      `json:"request_headers,omitempty"`
	RequestBody         string     `json:"request_body,omitempty"`
	FollowRedirects     bool       `json:"follow_redirects,omitempty"`
	RememberCookies     bool       `json:"remember_cookies,omitempty"`
	CreatedAt           *time.Time `json:"created_at,omitempty"`
	UpdatedAt           *time.Time `json:"updated_at,omitempty"`
	SslExpiration       any        `json:"ssl_expiration,omitempty"`
	DomainExpiration    any        `json:"domain_expiration,omitempty"`
	Regions             []string   `json:"regions,omitempty"`
	ExpectedStatusCodes []int      `json:"expected_status_codes,omitempty"`
	Port                any        `json:"port,omitempty"`
	ConfirmationPeriod  int        `json:"confirmation_period,omitempty"`
	PausedAt            any        `json:"paused_at,omitempty"`
	Paused              *bool      `json:"paused,omitempty"`
	MaintenanceFrom     any        `json:"maintenance_from,omitempty"`
	MaintenanceTo       any        `json:"maintenance_to,omitempty"`
	MaintenanceTimezone string     `json:"maintenance_timezone,omitempty"`
}

func (b *Monitor) SetID(id string) error {
	b.ID = id
	return nil
}

func (b *Monitor) SetType(t string) error {
	b.Type = t
	return nil
}

func (b *Monitor) SetData(to func(target interface{}) error) error {
	return to(b)
}
func (b Monitor) GetID() string {
	return b.ID
}

func (b Monitor) GetType() string {
	return b.Type
}

func (b Monitor) GetData() interface{} {
	return b
}

type Monitors []Monitor

func (b *Monitors) SetData(to func(target interface{}) error) error {
	return to(b)
}

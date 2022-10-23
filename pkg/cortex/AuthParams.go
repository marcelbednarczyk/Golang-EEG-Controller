package cortex

type AuthParams struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	License      string `json:"license,omitempty"`
	Debit        int    `json:"debit,omitempty"`
}

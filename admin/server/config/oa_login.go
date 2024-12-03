package config

type OaLogin struct {
	Url                   string `mapstructure:"url" json:"url" yaml:"url"`
	Oauth2ClientId        string `mapstructure:"oauth2-client-id" json:"oauth2-client-id" yaml:"oauth2-client-id"`
	Oauth2ClientSecret    string `mapstructure:"oauth2-client-secret" json:"oauth2-client-secret" yaml:"oauth2-client-secret"`
	GetUserApiPath        string `mapstructure:"get-user-info-api-path" json:"get-user-info-api-path" yaml:"get-user-info-api-path"`
	GetTokenByCodeApiPath string `mapstructure:"get-token-by-code-api-path" json:"get-token-by-code-api-path" yaml:"get-token-by-code-api-path"`
}

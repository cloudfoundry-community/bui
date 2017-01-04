package uaa

type TokenResp struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	JTI          string `json:"jti"`
}

type Info struct {
	App struct {
		Version string `json:"version"`
	} `json:"app"`
	ZoneName string `json:"zone_name"`
	Links    struct {
		UAA      string `json:"uaa"`
		Password string `json:"passwd"`
		Login    string `json:"login"`
		Register string `json:"register"`
	}
}

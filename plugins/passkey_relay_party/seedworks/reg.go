package seedwork

type Registration struct {
	Origin      string `json:"origin"`
	Account     string `json:"account"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
}

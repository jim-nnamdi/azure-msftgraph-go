package azure

type AzUser struct {
	OdataContext    string               `json:"odatacontext"`
	Id              string               `json:"id"`
	DisplayName     string               `json:"displayName"`
	AccountEnabled  bool                 `json:"accountEnabled"`
	GivenName       string               `json:"givenName"`
	Mail            string               `json:"mail"`
	TelephoneNumber string               `json:"mobilePhone"`
	Surname         string               `json:"surname"`
	CreatedDateTime string               `json:"createdDateTime"`
	Success         bool                 `json:"success"`
	Name            string               `json:"name"`
	Identities      []IdData             `json:"identities"`
	Authed          bool                 `json:"authed"`
	AzAuth          AzRetrieveAuthedData `json:"azauth"`
}

type AzRetrieveAuthedData struct {
	Success        bool   `json:"success"`
	Surname        string `json:"surname"`
	Telephone      string `json:"telephone"`
	CreatedAt      string `json:"createdDateTime"`
	GivenName      string `json:"given_name"`
	AccountEnabled bool   `json:"account_enabled"`
	Email          string `json:"email"`
	OldAuthed      bool   `json:"old_authed"`
}

type IdData struct {
	SignInType       string `json:"signInType"`
	Issuer           string `json:"issuer"`
	IssuerAssignedId string `json:"issuerAssignedId"`
}

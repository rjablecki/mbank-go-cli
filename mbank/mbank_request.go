package mbank

type LoginRequest struct {
	HrefHasHash        bool
	Scenario           string
	DfpData            map[string]string
	UWAdditionalParams struct {
		InOut         string
		ReturnAddress string
		Source        string
	}
	UserName int
	Password string
	Seed     string
	Lang     string
}

package settings

// Branding contains the branding settings of the app.
type Branding struct {
	Name            string `json:"name"`
	DisableExternal bool   `json:"disableExternal"`
	EmbededMode     bool   `json:"embededMode"`
	Files           string `json:"files"`
}

package validators

type Validator struct {
	Attrs   map[string]string
	Message string
	Error   string
	User    bool
}

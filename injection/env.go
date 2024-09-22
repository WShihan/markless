package injection

type Env struct {
	BaseURL     string
	Title       string
	DataBaseURL string
	Port        int
	Version     string
	Commit      string
	BuildTime   string
	HmacSecret  string
	SecretKey   string
	JWTExpire   int
}

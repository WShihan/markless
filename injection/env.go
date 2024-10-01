package injection

type Env struct {
	BaseURL     string `json:"base_url"`
	Title       string `json:"title"`
	DataBaseURL string `json:"database_url"`
	Port        int    `json:"port"`
	Version     string `json:"version"`
	Commit      string `json:"commit"`
	BuildTime   string `json:"build_time"`
	HmacSecret  string `json:"hmacsecret"`
	SecretKey   string `json:"secret_key"`
	JWTExpire   int    `json:"jwt_expires"`
}

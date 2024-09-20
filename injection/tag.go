package injection

type TagsPage struct {
	Page PageInjection
	Env  Env
	Data interface{}
}

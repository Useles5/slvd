package platform

type Provider interface {
	FetchRecent(handle string) ([]string, error)
}

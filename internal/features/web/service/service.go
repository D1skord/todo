package web_service

type WebSerice struct {
	webRepository WebRepository
}

type WebRepository interface {
	GetFile(filePath string) ([]byte, error)
}

func NewWebService(
	webRepository WebRepository,
) *WebSerice {
	return &WebSerice{
		webRepository: webRepository,
	}
}

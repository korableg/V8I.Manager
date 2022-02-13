package handlers

type (
	Handlers struct {
		service user.Service
	}
)

func NewHandlers(s user.Service) *Handlers {
	h := &Handlers{
		service: s,
	}

	return h
}

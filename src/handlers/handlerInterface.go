package handlers

type HandlerInterface interface {
	ProxyRequest() error
	Defer()
}

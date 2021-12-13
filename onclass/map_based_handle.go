package onclass

import "net/http"

type HandlerBasedMap struct {
	handlers map[string]func(ctx *Context)
}

func (h *HandlerBasedMap) ServeHTTP(writer http.ResponseWriter,
	request *http.Request) {
	key := h.Key(request.Method, request.URL.Path)
	if handler, ok := h.handlers[key]; ok {
		handler(NewContext(writer, request))
	} else {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte("Not Found"))
	}
}

func (h *HandlerBasedMap) Key(method string, pattern string) string {
	return method + "#" + pattern
}

package web_transport_http

import (
	"net/http"

	core_logger "github.com/D1skord/todo/internal/core/logger"
	core_http_response "github.com/D1skord/todo/internal/core/transport/http/response"
)

func (h *WebHTTPHandler) GetMainPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	html, err := h.webSerice.GetMainPage()
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get index.html for main page",
		)

		return
	}

	responseHandler.HTMLResponse(html)
}

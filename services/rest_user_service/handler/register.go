package handler

import (
	errm "authentication_service/core/errmodule"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

// ErrorResponse структура для возврата ошибок
type ErrorResponse struct {
	ErrorCode        int    `json:"error_code"`        // Код ошибки
	ErrorDescription string `json:"error_description"` // Описание ошибки
}

// WrapHandlerParams структура для параметров обертки
type WrapHandlerParams struct {
	HandlerFunc func(http.ResponseWriter, *http.Request) (interface{}, *errm.Error)
}

// responseWriterWrapper предотвращает повторный вызов WriteHeader
type responseWriterWrapper struct {
	http.ResponseWriter
	wroteHeader bool
}

// WriteHeader перехватывает вызов WriteHeader
func (rw *responseWriterWrapper) WriteHeader(statusCode int) {
	if !rw.wroteHeader {
		rw.ResponseWriter.WriteHeader(statusCode)
		rw.wroteHeader = true
	}
}

// Write перехватывает вызов Write
func (rw *responseWriterWrapper) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK) // Устанавливаем статус по умолчанию
	}
	return rw.ResponseWriter.Write(b)
}

// WrapHandlerF обертка для обработчиков
func WrapHandlerF(p WrapHandlerParams) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriterWrapper{ResponseWriter: w}
		payload, err := p.HandlerFunc(rw, r)
		respondWithJSON(rw, err, payload)
	}
}

// RegisterRoute вспомогательная функция для регистрации маршрутов
func RegisterRoute(r chi.Router, method, path string, handlerFunc func(http.ResponseWriter, *http.Request) (interface{}, *errm.Error)) {
	wrappedHandler := WrapHandlerF(WrapHandlerParams{
		HandlerFunc: handlerFunc,
	})
	switch method {
	case http.MethodGet:
		r.Get(path, wrappedHandler)
	case http.MethodPost:
		r.Post(path, wrappedHandler)
	case http.MethodPatch:
		r.Patch(path, wrappedHandler)
	case http.MethodPut:
		r.Put(path, wrappedHandler)
	case http.MethodDelete:
		r.Delete(path, wrappedHandler)
	}
}

// ParseRequestBodyPost обрабатывает тело запроса и заполняет структуру запроса
func ParseRequestBodyPost(r *http.Request, v interface{}) *errm.Error {
	defer r.Body.Close()

	if r.ContentLength == 0 {
		return errm.NewError("empty_request_body", errors.New("empty request body"))
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return errm.NewError("read_request_body", err)
	}

	if err := json.Unmarshal(body, v); err != nil {
		return errm.NewError("unmarshal_request_body", err)
	}

	return nil
}

// respondWithJSON отправляет ответ в формате JSON
func respondWithJSON(w http.ResponseWriter, errObj *errm.Error, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	var status int
	var response []byte
	var err error

	if errObj != nil {
		status = http.StatusInternalServerError
		errorResponse := ErrorResponse{
			ErrorCode:        errObj.Code,
			ErrorDescription: errObj.Description,
		}
		response, _ = json.Marshal(errorResponse)
	} else {
		response, err = json.Marshal(payload)
		if err != nil {
			status = http.StatusInternalServerError
			errorResponse := ErrorResponse{
				ErrorCode:        http.StatusInternalServerError,
				ErrorDescription: err.Error(),
			}
			response, _ = json.Marshal(errorResponse)
		} else {
			status = http.StatusOK
		}
	}

	// Проверяем, является ли w нашим responseWriterWrapper
	if rw, ok := w.(*responseWriterWrapper); ok {
		if !rw.wroteHeader {
			rw.WriteHeader(status)
		}
	} else {
		w.WriteHeader(status)
	}

	_, _ = w.Write(response)
}

package caphttp

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/lin-snow/ech0/pkg/gocap/core"
)

type errorBody struct {
	Success bool   `json:"success"`
	Code    string `json:"code,omitempty"`
	Error   string `json:"error"`
}

type decodeOptions struct {
	Strict bool
}

func decodeJSON(r *http.Request, out any, opts decodeOptions) error {
	dec := json.NewDecoder(r.Body)
	if opts.Strict {
		dec.DisallowUnknownFields()
	}
	if err := dec.Decode(out); err != nil {
		return err
	}
	// Reject multiple JSON values.
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return core.NewBadRequest("Malformed JSON body")
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeCoreError(w http.ResponseWriter, err error) {
	var ce *core.Error
	if errors.As(err, &ce) {
		writeJSON(w, ce.StatusCode, errorBody{
			Success: false,
			Code:    string(ce.Code),
			Error:   ce.Message,
		})
		return
	}
	writeJSON(w, http.StatusInternalServerError, errorBody{
		Success: false,
		Code:    string(core.ErrCodeInternal),
		Error:   "Internal server error",
	})
}

func writeDecodeError(w http.ResponseWriter, err error) {
	var ce *core.Error
	if errors.As(err, &ce) {
		writeCoreError(w, ce)
		return
	}
	var mbe *http.MaxBytesError
	if errors.As(err, &mbe) {
		writeCoreError(w, core.NewBadRequest("Request body too large"))
		return
	}
	writeCoreError(w, core.NewBadRequest("Malformed JSON body"))
}

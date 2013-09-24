package libLiveSrv

import (
	"fmt"
	"io"
	"net/http"
)

// HTTPResolver resolves agora packages over HTTP
type HTTPResolver struct {
	baseURL string
}

// Resolve a package over HTTP instead of on the local filesystem
func (h *HTTPResolver) Resolve(modPath string) (io.Reader, error) {
	if h.baseURL == "" {
		h.baseURL = "http://localhost:8000"
	}
	resp, err := http.Get(h.baseURL + "/" + modPath)
	if err != nil {
		fmt.Println("Error in HTTPResolver getting module", err.Error())
		return nil, err
	}
	return resp.Body, nil
}

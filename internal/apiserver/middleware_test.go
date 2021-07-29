package apiserver

import (
	"context"
	"net/http"
	"testing"

	"github.com/alisavch/image-service/internal/service"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestHandler_getUserID(t *testing.T) {
	var getContext = func(id int) *http.Request {
		ctx := context.WithValue(context.Background(), userCtx, id)
		request := http.Request{}
		return request.WithContext(ctx)
	}

	tests := []struct {
		name string
		ctx  *http.Request
		id   int
		isOk bool
	}{
		{
			name: "Test with correct values",
			ctx:  getContext(1),
			id:   1,
			isOk: true,
		},
		{
			name: "Test with incorrect values",
			ctx:  &http.Request{},
			isOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			services := &service.Service{}
			s := Server{router: mux.NewRouter(), service: services}

			id, err := s.getUserID(tt.ctx)
			if tt.isOk {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
			require.Equal(t, id, tt.id)
		})
	}
}

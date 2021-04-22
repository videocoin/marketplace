package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/videocoin/marketplace/api/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
	"net/http"
	"net/textproto"
)

type errorBody struct {
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields"`
}

func WithProtoHTTPErrorHandler() runtime.ServeMuxOption {
	return runtime.WithProtoErrorHandler(DefaultHandleHTTPError)
}

func DefaultHandleHTTPError(
	ctx context.Context,
	mux *runtime.ServeMux,
	marshaler runtime.Marshaler,
	w http.ResponseWriter,
	_ *http.Request,
	err error,
) {
	w.Header().Del("Trailer")
	w.Header().Set("Content-Type", marshaler.ContentType())

	s, ok := status.FromError(err)
	if !ok {
		s = status.New(codes.Unknown, err.Error())
	}

	fieldsErr := map[string]string{}
	detailErrs := s.Details()
	if detailErrs != nil {
		for _, detailErr := range detailErrs {
			verrs, ok := detailErr.(*rpc.MultiValidationError)
			if ok {
				for _, fieldErr := range verrs.Errors {
					fieldsErr[fieldErr.Field] = fieldErr.Message
				}
			}
		}
	}

	if len(fieldsErr) == 0 {
		fieldsErr = nil
	}

	body := &errorBody{
		Message: s.Message(),
		Fields:  fieldsErr,
	}

	buf, merr := marshaler.Marshal(body)
	if merr != nil {
		grpclog.Printf("failed to marshal error message %q: %v", body, merr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		grpclog.Printf("failed to extract ServerMetadata from context")
	}

	handleForwardResponseServerMetadata(w, mux, md)
	handleForwardResponseTrailerHeader(w, md)

	st := runtime.HTTPStatusFromCode(s.Code())
	w.WriteHeader(st)

	if st >= 500 {
		grpclog.Errorf("failed: %s", string(buf))
		jsonErrMsg, _ := json.Marshal(map[string]string{"message": "Oops. Something went wrong! Sorry. We've let our engineers know."})
		if st == 502 || st == 503 {
			jsonErrMsg, _ = json.Marshal(map[string]string{"message": "Service Unavailable"})
		}
		if st == 501 {
			jsonErrMsg, _ = json.Marshal(map[string]string{"message": s.Message()})
		}
		if _, err := w.Write(jsonErrMsg); err != nil {
			grpclog.Printf("failed to write response: %v", err)
		}
	} else {
		if _, err := w.Write(buf); err != nil {
			grpclog.Printf("failed to write response: %v", err)
		}
	}

	handleForwardResponseTrailer(w, md)
}

func outgoingHeaderMatcher(key string) (string, bool) {
	return fmt.Sprintf("%s%s", runtime.MetadataHeaderPrefix, key), true
}

func handleForwardResponseServerMetadata(w http.ResponseWriter, mux *runtime.ServeMux, md runtime.ServerMetadata) {
	for k, vs := range md.HeaderMD {
		if h, ok := outgoingHeaderMatcher(k); ok {
			for _, v := range vs {
				w.Header().Add(h, v)
			}
		}
	}
}

func handleForwardResponseTrailerHeader(w http.ResponseWriter, md runtime.ServerMetadata) {
	for k := range md.TrailerMD {
		tKey := textproto.CanonicalMIMEHeaderKey(fmt.Sprintf("%s%s", runtime.MetadataTrailerPrefix, k))
		w.Header().Add("Trailer", tKey)
	}
}

func handleForwardResponseTrailer(w http.ResponseWriter, md runtime.ServerMetadata) {
	for k, vs := range md.TrailerMD {
		tKey := fmt.Sprintf("%s%s", runtime.MetadataTrailerPrefix, k)
		for _, v := range vs {
			w.Header().Add(tKey, v)
		}
	}
}

func IsNotFoundError(err error) bool {
	if s, ok := status.FromError(err); ok {
		if s.Code() == codes.NotFound {
			return true
		}
	}

	return false
}

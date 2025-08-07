package grpc

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	pb "go-backend-skeleton/app/internal/transport/grpc/pb"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CustomHTTPErrorResponse defines the structure for your HTTP error JSON.
type CustomHTTPErrorResponse struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Details *CustomErrorDetails `json:"details,omitempty"`
}

// CustomErrorDetails is a simplified structure to hold your custom error data
// without the @type field.
type CustomErrorDetails struct {
	ValidationErrors []ValidationErrorDetail `json:"validationErrors,omitempty"`
	NotFoundDetails  *NotFoundErrorDetail    `json:"notFoundDetails,omitempty"`
}

type ValidationErrorDetail struct {
	Reason  string `json:"reason"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

type NotFoundErrorDetail struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// CustomHTTPErrorHandler is a custom error handler for grpc-gateway.
func CustomHTTPErrorHandler(
	ctx context.Context,
	mux *runtime.ServeMux,
	marshaler runtime.Marshaler,
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	w.Header().Set("Content-Type", "application/json")

	statusCode := runtime.HTTPStatusFromCode(codes.Internal)
	var customDetails *CustomErrorDetails

	s, ok := status.FromError(err)
	if ok {
		statusCode = runtime.HTTPStatusFromCode(s.Code())
		for _, detail := range s.Details() {
			if pbErr, isPbError := detail.(*pb.Error); isPbError {
				customDetails = &CustomErrorDetails{}
				if len(pbErr.ValidationErrorDetails) > 0 {
					for _, vd := range pbErr.ValidationErrorDetails {
						customDetails.ValidationErrors = append(customDetails.ValidationErrors,
							ValidationErrorDetail{
								Reason:  vd.GetReason(),
								Field:   vd.GetField(),
								Message: vd.GetMessage(),
							},
						)
					}

					sendValidationError(w, statusCode, customDetails.ValidationErrors)
					return
				}
				if pbErr.NotFoundErrorDetails != nil {
					sendNotFoundError(w, statusCode, pbErr.NotFoundErrorDetails)
					return
				}
			}
		}
	}

	sendInternalServerError(w, statusCode)
}

func sendValidationError(
	w http.ResponseWriter,
	statusCode int,
	validationErrorDetailsSlice []ValidationErrorDetail,
) {
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(CustomHTTPErrorResponse{
		Code:    statusCode,
		Message: "invalid_values",
		Details: &CustomErrorDetails{
			ValidationErrors: validationErrorDetailsSlice,
		},
	})
}

func sendNotFoundError(
	w http.ResponseWriter,
	statusCode int,
	err *pb.NotFoundErrorDetails,
) {
	w.WriteHeader(http.StatusNotFound)
	underscoredType := strings.ReplaceAll(string(err.Type), " ", "_")
	_ = json.NewEncoder(w).Encode(CustomHTTPErrorResponse{
		Code:    statusCode,
		Message: underscoredType + " " + err.Id + " not found",
		Details: &CustomErrorDetails{
			NotFoundDetails: &NotFoundErrorDetail{
				ID:   err.Id,
				Type: err.Type,
			},
		},
	})
}

func sendInternalServerError(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(CustomHTTPErrorResponse{
		Code:    statusCode,
		Message: "internal server error",
	})
}

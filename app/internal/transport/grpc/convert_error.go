package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ConvertError2Pb(err error, method string) error {
	// var nfErr notfound.Error
	// if errors.As(err, &nfErr) {
	// 	st := status.New(codes.NotFound, "grpc."+method+": not found")
	//
	// 	pbNotFoundError := &pb.Error{
	// 		NotFoundErrorDetails: &pb.NotFoundErrorDetails{
	// 			Id:   nfErr.ID,
	// 			Type: string(nfErr.Type),
	// 		},
	// 	}
	//
	// 	st, err = st.WithDetails(pbNotFoundError)
	// 	if err != nil {
	// 		return st.Err()
	// 	}
	//
	// 	return st.Err()
	// }
	// if errors.Is(err, validate.Error{}) {
	// 	st := status.New(codes.NotFound, "grpc."+method+": validation error")
	//
	// 	var multiErr *multierror.Error
	// 	var errs []error
	// 	if errors.As(err, &multiErr) {
	// 		errs = multiErr.Errors
	// 	} else {
	// 		errs = []error{err}
	// 	}
	//
	// 	details := make([]*pb.ValidationErrorDetails, len(errs))
	// 	for i, e := range errs {
	// 		var validationErr validate.Error
	// 		if errors.As(e, &validationErr) {
	// 			details[i] = &pb.ValidationErrorDetails{
	// 				Field:   validationErr.Field,
	// 				Reason:  string(validationErr.Reason),
	// 				Message: validationErr.Message,
	// 			}
	// 		}
	// 	}
	//
	// 	pbValidationError := &pb.Error{
	// 		ValidationErrorDetails: details,
	// 	}
	//
	// 	st, err = st.WithDetails(pbValidationError)
	// 	if err != nil {
	// 		return st.Err()
	// 	}
	//
	// 	return st.Err()
	// }

	// Test the pb error type
	// testPbErr := &pb.Error{
	// 	ValidationErrorDetails: []*pb.ValidationErrorDetails{
	// 		{
	// 			Reason:  "reason",
	// 			Field:   "field",
	// 			Message: "message",
	// 		},
	// 	},
	// 	NotFoundErrorDetails: &pb.NotFoundErrorDetails{
	// 		Id:   "id",
	// 		Type: "type",
	// 	},
	// }

	st := status.New(codes.Internal, "grpc."+method+": error: "+err.Error())
	// st, err = st.WithDetails(testPbErr)
	if err != nil {
		return st.Err()
	}

	return st.Err()
}

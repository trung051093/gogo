package decorator

import (
	"context"

	"go.opencensus.io/trace"
)

// Trace with function of struct
// Example: TraceService(ctx, "test")(&userService, "SearchUsers")(agrs...)
func TraceService(ctx context.Context, name string) func(interface{}, string) func(...interface{}) (interface{}, error) {
	_, span := trace.StartSpan(ctx, name)
	return func(service interface{}, method string) func(...interface{}) (interface{}, error) {
		return func(agrs ...interface{}) (interface{}, error) {
			result, err := Invoke(service, method, agrs...)
			defer span.End()
			return result.Interface(), err
		}
	}
}

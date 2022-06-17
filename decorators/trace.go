package decorator

import (
	"context"

	"go.opencensus.io/trace"
)

// Trace service
// Example: TraceService(ctx, "test")(&userService, "SearchUsers")(agrs...)
func TraceService(ctx context.Context, name string) func(interface{}, string) func(...interface{}) (interface{}, error) {
	return func(service interface{}, method string) func(...interface{}) (interface{}, error) {
		return func(agrs ...interface{}) (interface{}, error) {
			_, span := trace.StartSpan(ctx, name)
			result, err := Invoke(service, method, agrs...)
			span.End()
			return result.Interface(), err
		}
	}
}

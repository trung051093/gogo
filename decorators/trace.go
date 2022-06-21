package decorator

import (
	"context"

	"go.opencensus.io/trace"
)

// Trace service
// Example: TraceService[Result type](ctx, "test")(&userService, "SearchUsers")(agrs...)
func TraceService[R any](ctx context.Context, traceName string) func(interface{}, string) func(...interface{}) (R, error) {
	return func(service interface{}, method string) func(...interface{}) (R, error) {
		return func(agrs ...interface{}) (R, error) {
			_, span := trace.StartSpan(ctx, traceName)
			result, err := Invoke(service, method, agrs...)
			span.End()
			return result.Interface().(R), err
		}
	}
}

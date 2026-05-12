package handlers

// MetricsHandler
// @Title Metrics
// @Description Returns a predefined set of metrics to be consumed by a monitoring system.
// @Success 200 text/plain
// @Route /metrics [get]
func MetricsHandler() {
	// This function is just to be able to add the OpenAPI documentation for the /metrics endpoint.
	// The actual handler is provided by the Prometheus client library.
}

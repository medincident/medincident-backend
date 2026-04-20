package config

// NATSConfig is the JetStream connection block used by query-server's
// identity projector. Stream and subjects are operator-provisioned;
// query-server only attaches a durable consumer to an already-existing
// stream.
type NATSConfig struct {
	URL         string   `yaml:"url"          validate:"required"`
	Stream      string   `yaml:"stream"       validate:"required"`
	Subjects    []string `yaml:"subjects"     validate:"required,min=1,dive,required"`
	DurableName string   `yaml:"durable_name" validate:"required"`
}

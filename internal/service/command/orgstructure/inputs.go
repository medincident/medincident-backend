package orgstructure

// AddressInput is the command-layer input for an address. Text is
// required and bounded; Point is optional and, when present, both
// coordinates must be inside their valid ranges.
type AddressInput struct {
	Text  string      `validate:"required,min=4,max=128"`
	Point *PointInput `validate:"omitnil"`
}

// PointInput is the command-layer input for optional geographic
// coordinates.
type PointInput struct {
	Longitude float64 `validate:"min=-180,max=180"`
	Latitude  float64 `validate:"min=-90,max=90"`
}

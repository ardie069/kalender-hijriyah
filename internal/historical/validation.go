package historical

// Validator validates historical claims against retrocalculated astronomy
type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

// IsAstronomicallyPossible checks if a sighting claim was physically possible
func (v *Validator) IsAstronomicallyPossible(record RukyatRecord) bool {
	// Placeholder: In a real implementation, this would call astronomy logic
	// to verify if elongation/altitude allowed for naked-eye visibility
	return true
}

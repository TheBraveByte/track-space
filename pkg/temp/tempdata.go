package temp

// probably wont be making use of this package
type TemplateData struct {
	StringData      map[string]string
	IntData         map[string]int
	Float64Data     map[string]float64
	Error           string
	RandomData      map[string]interface{}
	IsAuthenticated int
}

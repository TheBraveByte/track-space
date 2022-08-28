package temp

// probably wont be making use of this package
type TemplateData struct {
	AuthData        map[string][]string
	Count           map[string]int
	Float64Data     map[string]float64
	Error           string
	RandomData      map[string]interface{}
	IsAuthenticated int
}

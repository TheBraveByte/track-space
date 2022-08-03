package temp

type TemplateData struct {
	StringData   map[string]string
	IntData      map[string]int
	Float64Data  map[string]float64
	Error        string
	FlashMessage string
	RandomData   map[string]interface{}
}

package criteria

// Registry holds a collection of visibility evaluators
type Registry struct {
	evaluators map[string]Evaluator
}

func NewRegistry() *Registry {
	return &Registry{
		evaluators: make(map[string]Evaluator),
	}
}

func (r *Registry) Register(name string, e Evaluator) {
	r.evaluators[name] = e
}

func (r *Registry) Get(name string) (Evaluator, bool) {
	e, exists := r.evaluators[name]
	return e, exists
}
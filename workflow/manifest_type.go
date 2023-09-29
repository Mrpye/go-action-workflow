package workflow

// Manifest is a workflow manifest and contains all the information needed to run a workflow
type Manifest struct {
	Meta       MetaData               `json:"meta_data" yaml:"meta_data"`
	Actions    []Action               `json:"actions" yaml:"actions"`
	Jobs       []Job                  `json:"jobs" yaml:"jobs"`
	Parameters []Parameter            `json:"parameters" yaml:"parameters"`
	Data       map[string]interface{} `json:"data" yaml:"data"`
}

// ManifestOption is a function that sets an option on the manifest
type ManifestOption func(*Manifest)

// OptionManifestMeta sets the meta data on the manifest
// v is the meta data
func OptionManifestMeta(v MetaData) ManifestOption {
	return func(h *Manifest) {
		h.Meta = v
	}
}

// OptionManifestJobs sets the jobs on the manifest
// v is the slice of jobs
func OptionManifestJobs(v []Job) ManifestOption {
	return func(h *Manifest) {
		h.Jobs = v
	}
}

// OptionManifestParameters sets the parameters on the manifest
// v is the slice of parameters
func OptionManifestParameters(v []Parameter) ManifestOption {
	return func(h *Manifest) {
		h.Parameters = v
	}
}

// OptionManifestData sets the data on the manifest
// v is the data as a map[string]interface{}
func OptionManifestData(v map[string]interface{}) ManifestOption {
	return func(h *Manifest) {
		h.Data = v
	}
}

// CreateManifest creates a new manifest
// opts are the options to set on the manifest
// returns the manifest
func CreateManifest(opts ...ManifestOption) *Manifest {
	manifest := &Manifest{}
	for _, opt := range opts {
		opt(manifest)
	}
	return manifest
}

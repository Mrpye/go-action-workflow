package workflow

type Manifest struct {
	Meta       MetaData               `json:"meta_data" yaml:"meta_data"`
	Jobs       []Job                  `json:"jobs" yaml:"jobs"`
	Parameters []Parameter            `json:"parameters" yaml:"parameters"`
	Data       map[string]interface{} `json:"data" yaml:"data"`
}

type ManifestOption func(*Manifest)

func OptionManifestMeta(v MetaData) ManifestOption {
	return func(h *Manifest) {
		h.Meta = v
	}
}

func OptionManifestJobs(v []Job) ManifestOption {
	return func(h *Manifest) {
		h.Jobs = v
	}
}
func OptionManifestParameters(v []Parameter) ManifestOption {
	return func(h *Manifest) {
		h.Parameters = v
	}
}
func CreateManifest(opts ...ManifestOption) *Manifest {
	manifest := &Manifest{}
	for _, opt := range opts {
		opt(manifest)
	}
	return manifest
}

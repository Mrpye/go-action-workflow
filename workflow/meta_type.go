package workflow

// MetaData is the metadata for a workflow
type MetaData struct {
	Name        string                 `json:"name" yaml:"name" flag:"name n" desc:"Project Name"`
	Description string                 `json:"description" yaml:"description" flag:"description d" desc:"Project Description"`
	Version     string                 `json:"version" yaml:"version" flag:"version v" desc:"Project Version"`
	Author      string                 `json:"author" yaml:"author" flag:"author a" desc:"Project Author"`
	Contact     string                 `json:"contact" yaml:"contact" flag:"contact c" desc:"Contact Details"`
	CreatedDate string                 `json:"create_date" yaml:"create_date" flag:"create_date" desc:"Created Date"`
	UpdateDate  string                 `json:"update_date" yaml:"update_date" flag:"update_date" desc:"Updated Date"`
	Vars        map[string]interface{} `json:"vars" yaml:"vars"`
}

// MetaDataOption is a function that sets an option on the metadata
type MetaDataOption func(*MetaData)

// OptionMetaDataName sets the name on the metadata
// v is the name
func OptionMetaDataName(v string) MetaDataOption {
	return func(h *MetaData) {
		h.Name = v
	}
}

// OptionMetaDataDescription sets the description on the metadata
// v is the description
func OptionMetaDataDescription(v string) MetaDataOption {
	return func(h *MetaData) {
		h.Description = v
	}
}

// OptionMetaDataVersion sets the version on the metadata
// v is the version
func OptionMetaDataVersion(v string) MetaDataOption {
	return func(h *MetaData) {
		h.Version = v
	}
}

// OptionMetaDataAuthor sets the author on the metadata
// v is the author
func OptionMetaDataAuthor(v string) MetaDataOption {
	return func(h *MetaData) {
		h.Author = v
	}
}

// OptionMetaDataContact sets the contact on the metadata
// v is the contact
func OptionMetaDataContact(v string) MetaDataOption {
	return func(h *MetaData) {
		h.Contact = v
	}
}

// OptionMetaDataCreatedDate sets the created date on the metadata
// v is the created date
func OptionMetaDataCreatedDate(v string) MetaDataOption {
	return func(h *MetaData) {
		h.CreatedDate = v
	}
}

// OptionMetaDataUpdateDate sets the update date on the metadata
// v is the update date
func OptionMetaDataUpdateDate(v string) MetaDataOption {
	return func(h *MetaData) {
		h.UpdateDate = v
	}
}

// OptionMetaDataVars sets the vars on the metadata
// v is the vars
func OptionMetaDataVars(v map[string]interface{}) MetaDataOption {
	return func(h *MetaData) {
		h.Vars = v
	}
}

// CreateMetaData creates a new metadata
// opts are the options to set on the metadata
// returns the metadata
func CreateMetaData(opts ...MetaDataOption) *MetaData {
	meta := &MetaData{}
	for _, opt := range opts {
		opt(meta)
	}
	return meta
}

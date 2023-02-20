package workflow

type MetaData struct {
	Name        string                 `json:"name" yaml:"name"`
	Description string                 `json:"description" yaml:"description"`
	Version     string                 `json:"version" yaml:"version"`
	Author      string                 `json:"author" yaml:"author"`
	Contact     string                 `json:"contact" yaml:"contact"`
	CreatedDate string                 `json:"create_date" yaml:"create_date"`
	UpdateDate  string                 `json:"update_date" yaml:"update_date"`
	Vars        map[string]interface{} `json:"vars" yaml:"vars"`
}

type MetaDataOption func(*MetaData)

func OptionMetaDataName(v string) MetaDataOption {
	return func(h *MetaData) {
		h.Name = v
	}
}
func OptionMetaDataDescription(v string) MetaDataOption {
	return func(h *MetaData) {
		h.Description = v
	}
}
func OptionMetaDataVersion(v string) MetaDataOption {
	return func(h *MetaData) {
		h.Version = v
	}
}
func OptionMetaDataAuthor(v string) MetaDataOption {
	return func(h *MetaData) {
		h.Author = v
	}
}
func OptionMetaDataContact(v string) MetaDataOption {
	return func(h *MetaData) {
		h.Contact = v
	}
}
func OptionMetaDataCreatedDate(v string) MetaDataOption {
	return func(h *MetaData) {
		h.CreatedDate = v
	}
}
func OptionMetaDataUpdateDate(v string) MetaDataOption {
	return func(h *MetaData) {
		h.UpdateDate = v
	}
}
func OptionMetaDataVars(v map[string]interface{}) MetaDataOption {
	return func(h *MetaData) {
		h.Vars = v
	}
}

func CreateMetaData(opts ...MetaDataOption) *MetaData {
	meta := &MetaData{}
	for _, opt := range opts {
		opt(meta)
	}
	return meta
}

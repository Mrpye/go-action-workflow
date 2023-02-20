package workflow

// dataBucket is a map of maps that can be used to store data between actions
// The first key is the name of the bucket
// The second key is the name of the value
// The value is the data
func (m *Workflow) SetValueToDataBucket(key string, name string, value interface{}) {
	if m.dataBucket == nil {
		m.dataBucket = make(map[string]map[string]interface{})
	}
	if m.dataBucket[key] == nil {
		m.dataBucket[key] = make(map[string]interface{})
	}
	m.dataBucket[key][name] = value
}

// GetValueFromDataBucket returns the value from the data bucket
// The first key is the name of the bucket
// The second key is the name of the value
// The value is the data or nil if the value does not exist
func (m *Workflow) GetValueFromDataBucket(key string, name string) interface{} {
	if m.dataBucket == nil {
		m.dataBucket = make(map[string]map[string]interface{})
	}
	if m.dataBucket[key] == nil {
		m.dataBucket[key] = make(map[string]interface{})
	}
	if val, ok := m.dataBucket[key][name]; ok {
		return val
	}
	return nil
}

// GetValueFromDataBucketAsStrings returns the value from the data bucket as a string
// The first key is the name of the bucket
// The second key is the name of the value
// The value is the data or nil if the value does not exist
func (m *Workflow) GetValueFromDataBucketAsStrings(key string, name string) string {
	if m.dataBucket == nil {
		m.dataBucket = make(map[string]map[string]interface{})
	}
	if m.dataBucket[key] == nil {
		m.dataBucket[key] = make(map[string]interface{})
	}
	if val, ok := m.dataBucket[key][name]; ok {
		return val.(string)
	}
	return ""
}

// GetValueFromDataBucketAsInt returns the value from the data bucket as an int
// The first key is the name of the bucket
// The second key is the name of the value
// The value is the data or nil if the value does not exist
func (m *Workflow) GetValueFromDataBucketAsInt(key string, name string) int {
	if m.dataBucket == nil {
		m.dataBucket = make(map[string]map[string]interface{})
	}
	if m.dataBucket[key] == nil {
		m.dataBucket[key] = make(map[string]interface{})
	}
	if val, ok := m.dataBucket[key][name]; ok {
		return val.(int)
	}
	return 0
}

// GetValueFromDataBucketAsFloat returns the value from the data bucket as a float
// The first key is the name of the bucket
// The second key is the name of the value
// The value is the data or nil if the value does not exist
func (m *Workflow) GetValueFromDataBucketAsFloat(key string, name string) float64 {
	if m.dataBucket == nil {
		m.dataBucket = make(map[string]map[string]interface{})
	}
	if m.dataBucket[key] == nil {
		m.dataBucket[key] = make(map[string]interface{})
	}
	if val, ok := m.dataBucket[key][name]; ok {
		return val.(float64)
	}
	return 0
}

// GetValueFromDataBucketAsBool returns the value from the data bucket as a bool
// The first key is the name of the bucket
// The second key is the name of the value
// The value is the data or nil if the value does not exist
func (m *Workflow) GetValueFromDataBucketAsBool(key string, name string) bool {
	if m.dataBucket == nil {
		m.dataBucket = make(map[string]map[string]interface{})
	}
	if m.dataBucket[key] == nil {
		m.dataBucket[key] = make(map[string]interface{})
	}
	if val, ok := m.dataBucket[key][name]; ok {
		return val.(bool)
	}
	return false
}

// GetDataBucket returns the data bucket
func (m *Workflow) GetDataBucket() map[string]map[string]interface{} {
	return m.dataBucket
}

//GetDataBucketItem returns the Map from the data bucket
// The first key is the name of the bucket
// returns Bucket Content nil if the bucket does not exist
func (m *Workflow) GetDataBucketContent(key string) map[string]interface{} {
	if m.dataBucket == nil {
		m.dataBucket = make(map[string]map[string]interface{})
	}
	if m.dataBucket[key] == nil {
		m.dataBucket[key] = make(map[string]interface{})
	}
	if val, ok := m.dataBucket[key]; ok {
		return val
	}
	return nil
}

// ClearDataBuckets clears the data bucket
func (m *Workflow) ClearDataBuckets() {
	m.dataBucket = make(map[string]map[string]interface{})
}

// ClearDataBucket values from the data bucket
// The first key is the name of the bucket
func (m *Workflow) ClearDataBucket(key string) {
	m.dataBucket[key] = make(map[string]interface{})
}

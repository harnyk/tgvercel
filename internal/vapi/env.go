package vapi

type EnvDescriptor struct {
	Type   string   `json:"type"`
	Value  string   `json:"value"`
	Target []Target `json:"target"`
	// ConfigurationId         string `json:"configurationId"`
	// Comment                 string `json:"comment"`
	ID  string `json:"id"`
	Key string `json:"key"`
	// CreatedAt               int    `json:"createdAt"`
	// UpdatedAt               int    `json:"updatedAt"`
	// CreatedBy               string `json:"createdBy"`
	// UpdatedBy               string `json:"updatedBy"`
	// Decrypted               bool   `json:"decrypted"`
	// LastEditedBy            string `json:"lastEditedBy"`
	// LastEditedByDisplayName string `json:"lastEditedByDisplayName"`
}

func (ed *EnvDescriptor) MatchTarget(target Target) bool {
	for _, t := range ed.Target {
		if t == target {
			return true
		}
	}
	return false
}

type EnvsResponse struct {
	Envs []EnvDescriptor `json:"envs"`
}

package vapi

type Deployment struct {
	TargetRaw string `json:"target"`
	Url       string `json:"url"`
	ProjectID string `json:"projectId"`
}

func (d *Deployment) Target() Target {
	if d.TargetRaw == "production" {
		return TargetProduction
	}
	return TargetPreview
}

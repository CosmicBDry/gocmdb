package forms

import (
	"encoding/json"
	"time"

	"github.com/CosmicBDry/gocmdb/server/models"
)

type AlertInfo struct {
	Status      string            `json:"status"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	StartsAt    *time.Time        `json:"startsAt"`
	EndsAt      *time.Time        `json:"endsAt"`
	FingerPrint string            `json:"fingerprint"`
}

type AlertForm struct {
	Alerts []AlertInfo `json:"alerts"`
}

func (f *AlertInfo) ToModel() *models.Alert {
	var (
		endtime    *time.Time
		formLabels string
	)

	if labels, err := json.Marshal(f.Labels); err == nil {
		formLabels = string(labels)
	}

	if f.EndsAt != nil && !f.EndsAt.IsZero() {
		endtime = f.EndsAt
	}
	return &models.Alert{
		AlertName:   f.Labels["alertname"],
		Instance:    f.Labels["instance"],
		Serverity:   f.Labels["serverity"],
		Status:      f.Status,
		FingerPrint: f.FingerPrint,
		Description: f.Annotations["description"],
		Summary:     f.Annotations["summary"],
		Labels:      formLabels,
		StartsAt:    f.StartsAt,
		EndsAt:      endtime,
	}
}

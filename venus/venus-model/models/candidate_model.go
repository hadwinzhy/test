package models

// Candidate ...
type Candidate struct {
	BaseModel
	EventID    uint    `gorm:"index" json:"event_id"`
	CustomerID uint    `gorm:"index" json:"customer_id"`
	Confidence float32 `json:"confidence"`
	PersonID   string  `gorm:"index" json:"person_id"`
	FaceURL    string  `gorm:"type:varchar(255)" json:"face_url"`
	// Source     string  `gorm:"default:readsense" json:"source"`
}

// CandidateBasicSerializer ...
type CandidateBasicSerializer struct {
	ID         uint    `json:"id"`
	PersonID   string  `json:"person_id"`
	Confidence float32 `json:"confidence"`
	FaceURL    string  `json:"face_url"`
}

// BasicSerialize ...
func (c *Candidate) BasicSerialize() CandidateBasicSerializer {
	return CandidateBasicSerializer{
		ID:         c.ID,
		PersonID:   c.PersonID,
		Confidence: c.Confidence,
		FaceURL:    c.FaceURL,
	}
}

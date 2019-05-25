package titan

import (
	"encoding/json"
)

// Face ...
type Face struct {
	FaceID string  `json:"face_id"`
	Age    int     `json:"age"`
	Gender int     `json:"gender"`
	Yaw    float64 `json:"yaw"`
	Roll   float64 `json:"roll"`
	Pitch  float64 `json:"pitch"`
	Pose   struct {
		Yaw   float64 `json:"yaw"`
		Roll  float64 `json:"roll"`
		Pitch float64 `json:"pitch"`
	} `json:"pose"`
	Rect struct {
		Left   int `json:"left"`
		Top    int `json:"top"`
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"rect"`
	Landmarks21 []struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"landmarks21"`
}

// FaceData ...
type FaceData struct {
	Status  string `json:"status"`
	ImageID string `json:"image_id"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Faces   []Face `json:"faces"`
}

// Person ...
type Person struct {
	Status    string `json:"status"`
	PersonID  string `json:"person_id"`
	Name      string `json:"name"`
	FaceCount int    `json:"face_count"`
}

// Group ...
type Group struct {
	Status string `json:"status"`
	// GroupID     string `json:"group_id"`
	GroupUUID   string `json:"group_uuid"`
	Name        string `json:"name"`
	PersonCount int    `json:"person_count"`
}

// CandidateData ...
type CandidateData struct {
	Status     string      `json:"status"`
	FaceID     string      `json:"face_id"`
	GroupID    string      `json:"group_id"`
	Candidates []Candidate `json:"candidates"`
}

// Candidate ...
type Candidate struct {
	ImageURL   string  `json:"image_url"`
	FaceID     string  `json:"face_id"`
	PersonID   string  `json:"person_id"`
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
}

////////////////////////////////// Parse Function /////////////////////////////

// ParseGroupData ...
func ParseGroupData(rawData interface{}, group *Group) {
	if error := json.Unmarshal([]byte(rawData.(string)), &group); error != nil {
		return
	}
}

// ParseFaceData ...
func ParseFaceData(rawData interface{}, faceData *FaceData) {
	if error := json.Unmarshal([]byte(rawData.(string)), &faceData); error != nil {
		return
	}
}

// ParsePersonData ...
func ParsePersonData(rawData interface{}, person *Person) {
	if error := json.Unmarshal([]byte(rawData.(string)), &person); error != nil {
		return
	}
}

// ParseCandidate ...
func ParseCandidate(rawData interface{}, candidateData *CandidateData) {
	if error := json.Unmarshal([]byte(rawData.(string)), &candidateData); error != nil {
		return
	}
}

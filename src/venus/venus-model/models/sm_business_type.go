package models

type SmBusinessType struct {
	BaseModel
	CompanyID uint   `gorm:"index" json:"company_id"`
	Name      string `gorm:"type:varchar(30);" json:"name"`
}

type SmBusinessTypes []SmBusinessType

func (s SmBusinessTypes) maxLength() bool {
	return len(s) <= 20
}

type SmBusinessTypeSerializer struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (s SmBusinessType) BasicSerializer() SmBusinessTypeSerializer {
	return SmBusinessTypeSerializer{
		ID:   s.ID,
		Name: s.Name,
	}
}

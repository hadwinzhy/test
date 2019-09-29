package models

// Province DB model
type Province struct {
	ID   uint   `gorm:"primary_key" json:"id"`
	Code string `gorm:"type:varchar(100);unique_index" json:"code"`
	Name string `gorm:"type:varchar(100);not null" json:"name"`
}

// City DB model
type City struct {
	ID         uint   `gorm:"primary_key" json:"id"`
	ProvinceID uint   `gorm:"index;not null" json:"-"`
	Code       string `gorm:"type:varchar(100);unique_index" json:"code"`
	Name       string `gorm:"type:varchar(100);not null" json:"name"`
	// District   District `gorm:"polymorphic:Province;" json:"-"`
}

// District DB model
type District struct {
	ID     uint   `gorm:"primary_key" json:"id"`
	CityID uint   `gorm:"index;not null" json:"-"`
	Code   string `gorm:"type:varchar(100);unique_index" json:"code"`
	Name   string `gorm:"type:varchar(100);not null" json:"name"`
}

// GaodeDistricts http://restapi.amap.com/v3/config/district?key=53e2f7980c7c67faca53dc4d2207690c&keywords=100000&subdistrict=3&extensions=base
type GaodeDistricts struct {
	Status     string `json:"status"`
	Info       string `json:"info"`
	Infocode   string `json:"infocode"`
	Count      string `json:"count"`
	Suggestion struct {
		Keywords []interface{} `json:"keywords"`
		Cities   []interface{} `json:"cities"`
	} `json:"suggestion"`
	Districts []struct {
		Citycode  []interface{} `json:"citycode"`
		Adcode    string        `json:"adcode"`
		Name      string        `json:"name"`
		Center    string        `json:"center"`
		Level     string        `json:"level"`
		Districts []struct {
			Citycode  string `json:"citycode"`
			Adcode    string `json:"adcode"`
			Name      string `json:"name"`
			Center    string `json:"center"`
			Level     string `json:"level"`
			Districts []struct {
				Citycode  string `json:"citycode"`
				Adcode    string `json:"adcode"`
				Name      string `json:"name"`
				Center    string `json:"center"`
				Level     string `json:"level"`
				Districts []struct {
					Citycode string `json:"citycode"`
					Adcode   string `json:"adcode"`
					Name     string `json:"name"`
					Center   string `json:"center"`
					Level    string `json:"level"`
				} `json:"districts"`
			} `json:"districts"`
		} `json:"districts"`
	} `json:"districts"`
}

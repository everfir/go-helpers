package structs

type BusinessStatus uint8

const (
	BUSINESS_STATUS_OFFLINE BusinessStatus = 0
	BUSINESS_STATUS_OK      BusinessStatus = 1
)

type BusinessConfig struct {
	BusinessList []Business `json:"business"`
	businessMap  map[string]Business
}

type Business struct {
	Id          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Ctime       string         `json:"ctime"`
	Status      BusinessStatus `json:"status"`
}

func (b *BusinessConfig) Format() {
	b.businessMap = make(map[string]Business)
	for _, business := range b.BusinessList {
		b.businessMap[business.Name] = business
	}
}

func (b BusinessConfig) Valid(business string) bool {
	if business == "" {
		return false
	}
	config, exist := b.businessMap[business]
	return exist && config.Status == BUSINESS_STATUS_OK
}

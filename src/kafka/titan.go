package kafka

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"siren/configs"
	"bitbucket.org/readsense/venus-model/models"
	"siren/pkg/database"
	"siren/pkg/utils"
	"siren/src/workers"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

func saveGroupInfo(companyID uint) (bool, *models.FrequentCustomerGroup) {
	var oneGroup models.FrequentCustomerGroup
	var (
		ok      bool
		groupID string
	)
	if dbError := database.POSTGRES.Where("company_id = ?", companyID).First(&oneGroup).Error; dbError != nil {
		oneGroup = models.FrequentCustomerGroup{
			CompanyID: companyID,
		}

		if ok, groupID = titanAddGroup(utils.GenerateUUID(20), fmt.Sprintf("%d_回头客", oneGroup.CompanyID)); !ok {
			return false, nil
		}
		oneGroup.GroupUUID = groupID

		if dbError := database.POSTGRES.Save(&oneGroup).Error; dbError != nil {
			return false, nil
		}

	}
	return true, &oneGroup
}

func titanAddGroup(groupUUID string, name string) (bool, string) {
	apiID := configs.FetchFieldValue("TitanAPIID")
	apiSecret := configs.FetchFieldValue("TitanAPISecret")
	response, err := http.PostForm(titanParams.groupCreateURL, url.Values{
		"api_id":     {apiID},
		"api_secret": {apiSecret},
		"group_id":   {groupUUID},
		"name":       {name},
	})
	if err != nil {
		log.Println("create group fail")
		return false, "-1"
	}
	defer response.Body.Close()
	log.Println("titan add group", response.StatusCode)
	content, _ := ioutil.ReadAll(response.Body)
	values := gjson.ParseBytes(content)
	if values.Get("status").String() != "ok" {
		log.Println("status is not ok, create group fail")
		return false, "-1"
	}
	groupID := values.Get("group_id").String()
	return true, groupID
}

func titanGroupAddPerson(groupUUID string, personID string) bool {
	// 将人加入到组中去
	apiID := configs.FetchFieldValue("TitanAPIID")
	apiSecret := configs.FetchFieldValue("TitanAPISecret")
	response, err := http.PostForm(titanParams.groupAddPerson, url.Values{
		"api_id":     {apiID},
		"api_secret": {apiSecret},
		"group_id":   {groupUUID},
		"person_id":  {personID},
	})
	if err != nil {
		return false
	}
	defer response.Body.Close()
	content, _ := ioutil.ReadAll(response.Body)
	values := gjson.ParseBytes(content)
	existsStatus := values.Get("status").Exists()
	if !existsStatus {
		return false
	}
	if values.Get("status").String() == "ok" {
		return true
	}
	return false
}

func fetchDataByTitan(group *models.FrequentCustomerGroup, info InfoForKafkaProducer) bool {
	log.Println("URL", titanParams.identificationURL)
	response, err := http.PostForm(titanParams.identificationURL, url.Values{
		"api_id":     {configs.FetchFieldValue("TitanAPIID")},
		"api_secret": {configs.FetchFieldValue("TitanAPISecret")},
		"face_id":    {info.FaceID},
		"group_id":   {group.GroupUUID},
		"top":        {"20"},
	})
	log.Println("response", response.StatusCode)
	if err != nil {
		return false
	}
	defer response.Body.Close()

	responseByte, _ := ioutil.ReadAll(response.Body)

	log.Println("titan values", string(responseByte))
	// todo: fix it if status is not ok
	if info.CompanyID != 0 {
		if ok := personIDHandler(info.EventID, group.ID, info.PersonID, responseByte, info.CapturedAt, info.EventStatus); !ok {
			return false
		}
	}

	return true

}

type result struct {
	PersonID string    `json:"person_id"`
	Day      time.Time `json:"day"`
}

type results []result

func personIDHandler(eventID uint, groupID uint, personUUID string, values []byte, capturedAt int64, status string) bool {
	if status != "analyzed" {
		var onePerson models.FrequentCustomerPeople
		if personUUID == "" {
			onePerson.PersonID = utils.GenerateUUID(20)
		} else {
			onePerson.PersonID = personUUID
		}
		onePerson.Date = utils.CurrentDate(time.Unix(capturedAt, 0))
		hour := utils.CurrentTime(time.Unix(capturedAt, 0), "hour")
		onePerson.Hour = hour
		onePerson.Frequency = 1 //  这次来的，加1
		onePerson.Interval = 0
		onePerson.FrequentCustomerGroupID = groupID
		onePerson.IsFrequentCustomer = false
		onePerson.EventID = eventID
		database.POSTGRES.Save(&onePerson)
		workers.MallCountFrequentCustomerHandler(onePerson, groupID, capturedAt)
		return true
	}
	valuesJson := gjson.ParseBytes(values)
	exists := valuesJson.Get("candidates").Exists()
	if !exists || (exists && len(valuesJson.Get("candidates").Array()) == 0) {
		var onePerson models.FrequentCustomerPeople
		if personUUID == "" {
			onePerson.PersonID = utils.GenerateUUID(20)
		} else {
			onePerson.PersonID = personUUID
		}
		onePerson.Date = utils.CurrentDate(time.Unix(capturedAt, 0))
		hour := utils.CurrentTime(time.Unix(capturedAt, 0), "hour")
		onePerson.Hour = hour
		onePerson.Frequency = 1 //  这次来的，加1
		onePerson.Interval = 0
		onePerson.FrequentCustomerGroupID = groupID
		onePerson.IsFrequentCustomer = false
		onePerson.EventID = eventID
		database.POSTGRES.Save(&onePerson)
		workers.MallCountFrequentCustomerHandler(onePerson, groupID, capturedAt)
		return true
	} else {
		var personIDs []string
		for _, i := range valuesJson.Get("candidates").Array() {
			personIDs = append(personIDs, fmt.Sprintf("'%s'", i.Get("person_id").String()))
		}
		personIDString := strings.Join(personIDs, ",")
		now := time.Now()
		right := now.Format("2006-01-02 15:04:05")
		left := now.AddDate(0, -1, 0).Format("2006-01-02 15:04:05")
		sql := fmt.Sprintf(`SELECT person_id, date_trunc('day',max(capture_at)) as day FROM events WHERE person_id in (%s) AND capture_at BETWEEN '%s' AND '%s' group by person_id ORDER BY day desc`,
			personIDString, left, right)

		var resultsValues results

		database.POSTGRES.Raw(sql).Scan(&resultsValues)

		//personID
		var onePerson models.FrequentCustomerPeople
		hour := utils.CurrentTime(time.Unix(capturedAt, 0), "hour")
		if dbError := database.POSTGRES.Where("person_id = ? AND hour = ?", personUUID, hour).First(&onePerson).Error; dbError != nil {
			onePerson = models.FrequentCustomerPeople{
				PersonID:                personUUID,
				FrequentCustomerGroupID: groupID,
				Date:                    utils.CurrentDate(time.Unix(capturedAt, 0)),
				Hour:                    hour,
				Frequency:               uint(len(resultsValues)), // 算上这次 +1
				EventID:                 eventID,
			}
			if len(resultsValues) <= 1 {
				onePerson.Interval = 0 // 新客，间隔为 0
				onePerson.IsFrequentCustomer = false
			} else {
				onePerson.Interval = uint(float64(time.Now().Sub(resultsValues[0].Day).Hours()/24) + 1)
				onePerson.IsFrequentCustomer = true
			}
			if dbError := database.POSTGRES.Save(&onePerson).Error; dbError != nil {
				return false
			}
		}
		workers.MallCountFrequentCustomerHandler(onePerson, groupID, capturedAt)
	}
	return true
}

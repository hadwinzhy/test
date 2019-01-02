package titan

import (
	"fmt"
	"os"
	"siren/configs"
	"strings"

	"github.com/parnurzeal/gorequest"
)

const (
	facesDetectionPath           = "/faces/detection"
	facesVerificationPath        = "/faces/verification"
	facesIdentificationPath      = "/faces/identification"
	facesImageIdentificationPath = "/faces/image_identification"

	// people api
	peopleCreatePath     = "/people/create"
	peopleDeletePath     = "/people/delete"
	peopleAddFacePath    = "/people/add_face"
	peopleRemoveFacePath = "/people/remove_face"
	peopleEmptyPath      = "/people/empty"

	// groups api
	groupsCreatePath       = "/groups/create"
	groupsDeletePath       = "/groups/delete"
	groupsAddPersonPath    = "/groups/add_person"
	groupsRemovePersonPath = "/groups/remove_person"
	groupsEmptyPath        = "/groups/empty"
)

// APIManager used to handle it
type APIManager struct {
	Host      string
	APIID     string
	APISecret string
}

var managers map[string]*APIManager

func GetAPIManager(serverName string) *APIManager {
	if managers == nil {
		managers = make(map[string]*APIManager)
	}

	// check map
	if managers[serverName] != nil {
		return managers[serverName]
	}

	newManager := new(APIManager)
	newManager.Host = configs.FetchFieldValue("titan" + serverName + "_host")
	newManager.APIID = configs.FetchFieldValue("titan" + serverName + "_api_id")
	newManager.APISecret = configs.FetchFieldValue("titan" + serverName + "_api_secret")

	managers[serverName] = newManager

	return newManager
}

// 非文件类型参数检验和添加
func (m *APIManager) simpleParamsVerify(params map[string]interface{}) (bool, string, map[string]interface{}) {
	titanParamsMap := map[string]interface{}{
		"api_id":     m.APIID,
		"api_secret": m.APISecret,
	}

	for key, param := range params {
		if param == "" {
			return false, ("miss param " + key), titanParamsMap
		} else {
			titanParamsMap[key] = param
		}
	}
	return true, "", titanParamsMap
}

// 非文件上传的post request
func simplePostRequest(url string, p map[string]interface{}) (bool, interface{}) {
	_, body, errs := gorequest.New().
		Post(url).
		SendStruct(p).
		End()

	if errs != nil {
		return false, errs
	}
	return true, body
}

// face detection
func (m *APIManager) TitanFacesDetection(filePath string) (bool, interface{}) {
	fmt.Println("detect 1", filePath)

	if filePath == "" {
		return false, "miss param file"
	}

	titanParamsMap := map[string]interface{}{
		"api_id":     m.APIID,
		"api_secret": m.APISecret,
	}

	req := gorequest.New()

	if strings.Index(filePath, "http") == -1 {
		fileContent, err := os.Open(filePath)
		if err != nil {
			return false, err
		}
		defer fileContent.Close()

		_, body, errs := req.Post(m.Host+facesDetectionPath).
			Type("multipart").
			SendStruct(titanParamsMap).
			SendFile(fileContent, "", "file").
			End()

		fmt.Println("body and err 1", body, errs)
		if errs != nil {
			return false, errs
		}

		return true, body
	} else {
		titanParamsMap["url"] = filePath

		_, body, errs := req.Post(m.Host + facesDetectionPath).
			SendStruct(titanParamsMap).
			End()

		fmt.Println("body and err 1", body, errs)

		if errs != nil {
			return false, errs
		}

		return true, body
	}
}

func (m *APIManager) TitanFacesVerification(fid1, fid2, pid string) (bool, interface{}) {
	rb, str, p := m.simpleParamsVerify(map[string]interface{}{
		"face_id": fid1,
	})
	if !rb {
		return false, str
	}

	if fid2 == "" {
		if pid == "" {
			return false, "miss param face_id2 or person_id"
		} else {
			p["person_id"] = pid
		}
	} else {
		p["face_id2"] = fid2
	}

	return simplePostRequest(m.Host+facesVerificationPath, p)
}

func (m *APIManager) TitanFacesIdentification(fid, gid, threshold string) (bool, interface{}) {
	rb, str, p := m.simpleParamsVerify(map[string]interface{}{
		"face_id":  fid,
		"group_id": gid,
	})
	if !rb {
		return false, str
	}

	if threshold != "" {
		p["threshold"] = threshold
	}
	return simplePostRequest(m.Host+facesIdentificationPath, p)
}

func (m *APIManager) TitanFacesImageIdentification(filePath, gid string) (bool, interface{}) {
	if filePath == "" {
		return false, "miss param file"
	}

	rb, str, p := m.simpleParamsVerify(map[string]interface{}{
		"group_id": gid,
	})
	if !rb {
		return false, str
	}

	fileContent, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer fileContent.Close()

	_, body, errs := gorequest.New().Post(m.Host+facesImageIdentificationPath).
		Type("multipart").
		SendStruct(p).
		SendFile(fileContent, "", "file").
		End()

	if errs != nil {
		return false, errs
	}
	return true, body
}

func (m *APIManager) TitanPeopleCreate(fid, name string) (bool, interface{}) {
	rb, str, p := m.simpleParamsVerify(map[string]interface{}{
		"face_id": fid,
		"name":    name,
	})
	if !rb {
		return false, str
	}

	return simplePostRequest(m.Host+peopleCreatePath, p)
}

func (m *APIManager) TitanPeopleDelete(pid string) (bool, interface{}) {
	rb, str, p := m.simpleParamsVerify(map[string]interface{}{
		"person_id": pid,
	})
	if !rb {
		return false, str
	}

	return simplePostRequest(m.Host+peopleDeletePath, p)
}

func (m *APIManager) TitanPeopleAddFace(fid, pid string) (bool, interface{}) {
	rb, str, p := m.simpleParamsVerify(map[string]interface{}{
		"face_id":   fid,
		"person_id": pid,
	})
	if !rb {
		return false, str
	}

	return simplePostRequest(m.Host+peopleAddFacePath, p)
}

func (m *APIManager) TitanPeopleRemoveFace(fid, pid string) (bool, interface{}) {
	rb, str, p := m.simpleParamsVerify(map[string]interface{}{
		"face_id":   fid,
		"person_id": pid,
	})
	if !rb {
		return false, str
	}

	return simplePostRequest(m.Host+peopleRemoveFacePath, p)
}

func (m *APIManager) TitanPeopleEmpty(pid string) (bool, interface{}) {
	rb, str, p := m.simpleParamsVerify(map[string]interface{}{
		"person_id": pid,
	})
	if !rb {
		return false, str
	}

	return simplePostRequest(m.Host+peopleEmptyPath, p)
}

func (m *APIManager) TitanGroupsCreate(pid, name string) (bool, interface{}) {
	rb, str, p := m.simpleParamsVerify(map[string]interface{}{
		"name": name,
	})
	if !rb {
		return false, str
	}

	if pid != "" {
		p["person_id"] = pid
	}

	return simplePostRequest(m.Host+groupsCreatePath, p)
}

func (m *APIManager) TitanGroupsDelete(gid string) (bool, interface{}) {
	rb, str, p := m.simpleParamsVerify(map[string]interface{}{
		"group_id": gid,
	})
	if !rb {
		return false, str
	}

	return simplePostRequest(m.Host+groupsDeletePath, p)
}

func (m *APIManager) TitanGroupsAddPerson(pid, gid string) (bool, interface{}) {
	rb, str, p := m.simpleParamsVerify(map[string]interface{}{
		"person_id": pid,
		"group_id":  gid,
	})
	if !rb {
		return false, str
	}

	return simplePostRequest(m.Host+groupsAddPersonPath, p)
}

func (m *APIManager) TitanGroupsRemovePerson(pid, gid string) (bool, interface{}) {
	rb, str, p := m.simpleParamsVerify(map[string]interface{}{
		"person_id": pid,
		"group_id":  gid,
	})
	if !rb {
		return false, str
	}

	return simplePostRequest(m.Host+groupsRemovePersonPath, p)
}

func (m *APIManager) TitanGroupsEmpty(gid string) (bool, interface{}) {
	rb, str, p := m.simpleParamsVerify(map[string]interface{}{
		"group_id": gid,
	})
	if !rb {
		return false, str
	}

	return simplePostRequest(m.Host+groupsEmptyPath, p)
}

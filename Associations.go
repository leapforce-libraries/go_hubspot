package hubspot

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
)

const maxBatchSize int = 10000

type AssociationsSet struct {
	Results []Association `json:"results"`
}

type Association struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type AssociationsV4Set struct {
	Results []AssociationV4 `json:"results"`
}

type AssociationV4 struct {
	From struct {
		Id string `json:"id"`
	} `json:"from"`
	To []struct {
		ToObjectId       int64             `json:"toObjectId"`
		AssociationTypes []AssociationType `json:"associationTypes"`
	} `json:"to"`
}

type AssociationType struct {
	Category string `json:"category"`
	TypeId   int64  `json:"typeId"`
	Label    string `json:"label"`
}

type BatchGetAssociationsConfig struct {
	FromObjectType ObjectType
	ToObjectType   ObjectType
	Ids            []string
}

func (service *Service) BatchGetAssociations(config *BatchGetAssociationsConfig) (*AssociationsV4Set, *errortools.Error) {
	if len(config.Ids) == 0 {
		return nil, nil
	}

	endpoint := fmt.Sprintf("associations/%v/%v/batch/read", config.FromObjectType, config.ToObjectType)

	ids := config.Ids
	var associationsV4Set AssociationsV4Set

	for len(ids) > 0 {

		var body struct {
			Inputs []struct {
				Id string `json:"id"`
			} `json:"inputs"`
		}

		for i, id := range ids {
			if i == maxBatchSize {
				break
			}
			idStruct := struct {
				Id string `json:"id"`
			}{id}
			body.Inputs = append(body.Inputs, idStruct)
		}

		var associationsV4Set_ AssociationsV4Set

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodPost,
			Url:           service.urlV4(endpoint),
			BodyModel:     body,
			ResponseModel: &associationsV4Set_,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		associationsV4Set.Results = append(associationsV4Set.Results, associationsV4Set_.Results...)

		if len(ids) > maxBatchSize {
			ids = ids[maxBatchSize:]
		} else {
			break
		}
	}

	return &associationsV4Set, nil
}

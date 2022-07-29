package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/ktr03rtk/touring-log-service/api-backend/domain/model"
	"github.com/pkg/errors"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

func (h *handler) storePhoto(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(defaultMaxMemory); err != nil {
		h.errJSON(w, errors.Wrapf(err, "failed to parse multipartform photo data"))

		return
	}

	id, unit, err := getAccountInfo(r)
	if err != nil {
		h.errJSON(w, err)

		return
	}

	files := r.MultipartForm.File["images"]
	for _, file := range files {
		if err := h.photoUsecase.Execute(file, id, unit); err != nil {
			h.errJSON(w, err)

			return
		}
	}

	res := jsonResp{
		OK: true,
	}

	if err := h.writeJSON(w, http.StatusOK, res, "response"); err != nil {
		h.errJSON(w, err)
		return
	}
}

type TripPayload struct {
	Year  int    `json:"year"`
	Month int    `json:"month"`
	Day   int    `json:"day"`
	Unit  string `json:"unit"`
}

func (h *handler) storeTrip(w http.ResponseWriter, r *http.Request) {
	var payload TripPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.errJSON(w, errors.Wrapf(err, "failed to json decode trip"))

		return
	}

	if err := h.tripUsecase.Execute(payload.Year, payload.Month, payload.Day, payload.Unit); err != nil {
		h.errJSON(w, err)

		return
	}

	res := jsonResp{
		OK: true,
	}

	if err := h.writeJSON(w, http.StatusOK, res, "response"); err != nil {
		h.errJSON(w, err)
		return
	}
}

var touringLogs []*model.TouringLog

func (h *handler) graphQLFileds(r *http.Request) graphql.Fields {
	return graphql.Fields{
		"list": &graphql.Field{
			Type:        graphql.NewList(touringLogType),
			Description: "Get log by year and month",
			Args: graphql.FieldConfigArgument{
				"year": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"month": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				year, ok := p.Args["year"].(int)
				if !ok {
					return nil, errors.New("failed to specify year")
				}

				month, ok := p.Args["month"].(int)
				if !ok {
					return nil, errors.New("failed to specify month")
				}

				id, unit, err := getAccountInfo(r)
				if err != nil {
					return nil, errors.New("failed to specify identity")
				}

				touringLogs, err := h.listQueryUsecase.Execute(year, month, id, unit)
				if err != nil {
					return nil, errors.New("failed to fetch log")
				}

				return touringLogs, nil
			},
		},
	}
}

var touringLogType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "TouringLog",
		Fields: graphql.Fields{
			"year": &graphql.Field{
				Type: graphql.Int,
			},
			"month": &graphql.Field{
				Type: graphql.Int,
			},
			"day": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

func (h *handler) graphQL(w http.ResponseWriter, r *http.Request) {
	q, _ := io.ReadAll(r.Body)
	query := string(q)

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: h.graphQLFileds(r)}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}

	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		h.errJSON(w, err)
		return
	}

	params := graphql.Params{Schema: schema, RequestString: query}
	resp := graphql.Do(params)

	if len(resp.Errors) > 0 {
		h.errJSON(w, errors.New(fmt.Sprintf("failed: %+v", resp.Errors)))
		return
	}

	j, _ := json.MarshalIndent(resp, "", " ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

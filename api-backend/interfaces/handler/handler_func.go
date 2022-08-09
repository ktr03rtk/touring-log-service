package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/julienschmidt/httprouter"
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
		if err := h.photoStoreUsecase.Execute(file, id, unit); err != nil {
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

func (h *handler) getPhoto(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	reader, err := h.photoGetUsecase.Execute(params.ByName("id"))
	if err != nil {
		h.errJSON(w, err)
		return
	}

	buf, err := io.ReadAll(reader)
	if err != nil {
		h.errJSON(w, err)
		return
	}

	res := jsonResp{
		OK:      true,
		Message: base64.StdEncoding.EncodeToString(buf),
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

var logDate []*model.LogDate

func (h *handler) graphQLFileds(r *http.Request) graphql.Fields {
	return graphql.Fields{
		"dateList": &graphql.Field{
			Type:        graphql.NewList(logDateType),
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

				logDates, err := h.listQueryUsecase.Execute(year, month, id, unit)
				if err != nil {
					return nil, errors.New("failed to fetch log")
				}

				return logDates, nil
			},
		},
		"touringLog": &graphql.Field{
			Type:        touringLogType,
			Description: "Get log by year and month and day",
			Args: graphql.FieldConfigArgument{
				"year": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"month": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"day": &graphql.ArgumentConfig{
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

				day, ok := p.Args["day"].(int)
				if !ok {
					return nil, errors.New("failed to specify day")
				}

				id, unit, err := getAccountInfo(r)
				if err != nil {
					return nil, errors.New("failed to specify identity")
				}

				fmt.Printf("--------------- %+v\n", unit)
				// TODO: fetch trip

				photoLog, err := h.photoLogQueryUsecase.Execute(year, month, day, id)
				if err != nil {
					return nil, errors.New("failed to fetch log")
				}

				touringLogs := model.TouringLog{
					Photo: photoLog,
				}

				return touringLogs, nil
			},
		},
	}
}

var logDateType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "LogDate",
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

var touringLogType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "TouringLog",
		Fields: graphql.Fields{
			"trip": &graphql.Field{
				Type: graphql.NewList(tripLogType),
			},
			"photo": &graphql.Field{
				Type: graphql.NewList(photoLogType),
			},
		},
	},
)

var tripLogType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "TripLog",
		Fields: graphql.Fields{
			"lat": &graphql.Field{
				Type: graphql.Float,
			},
			"lng": &graphql.Field{
				Type: graphql.Float,
			},
		},
	},
)

var photoLogType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "PhotoLog",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"lat": &graphql.Field{
				Type: graphql.Float,
			},
			"lng": &graphql.Field{
				Type: graphql.Float,
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

package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/lfroomin/restaurant-container/internal/model"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubError struct {
	restaurant string
	location   string
}

func Test_Create(t *testing.T) {
	restName := "Rest 1"
	restaurantExp, _ := json.Marshal(model.Restaurant{
		Name: restName,
		Address: &model.Address{
			Location:     &model.Location{},
			TimezoneName: new(string),
		},
	})
	restaurantNoAddressExp, _ := json.Marshal(model.Restaurant{
		Name: restName,
	})
	testCases := []struct {
		name         string
		restaurant   model.Restaurant
		emptyReqBody bool
		responseCode int
		responseBody string
		stubError    stubError
	}{
		{"happy path",
			model.Restaurant{
				Name:    restName,
				Address: &model.Address{},
			},
			false,

			http.StatusCreated,
			string(restaurantExp),
			stubError{},
		},
		{"no address",
			model.Restaurant{
				Name: restName,
			},
			false,
			http.StatusCreated,
			string(restaurantNoAddressExp),
			stubError{},
		},
		{"storage error",
			model.Restaurant{},
			false,
			http.StatusInternalServerError,
			`{"Message":"an error occurred"}`,
			stubError{restaurant: "an error occurred"},
		},
		{"location error",
			model.Restaurant{
				Name:    restName,
				Address: &model.Address{},
			},
			false,
			http.StatusInternalServerError,
			`{"Message":"an error occurred"}`,
			stubError{location: "an error occurred"},
		},
		{"empty request body",
			model.Restaurant{},
			true,
			http.StatusBadRequest,
			`{"Message":"error binding request body"}`,
			stubError{},
		},
	}

	for _, tc := range testCases {
		// scoped variable
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rc := RestaurantController{
				Restaurant: restaurantStorerStub{error: tc.stubError.restaurant},
				Location:   locationServiceStub{error: tc.stubError.location},
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = &http.Request{
				Body: io.NopCloser(bytes.NewBuffer([]byte{})),
			}
			if !tc.emptyReqBody {
				b, _ := json.Marshal(tc.restaurant)
				c.Request.Body = io.NopCloser(bytes.NewBuffer(b))
			}

			rc.Create(c)

			assert.Equal(t, tc.responseCode, w.Code)

			if tc.responseCode != http.StatusCreated {
				assert.Equal(t, tc.responseBody, w.Body.String())
			} else {
				// Convert to type Restaurant so comparison can be done
				// without the "Id" field
				expRestaurant := model.Restaurant{}
				_ = json.Unmarshal([]byte(tc.responseBody), &expRestaurant)
				restaurant := model.Restaurant{}
				_ = json.Unmarshal(w.Body.Bytes(), &restaurant)

				diff := cmp.Diff(
					expRestaurant,
					restaurant,
					cmpopts.IgnoreFields(model.Restaurant{}, "Id"),
				)
				assert.Empty(t, diff)
			}
		})
	}
}

func Test_Read(t *testing.T) {
	testCases := []struct {
		name         string
		restaurantId string
		exists       bool
		responseCode int
		responseBody string
		stubError    string
	}{
		{"happy path",
			"restId",
			true,
			http.StatusOK,
			`{"name":""}`,
			"",
		},
		{"empty restaurantId",
			"",
			true,
			http.StatusBadRequest,
			`{"Message":"restaurantId is empty"}`,
			"",
		},
		{"storage error",
			"restId",
			true,
			http.StatusInternalServerError,
			`{"Message":"an error occurred"}`,
			"an error occurred",
		},
		{"restaurant does not exist",
			"restId",
			false,
			http.StatusNotFound,
			"",
			"",
		},
	}

	for _, tc := range testCases {
		// scoped variable
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rc := RestaurantController{
				Restaurant: restaurantStorerStub{exists: tc.exists, error: tc.stubError},
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = []gin.Param{{Key: "restaurantId", Value: tc.restaurantId}}

			c.Request = &http.Request{
				Body: io.NopCloser(bytes.NewBuffer([]byte{})),
			}

			rc.Read(c)

			assert.Equal(t, tc.responseCode, w.Code)
			assert.Equal(t, tc.responseBody, w.Body.String())
		})
	}
}

func Test_Update(t *testing.T) {
	restId, restName := "Rest1", "Rest 1"
	restaurantExp, _ := json.Marshal(model.Restaurant{
		Id:   &restId,
		Name: restName,
		Address: &model.Address{
			Location:     &model.Location{},
			TimezoneName: new(string),
		},
	})
	restaurantNoAddressExp, _ := json.Marshal(model.Restaurant{
		Id:   &restId,
		Name: restName,
	})

	testCases := []struct {
		name         string
		restaurantId string
		restaurant   model.Restaurant
		emptyReqBody bool
		responseCode int
		responseBody string
		stubError    stubError
	}{
		{"happy path",
			restId,
			model.Restaurant{
				Id:      &restId,
				Name:    restName,
				Address: &model.Address{},
			},
			false,

			http.StatusOK,
			string(restaurantExp),
			stubError{},
		},
		{"no address",
			restId,
			model.Restaurant{
				Id:   &restId,
				Name: restName,
			},
			false,
			http.StatusOK,
			string(restaurantNoAddressExp),
			stubError{},
		},
		{"restaurantId is nil",
			restId,
			model.Restaurant{
				Name: restName,
			},
			false,
			http.StatusBadRequest,
			`{"Message":"restaurantId in URL path parameters and restaurant in body do not match"}`,
			stubError{},
		},
		{"mismatch restaurantId",
			"differentRestId",
			model.Restaurant{
				Id:   &restId,
				Name: restName,
			},
			false,
			http.StatusBadRequest,
			`{"Message":"restaurantId in URL path parameters and restaurant in body do not match"}`,
			stubError{},
		},
		{"storage error",
			restId,
			model.Restaurant{Id: &restId},
			false,
			http.StatusInternalServerError,
			`{"Message":"an error occurred"}`,
			stubError{restaurant: "an error occurred"},
		},
		{"location error",
			restId,
			model.Restaurant{
				Id:      &restId,
				Name:    restName,
				Address: &model.Address{},
			},
			false,
			http.StatusInternalServerError,
			`{"Message":"an error occurred"}`,
			stubError{location: "an error occurred"},
		},
		{"empty request body",
			"",
			model.Restaurant{},
			true,
			http.StatusBadRequest,
			`{"Message":"error binding request body"}`,
			stubError{},
		},
	}

	for _, tc := range testCases {
		// scoped variable
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rc := RestaurantController{
				Restaurant: restaurantStorerStub{error: tc.stubError.restaurant},
				Location:   locationServiceStub{error: tc.stubError.location},
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = []gin.Param{{Key: "restaurantId", Value: tc.restaurantId}}

			c.Request = &http.Request{
				Body: io.NopCloser(bytes.NewBuffer([]byte{})),
			}
			if !tc.emptyReqBody {
				b, _ := json.Marshal(tc.restaurant)
				c.Request.Body = io.NopCloser(bytes.NewBuffer(b))
			}

			rc.Update(c)

			assert.Equal(t, tc.responseCode, w.Code)
			assert.Equal(t, tc.responseBody, w.Body.String())
		})
	}
}

func Test_Delete(t *testing.T) {
	testCases := []struct {
		name         string
		restaurantId string
		responseCode int
		responseBody string
		stubError    string
	}{
		{"happy path",
			"restId",
			http.StatusOK,
			`""`,
			"",
		},
		{"empty restaurantId",
			"",
			http.StatusBadRequest,
			`{"Message":"restaurantId is empty"}`,
			"",
		},
		{"storage error",
			"restId",
			http.StatusInternalServerError,
			`{"Message":"an error occurred"}`,
			"an error occurred",
		},
	}

	for _, tc := range testCases {
		// scoped variable
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rc := RestaurantController{
				Restaurant: restaurantStorerStub{error: tc.stubError},
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = []gin.Param{{Key: "restaurantId", Value: tc.restaurantId}}

			rc.Delete(c)

			assert.Equal(t, tc.responseCode, w.Code)
			assert.Equal(t, tc.responseBody, w.Body.String())
		})
	}
}

type restaurantStorerStub struct {
	exists bool
	error  string
}

func (s restaurantStorerStub) Save(_ model.Restaurant) error {
	if s.error != "" {
		return errors.New(s.error)
	}
	return nil
}

func (s restaurantStorerStub) Get(_ string) (model.Restaurant, bool, error) {
	if s.error != "" {
		return model.Restaurant{}, false, errors.New(s.error)
	}
	if !s.exists {
		return model.Restaurant{}, false, nil
	}
	return model.Restaurant{}, true, nil
}

func (s restaurantStorerStub) Update(_ model.Restaurant) error {
	if s.error != "" {
		return errors.New(s.error)
	}
	return nil
}

func (s restaurantStorerStub) Delete(_ string) error {
	if s.error != "" {
		return errors.New(s.error)
	}
	return nil
}

type locationServiceStub struct {
	error string
}

func (s locationServiceStub) Geocode(_ model.Address) (model.Location, string, error) {
	if s.error != "" {
		return model.Location{}, "", errors.New(s.error)
	}
	return model.Location{}, "", nil
}

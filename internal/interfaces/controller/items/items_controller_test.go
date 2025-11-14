package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"Aicon-assignment/internal/domain/entity"
	domainErrors "Aicon-assignment/internal/domain/errors"
	"Aicon-assignment/internal/usecase"
)

type mockItemUsecase struct {
	updateItemFunc func(ctx context.Context, id int64, input usecase.UpdateItemInput) (*entity.Item, error)
}

func (m *mockItemUsecase) GetAllItems(ctx context.Context) ([]*entity.Item, error) {
	return nil, nil
}

func (m *mockItemUsecase) GetItemByID(ctx context.Context, id int64) (*entity.Item, error) {
	return nil, nil
}

func (m *mockItemUsecase) CreateItem(ctx context.Context, input usecase.CreateItemInput) (*entity.Item, error) {
	return nil, nil
}

func (m *mockItemUsecase) UpdateItem(ctx context.Context, id int64, input usecase.UpdateItemInput) (*entity.Item, error) {
	if m.updateItemFunc != nil {
		return m.updateItemFunc(ctx, id, input)
	}
	return nil, nil
}

func (m *mockItemUsecase) DeleteItem(ctx context.Context, id int64) error {
	return nil
}

func (m *mockItemUsecase) GetCategorySummary(ctx context.Context) (*usecase.CategorySummary, error) {
	return nil, nil
}

func TestItemHandler_UpdateItem(t *testing.T) {
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		mockUsecase := &mockItemUsecase{}
		mockUsecase.updateItemFunc = func(ctx context.Context, id int64, input usecase.UpdateItemInput) (*entity.Item, error) {
			assert.Equal(t, int64(1), id)
			if assert.NotNil(t, input.Name) {
				assert.Equal(t, "Updated", *input.Name)
			}
			return &entity.Item{ID: 1, Name: "Updated"}, nil
		}

		handler := NewItemHandler(mockUsecase)

		body := []byte(`{"name":"Updated"}`)
		req := httptest.NewRequest(http.MethodPatch, "/items/1", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/items/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := handler.UpdateItem(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var actual entity.Item
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actual))
		assert.Equal(t, int64(1), actual.ID)
		assert.Equal(t, "Updated", actual.Name)
	})

	t.Run("invalid id", func(t *testing.T) {
		handler := NewItemHandler(&mockItemUsecase{})
		req := httptest.NewRequest(http.MethodPatch, "/items/abc", bytes.NewReader([]byte(`{"name":"Updated"}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/items/:id")
		c.SetParamNames("id")
		c.SetParamValues("abc")

		err := handler.UpdateItem(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		handler := NewItemHandler(&mockItemUsecase{})
		req := httptest.NewRequest(http.MethodPatch, "/items/1", bytes.NewReader([]byte(`{}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/items/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := handler.UpdateItem(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockUsecase := &mockItemUsecase{}
		mockUsecase.updateItemFunc = func(ctx context.Context, id int64, input usecase.UpdateItemInput) (*entity.Item, error) {
			return nil, domainErrors.ErrItemNotFound
		}

		handler := NewItemHandler(mockUsecase)
		req := httptest.NewRequest(http.MethodPatch, "/items/1", bytes.NewReader([]byte(`{"name":"Updated"}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/items/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := handler.UpdateItem(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("domain validation error", func(t *testing.T) {
		mockUsecase := &mockItemUsecase{}
		mockUsecase.updateItemFunc = func(ctx context.Context, id int64, input usecase.UpdateItemInput) (*entity.Item, error) {
			return nil, domainErrors.ErrInvalidInput
		}

		handler := NewItemHandler(mockUsecase)
		req := httptest.NewRequest(http.MethodPatch, "/items/1", bytes.NewReader([]byte(`{"name":"Updated"}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/items/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := handler.UpdateItem(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

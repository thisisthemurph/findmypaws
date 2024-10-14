package routes

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"paws/internal/store"
	"paws/internal/types"
	"time"
)

func NewPetsHandler(s *store.PostgresStore, logger *slog.Logger) PetsHandler {
	return PetsHandler{
		PetStore: s.PetStore,
		Logger:   logger,
	}
}

type PetsHandler struct {
	PetStore store.PetStore
	Logger   *slog.Logger
}

func (h PetsHandler) MakeRoutes(g *echo.Group) {
	g.GET("/pets/:id", h.GetPetByID())
	g.POST("/pets", h.CreateNewPet())
	g.PUT("/pets", h.UpdatePet())
}

type NewPetRequest struct {
	Name string        `json:"name" validate:"required"`
	Tags types.PetTags `json:"tags"`
}

func (h PetsHandler) GetPetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad identifier")
		}

		pet, err := h.PetStore.Pet(id)
		if err != nil {
			if notFound := errors.As(err, &store.ErrPetNotFound); notFound {
				return echo.NewHTTPError(http.StatusNotFound, err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, pet)
	}
}

func (h PetsHandler) CreateNewPet() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req NewPetRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		pet := &types.Pet{
			Name: req.Name,
			Tags: req.Tags,
			DOB:  nil,
		}

		if err := h.PetStore.Create(pet); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusCreated, pet)
	}
}

type UpdatePetRequest struct {
	ID   uuid.UUID      `json:"id" validate:"required"`
	Type *types.PetType `json:"type" validate:"max=16"`
	Name string         `json:"name"`
	DOB  *time.Time     `json:"dob"`
	Tags types.PetTags  `json:"tags"`
}

func (h PetsHandler) UpdatePet() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req UpdatePetRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		pet := &types.Pet{
			ID:   req.ID,
			Type: req.Type,
			Name: req.Name,
			Tags: req.Tags,
			DOB:  req.DOB,
		}

		if err := h.PetStore.Update(pet); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, pet)
	}
}

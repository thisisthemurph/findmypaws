package routes

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"io"
	"log/slog"
	"net/http"
	"os"
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
	g.GET("/pets", h.ListPets())
	g.POST("/pets", h.CreateNewPet())
	g.PUT("/pets", h.UpdatePet())
	g.POST("/pets/:id/tag", h.AddTag())
	g.DELETE("/pets/:id/tag/:key", h.DeleteTag())
	g.DELETE("/pets/:id", h.DeletePet())
	g.PUT("/pets/:id/avatar", h.UpdateImage())
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

func (h PetsHandler) ListPets() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := clerkUser(c)
		if !user.Authenticated {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		pets, err := h.PetStore.Pets(user.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, pets)
	}
}

type NewPetRequest struct {
	Name string         `json:"name" validate:"required"`
	Type *types.PetType `json:"type"`
	Tags types.PetTags  `json:"tags"`
	DOB  *time.Time     `json:"dob"`
}

func (h PetsHandler) CreateNewPet() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := clerkUser(c)
		if !user.Authenticated {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		var req NewPetRequest
		if err := c.Bind(&req); err != nil {
			h.Logger.Error("bad request", "error", err)
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		currentPets, err := h.PetStore.Pets(user.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		for _, p := range currentPets {
			if p.Name == req.Name {
				return echo.NewHTTPError(http.StatusBadRequest, "A pet with that name already exists.")
			}
		}

		pet := &types.Pet{
			UserID: user.ID,
			Name:   req.Name,
			Type:   req.Type,
			Tags:   req.Tags,
			DOB:    req.DOB,
		}

		if err := h.PetStore.Create(pet); err != nil {
			h.Logger.Error("error creating pet", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusCreated, pet)
	}
}

type UpdatePetRequest struct {
	ID    uuid.UUID      `json:"id" validate:"required"`
	Type  *types.PetType `json:"type" validate:"max=16"`
	Name  string         `json:"name"`
	DOB   *time.Time     `json:"dob"`
	Tags  types.PetTags  `json:"tags"`
	Blurb *string        `json:"blurb"`
}

func (h PetsHandler) UpdatePet() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := clerkUser(c)
		if !user.Authenticated {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		var req UpdatePetRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		pet := &types.Pet{
			ID:    req.ID,
			Type:  req.Type,
			Name:  req.Name,
			Tags:  req.Tags,
			DOB:   req.DOB,
			Blurb: req.Blurb,
		}

		if err := h.PetStore.Update(pet, user.ID); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, pet)
	}
}

type NewTagRequest struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

func (h PetsHandler) AddTag() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := clerkUser(c)
		if !user.Authenticated {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		petID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad identifier")
		}

		var req NewTagRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		if err := c.Validate(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		pet, err := h.PetStore.Pet(petID)
		if err != nil {
			if notFound := errors.As(err, &store.ErrPetNotFound); notFound {
				return echo.NewHTTPError(http.StatusNotFound, err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		if pet.Tags == nil {
			pet.Tags = make(types.PetTags)
		}
		pet.Tags[req.Key] = req.Value
		if err := h.PetStore.Update(&pet, user.ID); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, pet)
	}
}

func (h PetsHandler) DeleteTag() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := clerkUser(c)
		if !user.Authenticated {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		petID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad identifier")
		}
		tagKey := c.Param("key")

		pet, err := h.PetStore.Pet(petID)
		if err != nil {
			if notFound := errors.As(err, &store.ErrPetNotFound); notFound {
				return echo.NewHTTPError(http.StatusNotFound, err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		if pet.Tags == nil {
			pet.Tags = make(types.PetTags)
		}
		delete(pet.Tags, tagKey)
		if err := h.PetStore.Update(&pet, user.ID); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, pet)
	}
}

func (h PetsHandler) DeletePet() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := clerkUser(c)
		if !user.Authenticated {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		petID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad identifier")
		}

		if err := h.PetStore.Delete(petID, user.ID); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusNoContent)
	}
}

func (h PetsHandler) UpdateImage() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := clerkUser(c)
		if !user.Authenticated {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		petID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad identifier")
		}

		avatar, err := c.FormFile("file")
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		userDir := "./static/usr/" + user.ID
		if err := os.MkdirAll(userDir, os.ModePerm); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		src, err := avatar.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		defer src.Close()
		dst, err := os.Create(userDir + "/pet_avatar_" + petID.String() + ".jpg")
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		avatarURI := userDir[2:] + "/pet_avatar_" + petID.String() + ".jpg"
		if err := h.PetStore.UpdateAvatar(avatarURI, petID, user.ID); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		return c.JSON(http.StatusCreated, map[string]string{
			"avatar_uri": avatarURI,
		})
	}
}

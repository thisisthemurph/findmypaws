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
	"paws/pkg/blight"
	"time"
)

func NewPetsHandler(s *store.PostgresStore, logger *slog.Logger) PetsHandler {
	b, err := blight.New("./avatars.db")
	if err != nil {
		panic(err)
	}

	return PetsHandler{
		AlertStore: s.AlertStore,
		PetStore:   s.PetStore,
		Blight:     b,
		Logger:     logger,
	}
}

type PetsHandler struct {
	AlertStore store.AlertStore
	PetStore   store.PetStore
	Blight     *blight.Client
	Logger     *slog.Logger
}

func (h PetsHandler) MakeRoutes(g *echo.Group) {
	g.GET("/pets/:id", h.GetPetByID())
	g.GET("/pets", h.ListPets())
	g.POST("/pets", h.CreateNewPet())
	g.PUT("/pets/:id", h.UpdatePet())
	g.POST("/pets/:id/tag", h.AddTag())
	g.DELETE("/pets/:id/tag/:key", h.DeleteTag())
	g.DELETE("/pets/:id", h.DeletePet())
	g.PUT("/pets/:id/avatar", h.UpdateImageBlight())
	g.GET("/pets/:id/avatar", h.GetAvatar())
	g.GET("/pets/test", h.Test())
	g.POST("/pets/:id/alert", h.NewAlert())
}

func (h PetsHandler) Test() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "this is only a test",
		})
	}
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
	Type  *types.PetType `json:"type" validate:"required,max=16"`
	Name  string         `json:"name" validate:"required"`
	DOB   *time.Time     `json:"dob"`
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
		if err := c.Validate(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		petId, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad identifier")
		}

		pet, err := h.PetStore.Pet(petId)
		if err != nil {
			if notFound := errors.As(err, &store.ErrPetNotFound); notFound {
				return echo.NewHTTPError(http.StatusNotFound, "The pet could not be found")
			}
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		if pet.UserID != user.ID {
			return echo.NewHTTPError(http.StatusUnauthorized, "You do not have permission to update this pet")
		}

		pet.Name = req.Name
		pet.Type = req.Type
		pet.Blurb = req.Blurb
		pet.DOB = req.DOB

		if err := h.PetStore.Update(&pet, user.ID); err != nil {
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

		if pet.UserID != user.ID {
			return echo.NewHTTPError(http.StatusUnauthorized, "You do not have permission to update this pet")
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

		if pet.UserID != user.ID {
			return echo.NewHTTPError(http.StatusUnauthorized, "You do not have permission to update this pet")
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

		existingPet, err := h.PetStore.Pet(petID)
		if err != nil {
			if notFound := errors.As(err, &store.ErrPetNotFound); notFound {
				return echo.NewHTTPError(http.StatusNotFound)
			}
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		if existingPet.UserID != user.ID {
			return echo.NewHTTPError(http.StatusUnauthorized, "You do not have permission to update this pet")
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

func (h PetsHandler) UpdateImageBlight() echo.HandlerFunc {
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

		src, err := avatar.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		defer src.Close()

		path := petID.String()
		if err := h.Blight.Add(path, src); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		return c.JSON(http.StatusCreated, map[string]string{
			"avatar_uri": path,
		})
	}
}

func (h PetsHandler) GetAvatar() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad identifier")
		}

		r, err := h.Blight.Get(id.String())
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		c.Response().Header().Set(echo.HeaderContentType, "image/jpeg")
		if _, err := io.Copy(c.Response(), r.BLOB); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return nil
	}
}

type NewAlertRequest struct {
	UserId          string `json:"user_id"`
	AnonymousUserId string `json:"anonymous_user_id"`
}

func (h PetsHandler) NewAlert() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req NewAlertRequest
		if err := c.Bind(&req); err != nil {
			h.Logger.Error("error parsing request", "error", err)
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		if req.AnonymousUserId == "" && req.UserId == "" {
			h.Logger.Error("one of the user ids is required", "req", req)
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		petId, err := uuid.Parse(c.Param("id"))
		if err != nil {
			h.Logger.Error("error parsing pet id", "error", err)
			return echo.NewHTTPError(http.StatusBadRequest, "bad identifier")
		}

		alert := types.Alert{
			UserId:          req.UserId,
			AnonymousUserId: req.AnonymousUserId,
			PetId:           petId,
		}

		if err := h.AlertStore.Create(alert); err != nil {
			h.Logger.Error("error creating", "error", err)
			if errors.Is(err, store.ErrAlertAlreadyExists) {
				return c.JSON(http.StatusOK, map[string]interface{}{
					"alert_created": false,
				})
			}
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		h.Logger.Info("alert", "petId", petId, "req", req, "alert", alert)

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"alert_created": true,
		})
	}
}

package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/sonalys/file-manager/manager/controller"
)

type Handler struct {
	*controller.Service
}

func NewHandler(s *controller.Service) *Handler {
	return &Handler{
		s,
	}
}

func (h *Handler) Start() {
	e := echo.New()
	e.HideBanner = true

	e.Use(Logger(h.Logger))
	e.POST("/upload", h.FileHandler)
	e.Static("/download", "storage")

	h.Logger.Fatal(e.Start(":8000"))
}

func (h *Handler) FileHandler(ctx echo.Context) error {
	path := ctx.FormValue("destination")
	file, err := ctx.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err := h.ReceiveFile(src, file.Filename, path); err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}

package handler

import (
	"fmt"
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
	for mountName, resolve := range h.Mounts {
		e.Static(fmt.Sprintf("/download/%s", mountName), resolve)
	}

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

	createdFiles, err := h.ReceiveFile(src, file.Filename, path)
	if err != nil {
		h.Logger.Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, createdFiles)
}

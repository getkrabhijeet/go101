package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/getkrabhijeet/go101/password-sharing/config"
	pserror "github.com/getkrabhijeet/go101/password-sharing/error"
	"github.com/getkrabhijeet/go101/password-sharing/model"
	"github.com/getkrabhijeet/go101/password-sharing/service"
)

type createLinkController struct {
	service service.PasswordService
	config  *config.Config
}

func NewCreateLinkController(service service.PasswordService, config *config.Config) Controller {
	return &createLinkController{
		service: service,
		config:  config,
	}
}

func (ctrl *createLinkController) Hander() gin.HandlerFunc {
	type Body struct {
		Password string `json:"password"`
	}

	return func(c *gin.Context) {
		body := &Body{}
		err := c.BindJSON(body)

		if err != nil {
			c.JSON(pserror.BadRequestError())
			return
		}

		link, err := ctrl.service.CreateLinkFromPassword(c, body.Password)
		if err != nil {
			psError := pserror.AsPasswordSharingError(err)
			c.JSON(psError.ToResponse())
			return
		}

		url := fmt.Sprintf("%s/%s",
			ctrl.config.App.BasePath,
			link)
		c.JSON(http.StatusCreated, model.LinkResponse{
			Url: url,
		})
	}
}

func (ctrl *createLinkController) Route() string {
	return "/link"
}

func (ctrl *createLinkController) Method() string {
	return http.MethodPost
}

package handler

import (
	_ "genesis/docs"
	"genesis/pkg/api"
	"genesis/pkg/service"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"log"
	"math"
	"strconv"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

// GetCurrentBTCUAH @Summary Отримати поточний курс BTC до UAH
// @Description Запит має повертати поточний курс BTC до UAH використовуючи будь-який third party сервіс з публічним API
// @Tags rate
// @Produce json
// @Success 200 {object} integer
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /rate [get]
func (h *Handler) GetCurrentBTCUAH(c *fiber.Ctx) error {
	a := api.API{}
	currentBTCUAH, err := a.GetCurrentBTCUAH()
	if err != nil {
		log.Print(err)
		return fiber.NewError(fiber.StatusBadRequest, "Error getting the current BTC rate")
	}

	err = c.JSON(math.Round(currentBTCUAH))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

// SendEmails @Summary Відправити e-mail з поточним курсом на всі підписані електронні пошти
// @Description Запит має отримувати актуальний курс BTC до UAH за допомогою third-party сервісу та відправляти його на всі електронні адреси, які були підписані раніше
// @Tags subscription
// @Produce json
// @Success 200 {object} string "E-mailʼи відправлено"
// @Failure 500 {object} string "Помилка на сервері"
// @Router /sendEmails [post]
func (h *Handler) SendEmails(c *fiber.Ctx) error {
	toEmails, err := h.service.GetEmails()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	a := api.API{}
	result, err := a.GetRate()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	rate := strconv.Itoa(result)
	invalidEmails := a.SendMail(toEmails, rate)
	if invalidEmails != nil {
		//return fiber.NewError()
		err := c.Status(fiber.StatusOK).JSON(map[string]interface{}{
			"message":       "Emails have been sent, but some addresses are invalid",
			"invalidEmails": invalidEmails,
		})
		if err != nil {
			return err
		}
		return err
	}
	err = c.JSON("Emails sent")
	if err != nil {
		return err
	}
	return nil

}

// Subscribe @Summary Підписати емейл на отримання поточного курсу
// @Description Запит має перевірити, чи немає данної електронної адреси в поточній базі даних (файловій) і, в разі її відсутності, записувати її. Пізніше, за допомогою іншого запиту ми будемо відправляти лист на ті електронні адреси, які будуть в цій базі.
// @Tags subscription
// @Accept x-www-form-urlencoded
// @Produce json
// @Param email formData string true "Електронна адреса, яку потрібно підписати"
// @Success 200 {object} string "Email додано"
// @Failure 409 {object} string "Email вже є в базі даних"
// @Failure 500 {object} string "Помилка на сервері"
// @Router /subscribe [post]
func (h *Handler) Subscribe(c *fiber.Ctx) error {
	email := c.FormValue("email")
	if email == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing email")
	}
	err := h.service.AddEmail(email)
	if err != nil {
		if err.Error() == "email already subscribed" {
			return fiber.NewError(fiber.StatusConflict, err.Error())
		}
		return err
	}
	err = c.JSON("Email added successfully")
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) InitRoutes() *fiber.App {
	app := fiber.New()

	app.Use(recover.New())
	app.Use(cors.New())
	apiRoutes := app.Group("/api")
	{
		apiRoutes.Get("/rate", h.GetCurrentBTCUAH)
		apiRoutes.Post("/subscribe", h.Subscribe)
		apiRoutes.Post("/sendEmails", h.SendEmails)
	}
	app.Get("/swagger/*", swagger.HandlerDefault)
	return app
}

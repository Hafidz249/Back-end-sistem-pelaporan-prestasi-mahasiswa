package routes

import (
		"POJECT_UAS/middleware"
	. 	"POJECT_UAS/middleware"
		"POJECT_UAS/model"
		"POJECT_UAS/service"

	"github.com/gofiber/fiber/v2"
)

func PekerjaanAlumni(app *fiber.App, userRepo *model.UserRepository) {
	app.Get("/pekerjaan", JWTAuth(userRepo), RequireRole("admin", "user"), service.GetAllpekerjaanAlumniService)
	app.Get("/pekerjaan/:id", JWTAuth(userRepo), RequireRole("admin", "user"), service.CheckpekerjaanAlumniService)
	app.Get("/pekerjaan/alumni/:nim_alumni", JWTAuth(userRepo), RequireRole("admin"), service.CheckpekerjaanAlumniService)
	app.Post("/pekerjaan", JWTAuth(userRepo), RequireRole("admin"), service.CreatepekerjaanAlumniService)
	app.Put("/softdeleted/:id", middleware.JWTAuth(userRepo), service.SoftDeleteBynimService)
	app.Get("/trash", middleware.JWTAuth(userRepo), service.GetAllTrashService)
	app.Put("/restore/:id", middleware.JWTAuth(userRepo), service.RestoreBynimService)
	app.Delete("/deleted/:id",middleware.JWTAuth(userRepo), service.DeletePekerjaanAlumniService)
}
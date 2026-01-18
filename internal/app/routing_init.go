package app

import (
	"github.com/go-chi/chi"
)

func RegisterRoutes(r chi.Router, app *AppInstance) {
	r.Route("/chats", func(r chi.Router) {

		//Хэндлеры чата
		r.Post("/", app.API.ChatAPI.CreateChat)
		r.Post("/{id}", app.API.ChatAPI.GetChat)
		r.Delete("/{id}", app.API.ChatAPI.DeleteChat)
		r.Post("/{id}/messages/", app.API.ChatAPI.SendMessage)

	})
}

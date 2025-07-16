package main

import (
	"snap-rq/app/controller"
	logger "snap-rq/app/log"
	"snap-rq/app/service"
	"snap-rq/app/view"

	"github.com/gdamore/tcell/v2"
)

func main() {
	logger.Init("app.log")

	// Set up services: should not perform any initialisation logic that would affect any views.
	// Services are not hierarchical, they talk between each other, but don't need parent-child relationship
	var services = service.NewAppService()

	// Load app
	var app = view.NewAppView()
	// Init root app controller
	var controller = controller.NewAppController(app, services)

	app.Views.CollectionsList.SetListener(&controller)
	app.Views.RequestsList.SetListener(&controller)
	app.Views.MethodPickerModal.SetListener(&controller)
	app.Views.UrlInputView.SetListener(&controller)

	app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		return false // Allow normal drawing to continue
	})

	controller.Start()
	app.Init()
}

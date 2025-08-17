package main

import (
	"snap-rq/app/controller"
	logger "snap-rq/app/log"
	"snap-rq/app/service"
	"snap-rq/app/view"
)

func main() {
	logger.Init("app.log")

	var service = service.NewAppService()
	var controller = controller.NewAppController(service)
	var app = view.NewAppView(&controller)

	// The listener order is important
	app.Views.CollectionsList.SetListener(&controller)
	app.Views.RequestsList.SetListener(&controller)
	app.Views.RequestHeaderBar.SetListener(&controller)
	app.Views.EditorView.SetListener(&controller)
	app.Views.NameEditorModal.SetListener(&controller)
	app.Init()
	service.Start()
	controller.Start(&app)

	app.Start()
}

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

	var app = view.NewAppView()

	var controller = controller.NewAppController(app, service)

	app.SetAppViewListener(&controller)
	app.Views.CollectionsList.SetListener(&controller)
	app.Views.RequestsList.SetListener(&controller)
	app.Views.RequestHeaderBar.SetListener(&controller)
	app.Views.EditorView.SetListener(&controller)

	app.Init()
	service.Start()
	controller.Start()

	app.Start() // no function can run beyond this point due to UI loop start
}

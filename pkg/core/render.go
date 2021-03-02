package core

type IRender interface {
	RenderAll(task *Task, event []*Event) (string, error)
	RenderTask(task *Task) (string, error)
	RenderEvents(event []*Event) (string, error)
	RenderDesc(desc string) (string, error)
}

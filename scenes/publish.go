package scenes

import "fmt"

type PublishScenes struct {

}

func (ps PublishScenes) New() interface{} {
	return &PublishScenes{}
}

func (ps PublishScenes) Run() {
	fmt.Println("PublishScenes")
}
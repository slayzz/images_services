package main

import (
	"github.com/slayzz/images_services/pkg/imager/externalserv"
	"log"
)

func main() {
	imageRequester := externalserv.NewImageRequesterUnsplash()
	image, err := imageRequester.GetImage()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(image))
	//tasks := []*worker.Task{
	//	worker.NewTask(func() (interface {}, error) { time.Sleep(time.Second); return nil }),
	//}
	//
	//p := worker.NewPool(tasks, 8)
	//p.Run()
	//
	//var numErrors int
	//for _, task := range p.Tasks {
	//	if task.Err != nil {
	//		log.Println(task.Err)
	//		numErrors++
	//	}
	//	if numErrors >= 10 {
	//		log.Println("Too many errors.")
	//		break
	//	}
	//}

}

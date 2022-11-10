package main

import (
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"log"
	"net/http"
	"sync"
	"time"
)

func Init() {
	sentinelInit()
	createFlowRule("limit", 1, 2000)
}

func main01() {
	Init()

	http.HandleFunc("/flow", func(writer http.ResponseWriter, request *http.Request) {
		resName := request.FormValue("name")
		fmt.Println(resName)
		e, b := sentinel.Entry(resName, sentinel.WithTrafficType(base.Inbound))
		if b != nil {
			fmt.Fprint(writer, "限流！！！")
		} else {
			fmt.Fprint(writer, "正常访问...")
			e.Exit()
		}
	})

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func sentinelInit() {
	conf := config.NewDefaultConfig()
	err := sentinel.InitWithConfig(conf)
	if err != nil {
		log.Fatal(err)
	}
}

func createFlowRule(resourceName string, threshold float64, interval uint32) {
	_, err := flow.LoadRules([]*flow.Rule{
		{
			Resource:               resourceName,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              threshold,
			StatIntervalInMs:       interval,
		},
	})
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
		return
	}
}

var wg sync.WaitGroup

func main02() {
	baton := make(chan int)

	wg.Add(1)
	go Runner(baton)

	baton <- 1
	wg.Wait()
}

func Runner(baton chan int) {
	var newRunner int

	runner := <-baton

	fmt.Printf("Runner %d Running with Baton\n", runner)

	if runner != 4 {
		newRunner = runner + 1
		fmt.Printf("Runner %d to the line\n", newRunner)
		go Runner(baton)
	}

	time.Sleep(time.Second)

	if runner == 4 {
		fmt.Printf("Runner %d Finished, Race over\n", runner)
		wg.Done()
		return
	}

	fmt.Printf("Runner %d Exchange with Runner %d\n", runner, newRunner)

	baton <- newRunner
}

func main() {
	fmt.Println(len("f01ed02817c3fcf358e84bb153dc0b8e"))
}

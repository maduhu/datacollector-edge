package main

import (
	"flag"
	"fmt"
	"github.com/streamsets/datacollector-edge/container/common"
	"github.com/streamsets/datacollector-edge/container/edge"
	_ "github.com/streamsets/datacollector-edge/stages/destinations"
	_ "github.com/streamsets/datacollector-edge/stages/origins"
	_ "github.com/streamsets/datacollector-edge/stages/processors"
	"log"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strings"
	"syscall"
)

func main() {
	debugFlag := flag.Bool("debug", false, "Debug flag")
	logToConsoleFlag := flag.Bool("logToConsole", false, "Log to console flag")
	startFlag := flag.String("start", "", "Start Pipeline flag")
	runtimeParametersFlag := flag.String("runtimeParameters", "", "Runtime Parameters flag")
	flag.Parse()

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	baseDir := strings.TrimSuffix(path.Dir(ex), "/bin")

	fmt.Println("StreamSets Data Collector Edge (SDC Edge): ")
	fmt.Printf("OS: %s\nArchitecture: %s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println("Base Dir: ", baseDir)

	dataCollectorEdge, _ := edge.DoMain(baseDir, *debugFlag, *logToConsoleFlag, *startFlag, *runtimeParametersFlag)
	go shutdownHook(dataCollectorEdge)
	dataCollectorEdge.WebServerTask.Run()
}

func shutdownHook(dataCollectorEdge *edge.DataCollectorEdgeMain) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("[INFO] Program got a system signal %v\n", <-c)
	if pipelineInfos, er := dataCollectorEdge.PipelineStoreTask.GetPipelines(); er == nil {
		for _, pipelineInfo := range pipelineInfos {
			runner := dataCollectorEdge.Manager.GetRunner(pipelineInfo.PipelineId)
			if pipelineState, er := runner.GetStatus(); er == nil &&
				(pipelineState.Status == common.RUNNING || pipelineState.Status == common.STARTING) {
				log.Printf("[INFO] Stopping pipeline : %s\n", pipelineInfo.PipelineId)
				if runner.StopPipeline(); er != nil {
					log.Printf("[INFO] Error happened when stopping pipeline : %s\n", pipelineInfo.PipelineId)
				}
			}
		}
	}
	dataCollectorEdge.WebServerTask.Shutdown()
	if dataCollectorEdge.RuntimeInfo.DPMEnabled {
		dataCollectorEdge.DPMMessageEventHandler.Shutdown()
	}
	log.Println("[INFO] Data Collector Edge shutting down")
}

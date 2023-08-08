package daemon

import (
	"github.com/RaymondCode/simple-demo/amqp"
	"github.com/RaymondCode/simple-demo/api"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var outDeviceUrlMap = map[int][2]string{
	1: {"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did106h_7_0100452c2f257624_0000",
		"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did104h_7_0100452c2f257624_0000"},
	2: {"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did106h_7_0100452c2f257624_0001",
		"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did104h_7_0100452c2f257624_0001"},
}

var innerDeviceUrlMap = map[int][2]string{
	1: {"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did102h_9_0100452c2f257624_0001",
		"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did103h_9_0100452c2f257624_0001"},
	2: {"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did102h_9_0100452c2f257624_0002",
		"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did103h_9_0100452c2f257624_0002"},
	3: {"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did102h_5_0100452c2f257624_0003",
		"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did103h_5_0100452c2f257624_0003"},
	4: {"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did102h_9_0100452c2f257624_0004",
		"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did103h_9_0100452c2f257624_0004"},
	5: {"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did102h_69_0100452c2f257624_0005",
		"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did103h_69_0100452c2f257624_0005"},
	6: {"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did102h_9_0100452c2f257624_0006",
		"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did103h_9_0100452c2f257624_0006"},
	7: {"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did102h_9_0100452c2f257624_0007",
		"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did103h_9_0100452c2f257624_0007"},
	8: {"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did102h_9_0100452c2f257624_0008",
		"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did103h_9_0100452c2f257624_0008"},
	9: {"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did102h_5_0100452c2f257624_0009",
		"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did103h_5_0100452c2f257624_0009"},
	10: {"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did102h_9_0100452c2f257624_0010",
		"https://cloudmaster.hisensehitachi.com/data/tdengine/table/did103h_9_0100452c2f257624_0010"},
}

func InitDaemon() {
	// Create a channel for the queue
	outdoorQueue := make(chan *api.OutdoorDevice)
	indoorQueue := make(chan *api.IndoorDevice)

	// Create a wait group to ensure all goroutines finish before exiting the program
	var wg sync.WaitGroup

	// Start the scheduler to produce data
	wg.Add(1)
	go scheduler(&wg, outdoorQueue, indoorQueue)

	// Start the consumer to handle data in batches
	wg.Add(1)
	go consumer(&wg, outdoorQueue, indoorQueue)

	go amqp.HandlerMessage()

	// Wait for all goroutines to finish their jobs
	wg.Wait()
}

func scheduler(wg *sync.WaitGroup, outdoorQueue chan<- *api.OutdoorDevice, indoorQueue chan *api.IndoorDevice) {
	// Stop the scheduler when the function finishes
	defer wg.Done()

	// Create a ticker with 1-second interval
	ticker := time.NewTicker(1 * time.Second)

	// Run the scheduler until it's stopped
	for {
		select {
		case <-ticker.C:
			token := api.GetGlobalToken()
			// Call Getapi.OutdoorDeviceInfo and store the result in the queue
			for deviceId, urlArr := range outDeviceUrlMap {
				device, err := api.GetOutdoorDeviceInfo(urlArr[0], urlArr[1], token)
				if err != nil {
					log.Errorf("GetOutdoorDeviceInfo Error: %v", err)
				} else {
					device.DeviceId = deviceId
					select {
					case outdoorQueue <- device:
					default:
						log.Info("Queue is full. Skipping element.")
					}
				}
			}
			for deviceId, urlArr := range innerDeviceUrlMap {
				device, err := api.GetIndoorDeviceInfo(urlArr[0], urlArr[1], token)
				if err != nil {
					log.Errorf("GetOutdoorDeviceInfo Error: %v", err)
				} else {
					device.DeviceId = deviceId
					select {
					case indoorQueue <- device:
					default:
						log.Info("Queue is full. Skipping element.")
					}
				}
			}
		}
	}
}

func consumer(wg *sync.WaitGroup, outdoorQueue <-chan *api.OutdoorDevice, indoorQueue <-chan *api.IndoorDevice) {
	// Stop the consumer when the function finishes
	defer wg.Done()

	// Create a list to store the elements
	var deviceList []interface{}

	for {
		select {
		case device := <-outdoorQueue:
			deviceList = append(deviceList, device)
		case device := <-indoorQueue:
			deviceList = append(deviceList, device)
		}
		// If the list reaches 30 elements, process it
		if len(deviceList) == 30 {
			// Todo
			// Process the list as needed (in this example, we print the elements)
			log.Info("Received batch of 30 elements:", deviceList)
			amqp.SendMessage(deviceList)
			// Clear the list for the next batch
			deviceList = nil
		}
	}

	// Process the remaining elements (if any)
	if len(deviceList) > 0 {
		log.Info("Received final batch of elements:", deviceList)
	}
}

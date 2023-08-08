package daemon

import (
	"github.com/RaymondCode/simple-demo/amqp"
	"github.com/RaymondCode/simple-demo/api"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

func InitDaemon() {
	// Create a channel for the queue
	queue := make(chan *api.OutdoorDevice)

	// Create a wait group to ensure all goroutines finish before exiting the program
	var wg sync.WaitGroup

	// Start the scheduler to produce data
	wg.Add(1)
	go scheduler(&wg, queue)

	// Start the consumer to handle data in batches
	wg.Add(1)
	go consumer(&wg, queue)

	// Wait for all goroutines to finish their jobs
	wg.Wait()
}

func scheduler(wg *sync.WaitGroup, queue chan<- *api.OutdoorDevice) {
	// Stop the scheduler when the function finishes
	defer wg.Done()

	// Create a ticker with 1-second interval
	ticker := time.NewTicker(1 * time.Second)

	// Run the scheduler until it's stopped
	for {
		select {
		case <-ticker.C:
			// Call Getapi.OutdoorDeviceInfo and store the result in the queue
			device, err := api.GetOutdoorDeviceInfo("https://cloudmaster.hisensehitachi.com/auth/oauth/token?username=p_dcfyzd_shijingfeng&password=W%2Bx6Ljdj7ZlLz6wDkpju3w%3D%3D&grant_type=client_credentials&scope=server", "ff85e547-c47e-45fb-a9cf-e165870a8e50")
			if err != nil {
				log.Errorf("GetOutdoorDeviceInfo Error: %v", err)
			} else {
				select {
				case queue <- device:
				default:
					log.Info("Queue is full. Skipping element.")
				}
			}
		}
	}
}

func consumer(wg *sync.WaitGroup, queue <-chan *api.OutdoorDevice) {
	// Stop the consumer when the function finishes
	defer wg.Done()

	// Create a list to store the elements
	var deviceList []*api.OutdoorDevice

	for device := range queue {
		deviceList = append(deviceList, device)

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

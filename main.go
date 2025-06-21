package main

import (
	"fmt"
	"io"
	"log"
	gtfsnyct "mta-realtime"
	"net/http"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"

	gtfsnyct "proto/gtfsnyct"
)

func main() {
	feedURL := "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs"

	resp, err := http.Get(feedURL)
	if err != nil {
		log.Fatalf("Failed to fetch GTFS-R feed: %v", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	feed := &gtfs.FeedMessage{}
	if err := proto.Unmarshal(raw, feed); err != nil {
		log.Fatalf("Failed to unmarshal feed: %v", err)
	}

	fmt.Println("Received feed entities:", len(feed.Entity))

	for _, entity := range feed.Entity {
		if entity.TripUpdate == nil {
			continue
		}

		trip := entity.TripUpdate.Trip
		fmt.Printf("Trip ID: %s | Route ID: %s\n", trip.TripId, trip.RouteId)

		if ext, err := proto.GetExtension(trip, gtfsnyct.E_NyctTripDescriptor); err == nil {
			if nyctTrip, ok := ext.(*gtfsnyct.NyctTripDescriptor); ok {
				fmt.Printf("  NYCT Train ID: %s | Assigned: %t | Direction: %s\n", nyctTrip.TrainId, nyctTrip.IsAssigned, nyctTrip.Direction)
			}
		}

		for _, stu := range entity.TripUpdate.StopTimeUpdate {
			fmt.Println("  Stop ID:", stu.StopId)

			if ext, err := proto.GetExtension(stu, gtfsnyct.E_NyctStopTimeUpdate); err == nil {
				if nyctStop, ok := ext.(*gtfsnyct.NyctStopTimeUpdate); ok {
					fmt.Printf("    Scheduled Track: %s | Actual Track: %s\n", nyctStop.ScheduledTrack, nyctStop.ActualTrack)
				}
			}
		}

		fmt.Println()
	}

	// for _, entity := range feed.Entity {
	// 	if entity.TripUpdate != nil {
	// 		// Already included in the block below
	// 		for _, stopTimeUpdate := range entity.TripUpdate.StopTimeUpdate {
	// 			fmt.Printf("  Stop ID: %s\n", stopTimeUpdate.StopId)
	// 			if stopTimeUpdate.Arrival != nil {
	// 				fmt.Printf("    Arrival Time: %d\n", stopTimeUpdate.Arrival.Time)
	// 				fmt.Printf("    Arrival Delay: %d\n", stopTimeUpdate.Arrival.Delay)
	// 				fmt.Printf("    Arrival Uncertainty: %d\n", stopTimeUpdate.Arrival.Uncertainty)
	// 			}
	// 			if stopTimeUpdate.Departure != nil {
	// 				fmt.Printf("    Departure Time: %d\n", stopTimeUpdate.Departure.Time)
	// 				fmt.Printf("    Departure Delay: %d\n", stopTimeUpdate.Departure.Delay)
	// 				fmt.Printf("    Departure Uncertainty: %d\n", stopTimeUpdate.Departure.Uncertainty)
	// 			}
	// 			fmt.Printf("    Schedule Relationship: %s\n", stopTimeUpdate.ScheduleRelationship.String())
	// 		}

	// 		trip := entity.TripUpdate.Trip
	// 		fmt.Printf("Trip ID: %s\n", trip.TripId)
	// 		fmt.Printf("Route ID: %s\n", trip.RouteId)
	// 		fmt.Printf("Start Date: %s\n", trip.StartDate)
	// 		fmt.Printf("Start Time: %s\n", trip.StartTime)
	// 		fmt.Printf("Schedule Relationship: %s\n", trip.ScheduleRelationship.String())

	// 		vehicle := entity.TripUpdate.Vehicle
	// 		if vehicle != nil {
	// 			fmt.Printf("Vehicle ID: %s\n", vehicle.Id)
	// 			fmt.Printf("Label: %s\n", vehicle.Label)
	// 			fmt.Printf("License Plate: %s\n", vehicle.LicensePlate)
	// 		}
	// 	}
	// 	break
	// }
}

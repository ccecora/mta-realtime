#!/usr/bin/env python3
import sys
import requests
from google.protobuf import message
import com.google.transit.realtime.gtfs_realtime_pb2 as gtfs
import gtfs_realtime_NYCT_pb2 as nyct

FEED_URL = "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs"

def fetch_feed(url: str) -> gtfs.FeedMessage:
    r = requests.get(url)
    r.raise_for_status()
    feed = gtfs.FeedMessage()
    feed.ParseFromString(r.content)
    return feed

def print_feed(feed: gtfs.FeedMessage):
    print(f"Received feed entities: {len(feed.entity)}")
    for entity in feed.entity:
        if not entity.HasField("trip_update"):
            continue
        tu = entity.trip_update

        # ——— TripDescriptor ———
        trip = tu.trip
        print(f"Trip ID: {trip.trip_id} | Route ID: {trip.route_id}")
        print(f"  Start Date: {trip.start_date} | Start Time: {trip.start_time}")
        print(f"  Schedule Relationship: {gtfs.TripDescriptor.ScheduleRelationship.Name(trip.schedule_relationship)}")

        # NYCT extension on TripDescriptor
        if trip.HasExtension(nyct.nyct_trip_descriptor):
            ext = trip.Extensions[nyct.nyct_trip_descriptor]
            print(f"  ➜ NYCT Train ID: {ext.train_id}")
            print(f"    Assigned: {ext.is_assigned} | Direction: {ext.direction} | Line: {ext.line}")

        # ——— StopTimeUpdates ———
        for stu in tu.stop_time_update:
            print(f"  Stop ID: {stu.stop_id}")
            if stu.HasField("arrival"):
                arr = stu.arrival
                print(f"    Arrival → Time: {arr.time} | Delay: {arr.delay} | Uncertainty: {arr.uncertainty}")
            if stu.HasField("departure"):
                dep = stu.departure
                print(f"    Departure → Time: {dep.time} | Delay: {dep.delay} | Uncertainty: {dep.uncertainty}")
            print(f"    Schedule Rel: {gtfs.TripDescriptor.ScheduleRelationship.Name(stu.schedule_relationship)}")

            # NYCT extension on StopTimeUpdate
            if stu.HasExtension(nyct.nyct_stop_time_update):
                ext2 = stu.Extensions[nyct.nyct_stop_time_update]
                print(f"    Scheduled Track: {ext2.scheduled_track} | Actual Track: {ext2.actual_track}")

        # ——— VehicleDescriptor ———
        if tu.HasField("vehicle"):
            v = tu.vehicle
            print(f"  Vehicle ID: {v.id} | Label: {v.label} | License Plate: {v.license_plate}")

        print()

def main():
    try:
        feed = fetch_feed(FEED_URL)
        print_feed(feed)
    except (requests.HTTPError, message.DecodeError) as e:
        print("Error:", e, file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()

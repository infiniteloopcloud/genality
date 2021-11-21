# genality

General analytics solution provided as a library or service.


### Usage

```go

package main

import (
	"context"
	"log"
	"time"

	"github.com/infiniteloopcloud/genality"
)

func main() {
	var connectionString = "postgres://genality:genality@localhost:5435/genality?sslmode=disable"
	var ctx = context.Background()
	// create new instance
	instance, _ := genality.New(genality.Opts{ConnectionString: connectionString})
	// handle err

	instance.Add(ctx, "random_record")

	// get apikey useage in a time period
	cnt, _ := instance.GetCountFrom(ctx, "random_record", time.Now().Add(-24*time.Hour))
	// handle err
	log.Print(cnt)

	// get certain buckets, where every bucket stores the corresponding count(apikey)
	res, _ := instance.GetCountBuckets(ctx, "random_record", time.Now().Add(-24*time.Hour), time.Hour)
	// handle err
	log.Print(res)
	// returns something like this:
    // [
    //      { bucket: 2021-11-19 04:00:00 +0000 UTC, count: 2 }
    //      { bucket: 2021-11-19 05:00:00 +0000 UTC, count: 4 }
    //      { bucket: 2021-11-19 06:00:00 +0000 UTC, count: 4 }
    //      { bucket: 2021-11-19 07:00:00 +0000 UTC, count: 4 }
    //      { bucket: 2021-11-19 08:00:00 +0000 UTC, count: 4 }
    //      { bucket: 2021-11-19 09:00:00 +0000 UTC, count: 4 }
    //      { bucket: 2021-11-19 10:00:00 +0000 UTC, count: 4 }
    //      { bucket: 2021-11-19 11:00:00 +0000 UTC, count: 4 }
    //      { bucket: 2021-11-19 12:00:00 +0000 UTC, count: 4 }
    //      { bucket: 2021-11-19 13:00:00 +0000 UTC, count: 4 }
    //      { bucket: 2021-11-19 14:00:00 +0000 UTC, count: 102 }
    // ]
}
```
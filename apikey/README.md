# api_key_metrics

##Example

    func main() { 

        // create new metrics
        m := NewMetrics(databaseConnectionString)
        m.Add(context.Background(), apikey)
    
        // get apikey useage in a time period
        cnt, err := m.GetCountFrom(context.Background(), apikey, fromStartTime)
    
        // get certain buckets, where every bucket stores the corresponding count(apikey)
        res, err := m.GetCountBuckets(ctx, apikey, timeStart, bucketSize) 
        // e.g. timeStart = time.Now().Add(-1* time.Hour) and bucketSize = time.Hour
        // returns something like this:
        [
            { bucket: 2021-11-19 04:00:00 +0000 UTC, count: 2 }
            { bucket: 2021-11-19 05:00:00 +0000 UTC, count: 4 } 
            { bucket: 2021-11-19 06:00:00 +0000 UTC, count: 4 } 
            { bucket: 2021-11-19 07:00:00 +0000 UTC, count: 4 } 
            { bucket: 2021-11-19 08:00:00 +0000 UTC, count: 4 } 
            { bucket: 2021-11-19 09:00:00 +0000 UTC, count: 4 } 
            { bucket: 2021-11-19 10:00:00 +0000 UTC, count: 4 } 
            { bucket: 2021-11-19 11:00:00 +0000 UTC, count: 4 } 
            { bucket: 2021-11-19 12:00:00 +0000 UTC, count: 4 } 
            { bucket: 2021-11-19 13:00:00 +0000 UTC, count: 4 } 
            { bucket: 2021-11-19 14:00:00 +0000 UTC, count: 102 }
        ]
    
    }

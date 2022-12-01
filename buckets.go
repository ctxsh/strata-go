package apex

import "github.com/prometheus/client_golang/prometheus"

var (
	// DefBuckets are the default Histogram buckets.
	DefBuckets = prometheus.DefBuckets
)

// The bucket functions are convienience wrappers around the prometheus
// bucket functions.

// LinearBuckets creates 'count' buckets, each 'width' wide, where the lowest
// bucket has an upper bound of 'start'. The final +Inf bucket is not counted
// and not included in the returned slice.  The returned slice is meant to be
// used for the Buckets field of HistogramOpts. (from prometheus docs)
//
// The function panics if 'count' is zero or negative.
func LinearBuckets(start, width float64, count int) []float64 {
	return prometheus.LinearBuckets(start, width, count)
}

// ExponentialBuckets creates 'count' buckets, where the lowest bucket has an
// upper bound of 'start' and each following bucket's upper bound is 'factor'
// times the previous bucket's upper bound. The final +Inf bucket is not
// counted and not included in the returned slice. The returned slice is meant
// to be used for the Buckets field of HistogramOpts. (from prometheus docs)
//
// The function panics if 'count' is 0 or negative, if 'start' is 0 or negative,
// or if 'factor' is less than or equal 1.
func ExponentialBuckets(start, factor float64, count int) []float64 {
	return prometheus.ExponentialBuckets(start, factor, count)
}

// ExponentialBucketsRange creates 'count' buckets, where the lowest bucket is
// 'min' and the highest bucket is 'max'. The final +Inf bucket is not counted
// and not included in the returned slice. The returned slice is meant to be
// used for the Buckets field of HistogramOpts. (from prometheus docs)
//
// The function panics if 'count' is 0 or negative, if 'min' is 0 or negative.
func ExponentialBucketRange(min, max float64, count int) []float64 {
	return prometheus.ExponentialBucketsRange(min, max, count)
}

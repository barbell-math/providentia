package types

type (
	// A struct that is used to represent the bar path when it has been
	// calculated by an external source.
	RawTimeSeriesData struct {
		TimeData     []Second             // The time data for the set
		PositionData []Vec2[Meter, Meter] // The position data for the set
	}

	// A tagged union that either contains a [RawTimeSeriesData] struct or a
	// path to a video file.
	//
	// A zero initialized BarPathVariant will hold neither a video path or time
	// series data and can be used to represent having no data.
	//
	// Use [BarPathVariant] or [BarPathTimeSeriesData] to initialize.
	BarPathVariant struct {
		flag      BarPathFlag
		videoPath string
		posData   RawTimeSeriesData
	}
)

// Returns a [BarPathVariant] initialized with a video path as the data source.
func BarPathVideo(videoPath string) BarPathVariant {
	return BarPathVariant{
		flag:      VideoBarPathData,
		videoPath: videoPath,
	}
}

// Returns a [BarPathVariant] initialized with time series data as the data
// source.
func BarPathTimeSeriesData(data RawTimeSeriesData) BarPathVariant {
	return BarPathVariant{
		flag:    TimeSeriesBarPathData,
		posData: data,
	}
}

// Gets the underlying data source.
func (b *BarPathVariant) Source() BarPathFlag {
	return b.flag
}

// Gets the video data if it is present in the variant, otherwise returns an
// empty string. The boolean return argument indicates if the video data was
// present or not.
func (b *BarPathVariant) VideoData() (string, bool) {
	if b.flag == VideoBarPathData {
		return b.videoPath, true
	}
	return "", false
}

// Gets the time series data if it is present in the variant, otherwise returns
// an struct. The boolean return argument indicates if the time series data was
// present or not.
func (b *BarPathVariant) TimeSeriesData() (RawTimeSeriesData, bool) {
	if b.flag == TimeSeriesBarPathData {
		return b.posData, true
	}
	return RawTimeSeriesData{}, false
}

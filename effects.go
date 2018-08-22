package multirender

import (
	"sort"
)

// TODO: Also create a concurrent version of StretchConstrast
// TODO: Also create a version of StretchContrast that disregards alpha.
// TODO: Also create a version of StretchContrast that uses the uint32 colorvalue directly.

// A data structure to hold key/value pairs
type Pair struct {
	Key   uint8
	Value int
}

// A slice of pairs that implements sort.Interface to sort by values
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// StretchContrast uses "cores" CPU cores to concurrently stretch the contrast of the pixels
// in the given "pixels" slice (of width "pitch"), discarding the discardRatio ratio of the
// most unpopular pixel values, then scaling the remaining pixels to cover the full 0..255 range.
func StretchContrast(cores int, pixels []uint32, pitch int32, discardRatio float32) {

	// TODO: Use N cores

	// Find all pixel values, store them in a map[uint8]int, where the int is the count
	popularity := make(map[uint8]int)
	for i := range pixels {
		v := Value(pixels[i])
		popularity[v]++
	}

	//fmt.Println("POPULARITY", popularity)

	// How large ratio of the values should be discarded?
	lengthOfSelectedKeys := int(float32(len(popularity)) * (1.0 - discardRatio))

	// Sort the popularity map by value, by placing it in a slice of structs
	// Sort the map by the popularity of the combined value of the colors,
	// by placing it in a slice of structs that has a key and value,
	// and then sorting it with sort.Sort.
	sortablePopularity := make(PairList, len(popularity))
	i := 0
	for k, v := range popularity {
		sortablePopularity[i] = Pair{k, v}
		i++
	}
	sort.Sort(sortablePopularity)

	// Discard the least popular brightness values
	selectedKeyValues := sortablePopularity[lengthOfSelectedKeys:]

	minValue := uint8(255) // start high, reduce when smaller values are found
	maxValue := uint8(0)   // start low, increase when larger values are found

	for _, pair := range selectedKeyValues {
		pixelValue := pair.Key
		if pixelValue < minValue {
			minValue = pixelValue
		}
		if pixelValue > maxValue {
			maxValue = pixelValue
		}
	}

	lowestV := minValue
	highestV := maxValue

	widthV := highestV - lowestV

	//fmt.Println("lowestV", lowestV)
	//fmt.Println("highestV", highestV)
	//fmt.Println("widthV", widthV)

	// Clamp and scale all pixels
	var r, g, b, a uint8
	for i := range pixels {
		r, g, b, a = ColorValueToRGBA(pixels[i])

		ratioR := float32(r-lowestV) / float32(widthV)
		r = uint8(ratioR * float32(255))

		ratioG := float32(g-lowestV) / float32(widthV)
		g = uint8(ratioG * float32(255))

		ratioB := float32(b-lowestV) / float32(widthV)
		b = uint8(ratioB * float32(255))

		pixels[i] = RGBAToColorValue(r, g, b, a)
	}
}

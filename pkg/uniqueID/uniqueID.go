package uniqueID

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"time"
)

func randomNumWithTimeStamp() (roughTimeStamp int32, randomNum int32) {
	randomNum = rand.Int32()

	// Takes the highest 10 bits of the randomNum as the sleep time (in milliseconds).
	// So this function will sleep for around 1 second at most.
	// Take this a precaution for the multi-threading cases.
	randomSleepDuration := randomNum >> 22
	time.Sleep(time.Duration(randomSleepDuration) * time.Millisecond)

	// roughTimeStamp is masked by DecMask6Digit (524287),
	// to make sure the roughTimeStamp is 6 digits long.
	roughTimeStamp = int32(time.Now().UnixMilli() & DecMask6Digit)

	// The returned randomNum is masked by DecMask5Digit (65535),
	// to make sure the randomNum is 5 digits long.
	randomNum &= DecMask5Digit
	return
}

// format of the result:
//
// "roughTimeStamp-word1-word2-word3-randomNum"
//
// "123456-abc-def-ghi-12345"
func UniqueID() string {
	var result strings.Builder
	result.Grow(6 + 1 + (maxWordLength+1)*3 + 5)

	sep := byte('-')
	roughTimeStamp, randomNum := randomNumWithTimeStamp()

	result.WriteString(fmt.Sprintf("%06d", roughTimeStamp))
	result.WriteByte(sep)

	for i, randomBits := 0, rand.Int32(); i < 3; i++ {
		// Takes only the lowest 10 bits of the randomNum.
		randomPos := randomBits & 0b1111111111
		randomBits >>= 10

		result.WriteString(wordList[randomPos])
		result.WriteByte(sep)
	}

	result.WriteString(fmt.Sprintf("%05d", randomNum))
	return result.String()
}

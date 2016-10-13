package gokids

import (
    "time"
)

func UtilTimeFewDaysLater(baseTime time.Time, day int) time.Time {
    return UtilTimeFewDurationLater(baseTime, time.Duration(day) * 24 * time.Hour)
}

func UtilTimeTwentyFourHoursLater(baseTime time.Time) time.Time {
    return UtilTimeFewDurationLater(baseTime, time.Duration(24) * time.Hour)
}

func UtilTimeSixHoursLater(baseTime time.Time) time.Time {
    return UtilTimeFewDurationLater(baseTime, time.Duration(6) * time.Hour)
}

func UtilTimeFewDurationLater(baseTime time.Time, duration time.Duration) time.Time {
    fewDurationLater := baseTime.Add(duration)
    return fewDurationLater
}

func UtilTimeIsExpired(expirationTime time.Time) bool {
    after := time.Now().After(expirationTime)
    return after
}

func UtilTimeGetCommonDateYdmHis(t time.Time) string {
    return t.Format("2006-01-02 15:04:05")
}

package date

type DayOfWeek int

const (
	_ DayOfWeek = iota

	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

var key_DayOfWeek = map[DayOfWeek]string{
	Monday:    "Monday",
	Tuesday:   "Tuesday",
	Wednesday: "Wednesday",
	Thursday:  "Thursday",
	Friday:    "Friday",
	Saturday:  "Saturday",
	Sunday:    "Sunday",
}

var val_DayOfWeek = map[string]DayOfWeek{
	"Monday":    Monday,
	"Tuesday":   Tuesday,
	"Wednesday": Wednesday,
	"Thursday":  Thursday,
	"Friday":    Friday,
	"Saturday":  Saturday,
	"Sunday":    Sunday,
}

func (dow DayOfWeek) IsValid() bool {
	return (Monday <= dow) && (dow <= Sunday)
}

func (dow DayOfWeek) String() string {
	s, _ := key_DayOfWeek[dow]
	return s
}

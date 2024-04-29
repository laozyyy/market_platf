package common

const NoDataErr = NoDataError("no data")

type NoDataError string

func (e NoDataError) Error() string {
	return string(e)
}

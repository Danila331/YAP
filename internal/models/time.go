package models

import "github/Danila331/YAP/internal/store"

type TimeOperations struct {
	TimePulse string `db:"timepulse" json:"timepulse"`
	TimeMinus string `db:"timeminus" json:"timeminus"`
	TimeProz  string `db:"timeproz" json:"timeproz"`
	TimeDel   string `db:"timedel" json:"timedel"`
}

type TimeOperationsInterface interface {
	Update() error
	Read() (TimeOperations, error)
}

func (t *TimeOperations) Update() error {
	conn, err := store.ConnectToDatabase()
	defer conn.Close()
	if err != nil {
		return err
	}

	_, err = conn.Query("UPDATE time SET timepulse = ?, timeminus = ?, timeproz = ?, timedel = ?",
		t.TimePulse,
		t.TimeMinus,
		t.TimeProz,
		t.TimeDel)
	if err != nil {
		return err
	}
	return nil
}

func (t *TimeOperations) Read() (TimeOperations, error) {
	var time TimeOperations
	conn, err := store.ConnectToDatabase()
	defer conn.Close()
	if err != nil {
		return TimeOperations{}, err
	}

	err = conn.QueryRow("SELECT * FROM time").Scan(&time.TimePulse,
		&t.TimeMinus,
		&t.TimeProz,
		&t.TimeDel)

	if err != nil {
		return TimeOperations{}, err
	}

	return time, nil
}

package models

type (
	// GetUserID delcare
	GetUserID struct {
		ID int `json:"id"`
	}

	// GetBalance declare
	GetBalance struct {
		ID      int `json:"id"`
		UserID  int `json:"user_id"`
		Balance int `json:"balance"`
	}
)

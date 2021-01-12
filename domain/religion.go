package domain

import "time"

type Religion struct {
	ID 			int64	
	Name		string 		
	CreatedAt 	time.Time 	
	UpdatedAt 	time.Time
}
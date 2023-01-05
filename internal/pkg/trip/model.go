package trip

// {
//     "from_date": "",
//     "to_date": "",

//     "name": "",
//     "description": "",

//     "location": {
//         "city": "",
//         "state": "",
//         "country": ""
//     },
//     "participants": [
//         {
//             "name": "",
//             "photo": ""
//         },
//         {
//             "name": "",
//             "photo": ""
//         },
//         {
//             "name": "",
//             "photo": ""
//         },
//         {
//             "name": "",
//             "photo": ""
//         }
//     ]
// }

type Participant struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Photo string `json:"photo"`
}

type Location struct {
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
}

type Trip struct {
	PK          string   `json:"PK"`
	SK          string   `json:"SK"`
	ID          string   `json:"id"`
	CreatedBy   string   `json:"created_by"`
	FromDate    string   `json:"from_date"`
	ToDate      string   `json:"to_date"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Location    Location `json:"location"`
}

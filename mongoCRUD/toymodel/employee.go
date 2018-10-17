package toymodel

//Employee ...
//swagger:model
type Employee struct {
	ID      int    `json:"id" bson:"id"`
	Name    string `json:"name" bson:"name"`
	Dept    string `json:"dept" bson:"dept"`
	Address string `json:"address" bson:"address"`
}

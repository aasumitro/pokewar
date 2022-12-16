package rest

type pokeapiRESTRepository struct {
	//client *appHttp.Client
}

func NewPokeapiRESTRepository() any {
	return &pokeapiRESTRepository{}
}

package model

type ShippingData struct {
	AddressLine1       string `json:"address1"`
	AddressLine2       string `json:"address2"`
	AddressLine3       string `json:"address3"`
	AdministrativeArea string `json:"administrativeArea"`
	CountryCode        string `json:"countryCode"`
	Locality           string `json:"locality"`
	Name               string `json:"name"`
	PhoneNumber        string `json:"phoneNumber"`
	PostalCode         string `json:"postalCode"`
	SortingCode        string `json:"sortingCode"`
}

package main

type ShortenUrlResponse struct {
	ShortUrl string `json:"short_url"`
}

type GenerateUrlRequest struct {
	Url string `json:"long_url"`
}

type TinyUrlSchema struct {
	Id string `bson:"_id"`
	ShortUrl string `bson:"shortUrl"`
	OriginalUrl string `bson:"originalUrl"`
	Count int `bson:"count"`
}


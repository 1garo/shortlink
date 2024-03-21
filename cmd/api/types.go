package api

type ShortenUrlResponse struct {
	ShortUrl string `json:"short_url"`
}

type ShortenUrlRequest struct {
	Url string `bson:"url" json:"url" binding:"required"`
}

type TinyUrlSchema struct {
	Id string `bson:"_id"`
	ShortUrl string `bson:"shortUrl"`
	OriginalUrl string `bson:"originalUrl"`
	Count int `bson:"count"`
}


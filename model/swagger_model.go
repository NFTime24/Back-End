package model

type ExibitionCreateParam struct {
	Name        string `json:"name" example: "RED ROOM"`
	Description string `json:"description" example: "이미 사랑이 충만한 사람, 언젠가 다가올 사랑을 꿈꾸는 사람, 사랑에 상처받고 지겨운 사람. 그럼에도 불구하고, 사랑하는 그날을 위하여! 사랑이 넘치는 레드룸에서"`
	StartDate   string `json:"start_date" example: "2022-04-28"`
	EndDate     string `json:"end_date" example: "2022-11-06"`
	FileId      string `json:"file_id" example: "45"`
	link        string `json:"link" example: "https://tickets.interpark.com/goods/22003677?app_tapbar_state=hide&"`
}
type LikeCreateParam struct {
	UserId string `json:"user_id" example: "1"`
	WorkId string `json:"work_id" example: "1"`
}

type ArtistCreateParam struct {
	ArtistName      string `json:"artist_name" example: "Claude Monet"`
	ArtistAddress   string `json:"artist_address" example: "0x436c61756465204d6f6e6574"`
	ArtistProfileId string `json:"artist_profile_id" example: "41"`
}

type WorkInfoCreateParam struct {
	WorkName        string `json:"work_name" example: "monet the rising sun"`
	WorkPrice       string `json:"work_price" example: "20000"`
	WorkDescription string `json:"work_description" example: "moving sun"`
	WorkCategory    string `json:"work_category" example: "Image/GIF"`
	FileId          string `json:"file_id" example: "42"`
	ArtistId        string `json:"artist_id" example: "41"`
}

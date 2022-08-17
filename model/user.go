package model

type File struct {
	ID          uint `gorm:"unique;autoIncrement:false"`
	Filename    string
	Filesize    uint
	Filetype    string
	Path        string
	ThumbnailID *uint `json:"thumbnail_id" gorm:"unique; type:uuid"`
	Reference   *File `gorm:"foreignKey:ThumbnailID; References:ID"`
}

type Artist struct {
	ID        uint `gorm:"primarykey"`
	Name      string
	Address   string
	File      File `gorm:"foreignkey:ProfileID"`
	ProfileID uint
}

type User struct {
	ID        uint   `gorm:"primarykey;"`
	Address   string `gorm:"unique"`
	Nickname  string
	File      File `gorm:"foreignkey:ProfileID"`
	ProfileID uint
}

type Work struct {
	WorkID      uint `gorm:"primarykey;"`
	Name        string
	Price       uint
	Description string
	Category    string
	File        File `gorm:"foreignkey:FileID"`
	FileID      uint
	Artist      Artist `gorm:"foreignkey:ArtistID"`
	ArtistID    uint
}
type Exibition struct {
	ExibitionID uint `gorm:"primarykey"`
	Name        string
	// Nft         Nft `gorm:"foreignkey:NftsID"`
	// NftsID      uint
	Description string
	StartDate   string
	EndDate     string
	File        File `gorm:"foreignkey:FileID"`
	FileID      uint
	Link        string
}
type Nft struct {
	NftID        uint      `gorm:"primarykey"`
	Exibition    Exibition `gorm:"foreignkey:ExibitionsID"`
	ExibitionsID uint
	Work         Work `gorm:"foreignkey:WorksID"`
	WorksID      uint
	User         User `gorm:"foreignkey:OwnerID"`
	OwnerID      uint
}

type Like struct {
	LikeID  uint `gorm:"primarykey"`
	User    User `gorm:"foreignkey:OwnerID"`
	OwnerID uint
	Work    Work `gorm:"foreignkey:WorksID"`
	WorksID uint
}

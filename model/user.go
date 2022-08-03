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

type Nft struct {
	NftID   int  `gorm:"primarykey"`
	Work    Work `gorm:"foreignkey:WorksID"`
	WorksID uint
	User    User `gorm:"foreignkey:OwnerID"`
	OwnerID uint
}

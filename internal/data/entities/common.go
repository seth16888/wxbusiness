package entities

// CommonTimeFields Gorm common time fields: createdAt, updatedAt, deletedAt
type CommonTimeFields struct {
	CreatedAt int64 `bson:"created_at" json:"created_at"`
	UpdatedAt int64 `bson:"updated_at" json:"updated_at"`
	DeletedAt int64 `bson:"deleted_at" json:"deleted_at"`
}

// VersionField Gorm version field
type VersionField struct {
	Version int64 `bson:"version" json:"-"`
}

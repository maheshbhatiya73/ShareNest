package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type File struct {
    ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
    Filename     string             `json:"filename,omitempty" bson:"filename,omitempty"`
    Filepath     string             `json:"filepath,omitempty" bson:"filepath,omitempty"`
    DownloadLink string             `json:"downloadLink,omitempty" bson:"downloadLink,omitempty"`
    PasswordHash string             `json:"passwordHash,omitempty" bson:"passwordHash,omitempty"`
    CreatedAt    primitive.DateTime `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
    ExpiresAt    primitive.DateTime `json:"expiresAt,omitempty" bson:"expiresAt,omitempty"`
}

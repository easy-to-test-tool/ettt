package ettt

import "github.com/google/uuid"

/*
Evidence
証跡構造体
*/
type Evidence struct {
	Id   uuid.UUID
	Name string
	Path string
}

package model

import "time"

type PekerjaanAlumni struct {
	ID            int       `json:"id"`
	NimAlumni     string    `json:"nim_alumni"`
	StatusKerja   string    `json:"status_kerja"`
	JenisIndustri string    `json:"jenis_industri"`
	Jabatan       string    `json:"jabatan"`
	Pekerjaan     string    `json:"pekerjaan"`
	Gaji          int       `json:"gaji"`
	LamaBekerja   int       `json:"lama_bekerja"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Trash struct {
	PekerjaanAlumni
	IsDeleted   time.Time `json:"is_deleted"`
}
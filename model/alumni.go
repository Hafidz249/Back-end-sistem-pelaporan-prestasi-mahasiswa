package model

import "time"

type Alumni struct {
	UserID     *int      `json:user_id`
	NIM        string    `json:"nim"`
	Nama       string    `json:"nama"`
	Angkatan   *int      `json:"angkatan"`
	TahunLulus *int      `json:"tahun_lulus"`
	IDFakultas *int      `json:"id_fakultas"`
	IDProdi    *int      `json:"id_prodi"`
	IDSumber   *int      `json:"id_sumber"`
	Sumber     *string   `json:"sumber"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
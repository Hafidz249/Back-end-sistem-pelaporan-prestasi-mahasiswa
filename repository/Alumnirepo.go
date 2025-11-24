package repository

import (
	"POJECT_UAS/config"
	"POJECT_UAS/model"
)

func CheckAlumniByNim(nim string) (*model.Alumni, error) {
	alumni := new(model.Alumni)
	query := `SELECT nim, nama, angkatan, id_fakultas, id_prodi, tahun_lulus, sumber, id_sumber
        FROM alumni WHERE nim = $1 LIMIT 1`
	err := config.DB.QueryRow(query, nim).Scan(&alumni.NIM, &alumni.Nama, &alumni.Angkatan, &alumni.IDFakultas, &alumni.IDProdi,
		&alumni.TahunLulus, &alumni.Sumber, &alumni.IDSumber)
	if err != nil {
		return nil, err
	}
	return alumni, nil
}

func CreateAlumni(alumni *model.Alumni) error {
	query := `INSERT INTO alumni (nim, nama, angkatan, id_fakultas, id_prodi, tahun_lulus, sumber, id_sumber)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := config.DB.Exec(query, alumni.NIM, alumni.Nama, alumni.Angkatan, alumni.IDFakultas, alumni.IDProdi,
		alumni.TahunLulus, alumni.Sumber, alumni.IDSumber)
	return err
}

func UpdateAlumni(nim string, alumni *model.Alumni) error {
	query := `UPDATE alumni SET nama=$1, angkatan=$2, tahun_lulus=$3, id_fakultas=$4, id_prodi=$5, sumber=$6, id_sumber=$7
        WHERE nim=$8`
	_, err := config.DB.Exec(query, alumni.Nama, alumni.Angkatan, alumni.TahunLulus, alumni.IDFakultas, alumni.IDProdi,
		alumni.Sumber ,alumni.IDSumber , nim)
	return err
}

func DeleteAlumni(nim string) error {
	query := `DELETE FROM alumni WHERE nim=$1`
	_, err := config.DB.Exec(query, nim)
	return err
}

func GetAllAlumni() ([]model.Alumni, error) {
	query := `SELECT nim, nama, angkatan, tahun_lulus, id_fakultas, id_prodi, sumber, id_sumber FROM alumni`
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alumniList []model.Alumni
	for rows.Next() {
		alumni := model.Alumni{}
		err := rows.Scan(&alumni.NIM, &alumni.Nama, &alumni.Angkatan, &alumni.TahunLulus, &alumni.IDFakultas, &alumni.IDProdi,
			&alumni.Sumber, &alumni.IDSumber)
		if err != nil {
			return nil, err
		}
		alumniList = append(alumniList, alumni)
	}
	return alumniList, nil
}
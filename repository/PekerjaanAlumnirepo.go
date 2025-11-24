package repository

import (
	"POJECT_UAS/config"
	"POJECT_UAS/model"
)

func CheckpekerjaanAlumniByID(id string) (*model.PekerjaanAlumni, error) {
	pekerjaan := new(model.PekerjaanAlumni)
	query := `
		SELECT id, nim_alumni, status_kerja, jenis_industri, pekerjaan,
		    jabatan, gaji, lama_bekerja
		FROM pekerjaan_alumni WHERE id = $1 LIMIT 1`
	err := config.DB.QueryRow(query, id).Scan(
		&pekerjaan.ID,
		&pekerjaan.NimAlumni,
		&pekerjaan.StatusKerja,
		&pekerjaan.JenisIndustri,
		&pekerjaan.Pekerjaan,
		&pekerjaan.Jabatan,
		&pekerjaan.Gaji,
		&pekerjaan.LamaBekerja,
	)
	if err != nil {
		return nil, err
	}
	return pekerjaan, nil
}

func CreatepekerjaanAlumni(pekerjaan *model.PekerjaanAlumni) error {
	query := `
		INSERT INTO pekerjaan_alumni (
			nim_alumni, status_kerja, jenis_industri, pekerjaan,
			jabatan, gaji, lama_bekerja
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`
	err := config.DB.QueryRow(query,
		pekerjaan.NimAlumni,
		pekerjaan.StatusKerja,
		pekerjaan.JenisIndustri,
		pekerjaan.Pekerjaan,
		pekerjaan.Jabatan,
		pekerjaan.Gaji,
		pekerjaan.LamaBekerja,
	).Scan(&pekerjaan.ID)
	return err
}

func UpdatepekerjaanAlumni(NimAlumni string, pekerjaan *model.PekerjaanAlumni) error {
	query := `
		UPDATE pekerjaan_alumni
		SET status_kerja = $1, jenis_industri = $2, pekerjaan=$3, jabatan = $4,
		    gaji = $5, lama_bekerja = $6, pekerjaan = $7
		WHERE nim_alumni = $8`
	_, err := config.DB.Exec(query,
		pekerjaan.StatusKerja,
		pekerjaan.JenisIndustri,
		pekerjaan.Pekerjaan,
		pekerjaan.Jabatan,
		pekerjaan.Gaji,
		pekerjaan.LamaBekerja,
		pekerjaan.Pekerjaan,
		NimAlumni,
	)
	return err
}

func GetAllpekerjaanAlumni() ([]model.PekerjaanAlumni, error) {
	query := `SELECT id, nim_alumni, status_kerja, jenis_industri, pekerjaan, jabatan, gaji, lama_bekerja
		FROM pekerjaan_alumni`
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pekerjaanList []model.PekerjaanAlumni
	for rows.Next() {
		var pekerjaan model.PekerjaanAlumni
		err := rows.Scan(
			&pekerjaan.ID,
			&pekerjaan.NimAlumni,
			&pekerjaan.StatusKerja,
			&pekerjaan.JenisIndustri,
			&pekerjaan.Pekerjaan,
			&pekerjaan.Jabatan,
			&pekerjaan.Gaji,
			&pekerjaan.LamaBekerja,
		)
		if err != nil {
			return nil, err
		}
		pekerjaanList = append(pekerjaanList, pekerjaan)
	}
	return pekerjaanList, nil
}

func SoftDeleteBynim(NimAlumni string) error {
	query := `UPDATE pekerjaan_alumni SET is_deleted = NOW() WHERE id = $1`

	_, err := config.DB.Exec(query, NimAlumni)
	return err
}

func GetAllTrash(nimAlumni string) ([]*model.Trash, error) {
	var trashes []*model.Trash

	query := `
	SELECT 
            pa.id, pa.nim_alumni, pa.status_kerja, pa.jenis_industri, pa.jabatan, 
            pa.pekerjaan, pa.gaji, pa.lama_bekerja, pa.created_at, pa.updated_at, 
            pa.is_deleted
        FROM pekerjaan_alumni pa
		JOIN alumni a ON a.nim = pa.nim_alumni
        WHERE pa.is_deleted IS NOT NULL
    `
    var args []interface{}

    if nimAlumni != "" {
        query += " AND pa.nim_alumni = $1"
        args = append(args, nimAlumni)
    }

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		trash := new(model.Trash)
		if err := rows.Scan(
            &trash.ID,
            &trash.NimAlumni,
            &trash.StatusKerja,
            &trash.JenisIndustri,
            &trash.Jabatan,
            &trash.Pekerjaan,
            &trash.Gaji,
            &trash.LamaBekerja,
            &trash.CreatedAt,
            &trash.UpdatedAt,
            &trash.IsDeleted,
		); err != nil {
			return nil, err
		}
		trashes = append(trashes, trash)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return trashes, nil
}

func RestoreTrashBynim(NimAlumni string) error {
	query :=`UPDATE pekerjaan_alumni SET is_deleted = NULL WHERE id = $1`
		_, err := config.DB.Exec(query, NimAlumni)
	return err
}

func DeletePekerjaanByid(NimAlumni string) error {
	query := `DELETE FROM pekerjaan_alumni WHERE id = $1 AND is_deleted IS NOT NULL`

	_, err := config.DB.Exec(query, NimAlumni)
	return err
}
# Gambaran Umum

### Visi Misi

#### Visi

Menjadikan Puskesmas sebagai Pusat Pelayanan Kesehatan unggulan yang berkualitas, profesional dan terjangkau

#### Misi

- Memberikan pelayanan kesehatan yang berkualitas, merata dan terjangkau bagi seluruh lapisan masyarakat.
- Meningkatkan kualitas sumber daya manusia (SDM) yang profesional.
- Memelihara dan meningkatkan sarana dan prasarana pelayanan kesehatan.
- Memberdayakan masyarakat dalam bidang kesehatan dengan menjalin kerja sama lintas sektoral.

### Letak Geografis [(Google Maps)](https://goo.gl/maps/TvHWhvgrHLAi4vEQ8)

- Wilayah kerja Puskesmas Pasir Nangka Kecamatan Tigaraksa meliputi Desa Pasir Nangka, Desa Pasir Baolang, Desa Pete, Desa Pematang, Desa Tegalsari, Desa Matagara dan Desa Cisereh dengan luas wilayah kerja 48,74 Km².
- Puskesmas Pasir Nangka terletak di Jl. Aria Jaya Santika Desa Pasir Nangka Kecamatan Tigaraksa. Kabupaten Tangerang dengan batasan wilayah administrasi sebagai berikut:
- Sebelah Utara berbatasan dengan Kecamatan Cikupa dan Balaraja
- Sebelah Selatan berbatasan dengan Puskesmas Tigaraksa
- Sebelah Barat berbatasan dengan Kecamatan Cisoka
- Sebelah Timur berbatasan dengan Kecamatan Cikupa dan Puskesmas Tigaraksa.

### Wilayah Administrasi

1. Desa Pasir Nangka
2. Desa Pasir Bolang
3. Desa Pete
4. Desa Pematang
5. Desa Tegalsari
6. Desa Matagara
7. Desa Cisereh

### Keadaan Penduduk

Jumlah Penduduk di Wilayah Puskesmas Pasir Nangka tahun 2020, menurut data dari BPS Kecamatan Tigaraksa sebanyak 104.254 jiwa. Jumlah Penduduk tertinggi di Desa Pasir Nangka yang berjumlah 36.197 jiwa, sedangkan yang terendah di Desa Tegalsari yaitu 4.161 jiwa.

Dari data BPS Kecamatan Tigaraksa, Pasir Nangka menunjukan struktur penduduk usia produktif yaitu usia 15 s/d 64 tahun adalah 65,89% dari keseluruhan jumlah penduduk yaitu 68.685 jiwa, berumur 0 s/d 14 tahun adalah 31,10 % yaitu 32.424 jiwa dan 3.01 % adalah penduduk berumur > 65 tahun atau 3.145 jiwa. Dilihat dari jenis kelamin, penduduk berjenis kelamin laki-laki prosentasenya 51,06 % (53.238 jiwa) dan perempuan 48,93 % (50.016 jiwa).

### SDM

Jumlah dan Jenis tenaga di Puskesmas Pasir Nangka

import { useState, useEffect } from 'react';



export const Manpower = () => {
    const [manpowerCategories, setManpowerCategories] = useState([])
    const [manpowerStatuses, setManpowerStatuses] = useState([])
    const [manpowerStatusAmounts, setManpowerStatusAmounts] = useState([])
    useEffect(() => {
        fetchData()
    }, [])
    const fetchData = async () => {
        await Promise.all([
            fetchManpowerCategories(),
            fetchManpowerStatuses(),
            fetchManpowerStatusAmounts(),
        ])
    }
    const fetchManpowerCategories = async () => {
        try {
            const resp = await fetch(`${process.env.REACT_APP_BASE_URL}/manpowercategories`)
            if (resp.status === 200) {
                setManpowerCategories(await resp.json())
            } else {
                console.error('Status error')
            }
        } catch (e) {
        console.error(e)
        } 
    }
    const fetchManpowerStatuses = async () => {
        try {
            const resp = await fetch(`${process.env.REACT_APP_BASE_URL}/manpowerstatuses`)
            if (resp.status === 200) {
                setManpowerStatuses(await resp.json())
            } else {
                console.error('Status error')
            }
        } catch (e) {
        console.error(e)
        } 
    }
    const fetchManpowerStatusAmounts = async () => {
        try {
            const resp = await fetch(`${process.env.REACT_APP_BASE_URL}/manpowerstatusamounts`)
            if (resp.status === 200) {
                setManpowerStatusAmounts(await resp.json())
            } else {
                console.error('Status error')
            }
        } catch (e) {
        console.error(e)
        } 
    }
    return (
        <div>
            <div style={{ overflow: "auto", resize: "vertical", height: "65vh" }}>
                <table border="1" style={{ width: "100%", borderCollapse: "separate" }}>
                <thead>
                    <tr style={{ position: "sticky", top: 0, zIndex: 1 }}>
                    <th style={{ position: "sticky", top: 0, zIndex: 1 }} rowSpan={2}>#</th>
                    <th style={{ position: "sticky", top: 0 }} rowSpan={2}>Category 
                        
                    </th>
                    <th colSpan={manpowerStatuses.length ?? 0}>Status 
                      
                    </th>
                    <th style={{ position: "sticky", top: 0 }} rowSpan={2}>Total</th>
                    </tr>
                    <tr>
                    {manpowerStatuses.length > 0 
                        ? manpowerStatuses.map((manpowerStatus, i) => {
                            return (
                            <td style={{ position: "sticky", top: 0 }}>
                                {manpowerStatus?.name ?? ''}
                            </td>
                            )
                        })
                        : <></>
                    }
                    </tr>
                    {manpowerCategories.map((manpowerCategory, i) => {
                    return (
                        <tr>
                        <td>{i + 1}</td>
                        <td>
                            {manpowerCategory?.name ?? ''}
                        </td>
                        {manpowerStatuses.map((manpowerStatus, i) => {
                            const foundManpowerStatusAmount = manpowerStatusAmounts?.find(manpowerStatusAmount =>
                            manpowerStatusAmount?.manpowerCategoryUuid === manpowerCategory?.uuid &&
                            manpowerStatusAmount?.manpowerStatusUuid === manpowerStatus?.uuid
                            )
                            return (
                            <td>
                                {foundManpowerStatusAmount?.value ?? 0}
                            </td>
                            )
                        })}
                        <td>
                            {manpowerStatusAmounts
                            .filter(manpowerStatusAmount =>
                                manpowerStatusAmount?.manpowerCategoryUuid === manpowerCategory?.uuid
                            )
                            .map(m => {
                              console.log(manpowerCategories.find(mc => mc.uuid === m?.manpowerCategoryUuid)?.name)
                              return m;
                            })
                            .reduce((acc, manpowerStatusAmount) => acc + (manpowerStatusAmount?.value ?? 0), 0)
                            }
                        </td>
                        </tr>
                    )
                    })}
                    <tr>
                        <th colSpan={2}>
                            Total
                        </th>
                        {manpowerStatuses.map(manpowerStatus => {
                            return (
                                <th>
                                    {manpowerStatusAmounts
                                        .filter(manpowerStatusAmount => manpowerStatusAmount?.manpowerStatusUuid === manpowerStatus?.uuid )
                                        .reduce((acc, manpowerStatusAmount) => acc + (manpowerStatusAmount?.value ?? 0), 0)
                                    }
                                </th>
                            )
                        })}
                        <th>
                            {manpowerStatusAmounts
                                .reduce((acc, manpowerStatusAmount) => acc + (manpowerStatusAmount?.value ?? 0), 0)
                            }
                        </th>
                        
                    </tr>
                </thead>
                </table>
            </div>
        </div>
    )
}

<Manpower />

<!-- ![sdm](/img/05-sdm.png) -->

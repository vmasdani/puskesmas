# Sidumas

Halo, selamat datang di SIDUMAS (Sistim Pengaduan Masyarakat)!  

Dalam rangka  meningkatkan mutu pelayanan, Puskesmas Pasir Nangka membuka pengaduan pelayanan kesehatan  bagi masyarakat melalui website Puskesmas Pasir Nangka.  

Pengaduan secara langsung diterima oleh petugas dan apabila dapat diselesaikan langsung dan dilaporkan kepada Unit Pengaduan Masyarakat Puskesmas Pasir Nangka.

Kami menjamin kerahasiaan nomor telepon dan alamat anda. Nomor telepon dan alamat anda hanya diteruskan ke admin, tapi tidak dimunculkan di website.

import { useState, useEffect } from 'react'

export const Sidumas = () => {
  const [baseUrl,  setBaseUrl] = useState('')
  const [complaints, setComplaints] = useState([])
  const [name, setName] = useState('')
  const [phone, setPhone] = useState('')
  const [address, setAddress] = useState('')
  const [complaint, setComplaint] = useState('')
  
  useEffect(() => {
    setBaseUrl(process.env.REACT_APP_BASE_URL)
    fetchComplaints()
  }, [])
  const fetchComplaints = async () => {
    try {
      const resp = await fetch(`${process.env.REACT_APP_BASE_URL}/complaints`)
      if (resp.status !== 200) throw await resp.text()
      setComplaints(await resp.json())
    } catch(e) {
      console.error(e)
    }
  }
  const handleSubmit = async (e) => {
    e.preventDefault();
    const confirm = window.confirm('Kirim form sekarang?')
  
    if (confirm) {
      try {
        const resp = await fetch(`${baseUrl}/complaints`, {
          method: 'post',
          headers: {
            'content-type': 'application/json',
          },
          body: JSON.stringify({
            name: name,
            phone: phone,
            address: address,
            complaint: complaint
          })
        })
        if (resp.status !== 201) throw await resp.text()
        alert('Pengiriman sukses!')
        window.location.reload() 
      } catch (e) {
        console.error(e)
        alert('Pengiriman gagal! Harap coba kembali beberapa saat')
      }
    }
  }
  return (
    <div >
      <div 
        style={{ 
          display:"flex", 
          flexWrap: "wrap", 
          justifyContent: "space-between",
        }}
      >
        <div>
          <div>Form Pengaduan</div>
          <form 
            onSubmit={e => { handleSubmit(e); }}
            style={{display:"flex", flexDirection:"column"}}
          >
            <input 
              onChange={e => {
                setName(e.target.value)
              }} 
              placeholder="Nama..." style={{ marginTop: 5, marginBottom: 5, paddingTop: 5, paddingBottom: 5 }} 
            />
            <input
              onChange={e => {
                setPhone(e.target.value)
              }} 
              placeholder="Telp..." style={{ marginTop: 5, marginBottom: 5, paddingTop: 5, paddingBottom: 5 }} 
            />
            <input 
              onChange={e => {
                setAddress(e.target.value)
              }} 
              placeholder="Alamat..." style={{ marginTop: 5, marginBottom: 5, paddingTop: 5, paddingBottom: 5 }} 
            />
            <textarea
              onChange={e => {
                setComplaint(e.target.value)
              }} 
              placeholder="Pengaduan anda..." 
              style={{ 
                marginTop: 5, 
                marginBottom: 5, 
                paddingTop: 5, 
                paddingBottom: 5, 
                resize: "vertical",
                height: "25vh", 
              }} 
            />
            <input type="submit" className="btn btn-blue" value="Kirim" />
          </form>
        </div>
        
        <div
          style={{
            justifyContent: "start", 
            flexGrow: 1, 
            marginLeft: 15,
          }}
        >
          <div>Jumlah total pengaduan:{" "}
            {complaints?.length ?? 0},{" "}
            terjawab: {complaints?.filter(complaint => complaint?.answer !== "" )?.length ?? 0}/{complaints?.length ?? 0}
          </div>
          <div 
            style={{ 
              width: "100%", 
              marginLeft: 5, 
              marginRight: 5,
              border: "2px solid grey",
              padding: 10,
              overflow: "auto",
              height: "50vh",
              resize: "vertical",
            }}
          >
            <div>
              {complaints.reverse().map((complaint, i) => {
                return (
                  <div>
                    <div>
                      <strong>
                        {(complaints?.length ?? 0) - i}) Pengirim: {complaint?.name}
                      </strong>
                    </div>
                    <div>Pengaduan:</div>
                    <div>
                      <textarea 
                        readOnly 
                        value={complaint?.complaint} 
                        style={{ width: "100%", resize: "vertical", height: "10vh" }} 
                      />
                    </div>
                    <div>Jawaban:</div>
                    <div>
                      <textarea 
                        readOnly  
                        style={{ width: "100%", resize: "vertical", height: "7.5vh" }} 
                      />
                    </div>
                    <div>
                      Tanggal: {complaint?.createdAt 
                        ? Intl.DateTimeFormat(navigator.language ?? 'en-US',  {
                            dateStyle: "long",
                            timeStyle: "long"
                          }).format(new Date(complaint.createdAt)) 
                        : 'No date' 
                      }
                    </div>
                    <hr />
                  </div>
                )
              })}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

<Sidumas />

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
  const [complaintToEdit, setComplaintToEdit] = useState(null)
  const [authorized, setAuthorized] = useState(false)
  
  useEffect(() => {
    setBaseUrl(process.env.REACT_APP_BASE_URL)
    fetchComplaints()
    fetchAuthorizationStatus()
  }, [])
  const fetchAuthorizationStatus = async () => {
    try {
      console.log('apikey:', localStorage.getItem('apiKey'))
      if (localStorage.getItem('apiKey')) {
        const resp = await fetch(`${process.env.REACT_APP_BASE_URL}/authorize-admin`, {
          headers: {
            authorization: localStorage.getItem('apiKey')
          }
        })
        if (resp.status !== 200) throw await resp.text()
        setAuthorized(true)
      }
    } catch (e) {
      console.error(e)
    }
  }
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
  const handleSubmitForm = async (e) => {
    try {
      const resp = await fetch(`${process.env.REACT_APP_BASE_URL}/complaints-save`, {
        method: 'post',
        headers: {
          authorization: localStorage.getItem('apiKey'),
          'content-type': 'application/json'
        },
        body: JSON.stringify(complaintToEdit)
      })
      if (resp.status !== 201) throw await resp.text()
      window.location.reload()
    } catch (e) {
      console.error(e)
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
              {(() => {
                const complaintsReversed = [...complaints]
                complaintsReversed.reverse()

                return complaintsReversed.map((complaint, i) => {
                  return (
                    <div style={{ }}>
                      <div 
                        style={{
                          display:"flex",
                          width: "100%",
                          justifyContent: "space-between",
                          color: complaint?.answer === "" ? 'crimson' : 'green'
                        }}
                      >
                        <strong>
                          {(complaints?.length ?? 0) - i}) Pengirim: {complaint?.name}
                        </strong>
                        {authorized 
                          ? <button 
                              onClick={() => setComplaintToEdit(complaint)}
                            >
                              Edit
                            </button>
                          :<></>
                        }
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
                          value={complaint?.answer}
                          style={{ width: "100%", resize: "vertical", height: "10vh" }} 
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
                })
              })()}
            </div>
          </div>
        </div>
      </div>
      {complaintToEdit
          ? <form onSubmit={handleSubmitForm}>
              <div style={{marginTop: 15}}>
                <div style={{ display:"flex", width: "100%"}}>
                  Nama:<input value={complaintToEdit?.name} style={{width:"100%", marginLeft: 5}} />
                </div>
                <div style={{ display:"flex", width: "100%"}}>
                  Telp:<input value={complaintToEdit?.phone} style={{width:"100%", marginLeft: 5}} />
                </div>
                <div style={{ display:"flex", width: "100%"}}>
                  Alamat:<input value={complaintToEdit?.address} style={{width:"100%", marginLeft: 5}} />
                </div>
                <div style={{ display:"flex", width: "100%"}}>
                  Aduan:<textarea value={complaintToEdit?.complaint} style={{width:"100%", marginLeft: 5, resize: "vertical", height:"10vh"}} />
                </div>
                <div style={{ display:"flex", width: "100%"}}>
                  Jawaban:
                    <textarea
                      onChange={e => {
                        setComplaintToEdit({
                          ...complaintToEdit,
                          answer: e.target.value
                        })
                      }} 
                      value={complaintToEdit?.answer} 
                      style={{
                        width:"100%", 
                        marginLeft: 5, 
                        resize: "vertical", 
                        height:"10vh"
                      }} 
                    />
                </div>
                <input type="submit" value="save" />
              </div>
            </form>
            
          : <></>
        }
    </div>
  )
}

<Sidumas />

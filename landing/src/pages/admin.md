### Admin Page

import { useState, useEffect } from 'react'

export const AdminPage = () => {
  const [authorized, setAuthorized] = useState(false)
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  
  useEffect(() => {
    console.log(process.env.REACT_APP_BASE_URL)
    // todo: check authorization
    if (localStorage.getItem('apiKey'))
      setAuthorized(true)
  }, [])
  const handleLogin = async (e) => {
    try {
      e.preventDefault()
      const resp = await fetch(`${process.env.REACT_APP_BASE_URL}/login`, {
        method: 'post',
        headers: {
          'content-type':'application/json'
        },
        body: JSON.stringify({
          username: username,
          password: password
        })
      })
      console.log(e)
      if (resp.status !== 200) throw await resp.text()
      const apiKey = await resp.text()
      localStorage.setItem('apiKey', apiKey)
      setAuthorized(true)
    } catch (e) {
      alert(e)
      console.error(e)
    }
  }
  const handleLogout = async () => {
    try {
      localStorage.removeItem('apiKey')
      window.location.reload()
    } catch (e) {
      console.error(e)
    }
  }
  return (
    <div>
      {authorized 
        ? <div>
            Logged in.
            <div>
              <div>
                <button onClick={() => handleLogout()}>Logout</button>
              </div>
            </div>
          </div>
        : <div>
            <form onSubmit={e => handleLogin(e)} style={{ display: "flex", flexDirection: "column" }}>
              <input onChange={e => { setUsername(e.target.value) }} type="text" placeholder="Username..." style={{ marginTop: 10, padding: 10 }} />
              <input onChange={e => { setPassword(e.target.value) }} type="password" placeholder="Password..." style={{ marginTop: 10, padding: 10 }} />
              <input type="submit" value="Login"   style={{ marginTop: 10, padding: 10 }} />
            </form>
          </div>
      }
    </div>
  )
}

<AdminPage />
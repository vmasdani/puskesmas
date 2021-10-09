### Admin Page

import { useState, useEffect } from 'react'

export const AdminPage = () => {
  const [authorized, setAuthorized] = useState(false)
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [users, setUsers] = useState([])
  const [userDeleteIds, setUserDeleteIds] = useState([])
  const [manpowerCategories, setManpowerCategories] = useState([])
  const [manpowerCategoryDeleteIds, setManpowerCategoryDeleteIds] = useState([])
  const [manpowerStatuses, setManpowerStatuses] = useState([])
  const [manpowerStatusDeleteIds, setManpowerStatusDeleteIds] = useState([])
  const [manpowerStatusAmounts, setManpowerStatusAmounts] = useState([])
  const [manpowerStatusAmountDeleteIds, setManpowerStatusAmountDeleteIds] = useState([])
  const [loading, setLoading] = useState(false)
  let nanoid=(t=21)=>{let e="",r=crypto.getRandomValues(new Uint8Array(t));for(;t--;){let n=63&r[t];e+=n<36?n.toString(36):n<62?(n-26).toString(36).toUpperCase():n<63?"_":"-"}return e};
  useEffect(() => {
    console.log(process.env.REACT_APP_BASE_URL)
    // todo: check authorization
    if (localStorage.getItem('apiKey'))
      handleAuthorizeAndFetch()
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
  const handleAuthorizeAndFetch = async () => {
    try {
      const resp = await fetch(`${process.env.REACT_APP_BASE_URL}/authorize`, {
          method: 'post',
          headers: {
            'authorization': localStorage.getItem('apiKey')
          },
        })
      if (resp.status === 200) {
        setAuthorized(true)
        await Promise.all([
          fetchManpowerCategories(),
          fetchManpowerStatuses(),
          fetchManpowerStatusAmounts(),
          fetchUsersView(),
        ])
      } else {
        console.error('Status error')
      }
    } catch (e) {
      console.error(e)
    } 
  }
  const fetchManpowerCategories = async () => {
    try {
      const resp = await fetch(`${process.env.REACT_APP_BASE_URL}/manpowercategories`, {
          headers: {
            'authorization': localStorage.getItem('apiKey')
          },
        })
      if (resp.status === 200) {
        setManpowerCategories(await resp.json())
      } else {
        console.error('Status error')
      }
    } catch (e) {
      console.error(e)
    } 
  }
  const fetchUsersView = async () => {
    try {
      const resp = await fetch(`${process.env.REACT_APP_BASE_URL}/users-view`, {})
      if (resp.status === 200) {
        setUsers(await resp.json())
      } else {
        console.error('Status error')
      }
    } catch (e) {
      console.error(e)
    } 
  }
  
  const fetchManpowerStatuses = async () => {
    try {
      const resp = await fetch(`${process.env.REACT_APP_BASE_URL}/manpowerstatuses`, {
          headers: {
            'authorization': localStorage.getItem('apiKey')
          },
        })
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
      const resp = await fetch(`${process.env.REACT_APP_BASE_URL}/manpowerstatusamounts`, {
          headers: {
            'authorization': localStorage.getItem('apiKey')
          },
        })
      if (resp.status === 200) {
        setManpowerStatusAmounts(await resp.json())
      } else {
        console.error('Status error')
      }
    } catch (e) {
      console.error(e)
    } 
  }
  
  const handleSave = async ()  => {
    try {
      setLoading(true)
      const resp = await Promise.all([
        fetch(`${process.env.REACT_APP_BASE_URL}/manpowers-save`, {
          method: 'post',
          headers: {
            'content-type': 'application/json',
            'authorization': localStorage.getItem('apiKey')
          },
          body: JSON.stringify({
            manpowerCategories: manpowerCategories,
            manpowerCategoryDeleteIds: manpowerCategoryDeleteIds,
            manpowerStatuses: manpowerStatuses,
            manpowerStatusDeleteIds: manpowerStatusDeleteIds,
            manpowerStatusAmounts: manpowerStatusAmounts,
            manpowerStatusAmountDeleteIds: manpowerStatusAmountDeleteIds,
          })
        }),
        fetch(`${process.env.REACT_APP_BASE_URL}/users-save`, {
          method: 'post',
          headers: {
            'content-type': 'application/json',
            'authorization': localStorage.getItem('apiKey')
          },
          body: JSON.stringify({userBody: users, userDeleteIds: userDeleteIds})
        }),
      ]) 
      window.location.reload()
    } catch (e) {
      setLoading(false)
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
      {authorized
        ? <div style={{ marginTop: 15 }}>
            <h3>
              Settings 
              {loading
                ? <>Loading... Please wait!</>
                : <button onClick={handleSave}>Save</button>
              }
            </h3>
            <hr />
            <h6>Users</h6>
            <div style={{display:"flex"}}>
              <button onClick={() => {
                setUsers([
                  ...users,
                  {
                    username: ''
                  }
                ])
                console.log(users)
              }}>Add</button>
            </div>
            <div>
              {users.map((user, i ) => {
                return (
                  <div style={{display:"flex",width:"100%", marginTop: 5, marginBottom: 5}}>
                    <input 
                      style={{ flexGrow: 1 }} 
                      placeholder="Name..." 
                      value={user?.name} 
                      onChange={e => { 
                        setUsers(users.map((userX, ix) => ix === i
                          ? { ...userX, name: e.target.value }
                          : userX
                        )) 
                      }} 
                    />
                    <input 
                      style={{ flexGrow: 1 }} 
                      placeholder="Username..." 
                      value={user?.username} 
                      onChange={e => { 
                        setUsers(users.map((userX, ix) => ix === i
                          ? { ...userX, username: e.target.value }
                          : userX
                        )) 
                      }}
                    />
                    <button
                      onClick={async () => {
                        try {
                          const newPassword = window.prompt('Enter new password:')
                          if (newPassword && newPassword !== '') {
                            console.log('new pass', newPassword)
                            setUsers(users.map((userX, ix) => ix === i
                              ? { ...userX, newPassword: newPassword, changePassword: true }
                              : userX
                            ))
                          }
                        } catch (e) {
                          console.error(e)
                        }
                      }}
                    >Change Password</button>
                    <button onClick={() => {
                      console.log(i)
                      console.log(users.filter((_, ix) => i !== ix))
                      setUsers(users.filter((_, ix) => i !== ix))
                    }}>Delete</button>
                  </div>
                )
              })}
            </div>
            <hr />
            <h6>Manpower</h6> 
            <div style={{ overflow: "auto", resize: "vertical", height: "40vh" }}>
              <table border="1" style={{ width: "100%", borderCollapse: "separate" }}>
                  <tr style={{ position: "sticky", top: 0, zIndex: 1 }}>
                    <th style={{ position: "sticky", top: 0, zIndex: 1 }} rowSpan={2}>#</th>
                    <th style={{ position: "sticky", top: 0 }} rowSpan={2}>Category 
                      <button 
                        onClick={() => { setManpowerCategories([...manpowerCategories, { uuid: nanoid() }]) }}
                      >
                        Add
                      </button>
                    </th>
                    <th colSpan={manpowerStatuses.length ?? 0}>Status 
                      <button 
                        onClick={() => { setManpowerStatuses([...manpowerStatuses, { uuid: nanoid() }]) }}
                      >
                        Add
                      </button>
                    </th>
                    <th style={{ position: "sticky", top: 0 }} rowSpan={2}>Total</th>
                  </tr>
                  <tr>
                    {manpowerStatuses.length > 0 
                      ? manpowerStatuses.map((manpowerStatus, i) => {
                          return (
                            <td style={{ position: "sticky", top: 0 }}>
                              <input
                                style={{ width: 75 }}
                                value={manpowerStatus?.name ?? ''}
                                onChange={e => {
                                  setManpowerStatuses(
                                    manpowerStatuses.map(manpowerStatusX => manpowerStatus?.uuid === manpowerStatusX?.uuid
                                      ? { ...manpowerStatusX, name: e.target.value }
                                      : manpowerStatusX
                                    )
                                  )
                                }}
                              />
                              <button onClick={() => {
                                setManpowerStatusDeleteIds([
                                  ...manpowerStatusDeleteIds,
                                  manpowerStatus?.id ?? 0
                                ])
                                setManpowerStatuses(manpowerStatuses.filter((_, ix) => ix !== i))
                              }}>Delete</button>
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
                          <input
                            type="text" 
                            value={manpowerCategory?.name ?? ''} 
                            onChange={e => {
                              setManpowerCategories(
                                manpowerCategories.map(manpowerCategoryX => manpowerCategory?.uuid === manpowerCategoryX?.uuid
                                  ? { ...manpowerCategoryX, name: e.target.value }
                                  : manpowerCategoryX
                                )
                              )
                            }}
                          />
                          <button onClick={() => {
                            setManpowerCategoryDeleteIds([
                              ...manpowerCategoryDeleteIds,
                              manpowerCategory?.id ?? 0
                            ])
                            setManpowerCategories(manpowerCategories.filter((_, ix) => ix !== i))
                          }}>Delete</button>
                        </td>
                        {manpowerStatuses.map((manpowerStatus, i) => {
                          const foundManpowerStatusAmount = manpowerStatusAmounts?.find(manpowerStatusAmount =>
                            manpowerStatusAmount?.manpowerCategoryUuid === manpowerCategory?.uuid &&
                            manpowerStatusAmount?.manpowerStatusUuid === manpowerStatus?.uuid
                          )
                          return (
                            <td>
                              <input 
                                type="number" 
                                style={{ width: 75 }}
                                value={foundManpowerStatusAmount ? (foundManpowerStatusAmount?.value ?? 0) : 0}
                                onChange={e => {
                                  const newManpowerStatusAmounts = [...manpowerStatusAmounts]
                                  const foundManpowerStatusAmount = 
                                    newManpowerStatusAmounts
                                      .find(manpowerStatusAmount =>
                                        manpowerStatusAmount?.manpowerCategoryUuid === manpowerCategory?.uuid && 
                                        manpowerStatusAmount?.manpowerStatusUuid === manpowerStatus?.uuid
                                      )
                                  console.log('found:', foundManpowerStatusAmount)
                                  if (foundManpowerStatusAmount) {
                                    foundManpowerStatusAmount.value = isNaN(parseInt(e.target.value)) 
                                      ? foundManpowerStatusAmount?.value
                                      : parseInt(e.target.value)
                                    console.log(isNaN(parseInt(e.target.value)) 
                                      ? foundManpowerStatusAmount?.value
                                      : parseInt(e.target.value))
                                    setManpowerStatusAmounts(newManpowerStatusAmounts)
                                  } else {
                                    console.log({
                                        manpowerCategoryUuid: manpowerCategory?.uuid,
                                        manpowerStatusUuid: manpowerStatus?.uuid,
                                        value: isNaN(parseInt(e.target.value)) 
                                          ? foundManpowerStatusAmount?.value
                                          : parseInt(e.target.value),
                                      })
                                    
                                    setManpowerStatusAmounts([...
                                      newManpowerStatusAmounts,
                                      {
                                        manpowerCategoryUuid: manpowerCategory?.uuid,
                                        manpowerStatusUuid: manpowerStatus?.uuid,
                                        value: isNaN(parseInt(e.target.value)) 
                                          ? foundManpowerStatusAmount?.value
                                          : parseInt(e.target.value),
                                      }
                                    ])
                                  }
                                }}
                              />
                            </td>
                          )
                        })}
                        <td>
                          {manpowerStatusAmounts
                            .filter(manpowerStatusAmount =>
                              manpowerStatusAmount?.manpowerCategoryUuid === manpowerCategory?.uuid
                            )
                            .reduce((acc, manpowerStatusAmount) => acc + (manpowerStatusAmount?.value ?? 0), 0)
                          }
                        </td>
                      </tr>
                    )
                  })}
              </table>
            </div>
          </div>
        :<></>
      }
      
    </div>
  )
}

<AdminPage />
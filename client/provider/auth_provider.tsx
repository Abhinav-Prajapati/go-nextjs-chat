"use client"
import { useRouter } from 'next/navigation'
import React, { createContext, useEffect, useState } from 'react'

export type UserInfo = {
  username: string
  id: string
}

export const AuthContext = createContext<{
  authenticated: boolean
  setAuthenticated: (auth: boolean) => void
  user: UserInfo
  setUser: (user: UserInfo) => void
}>({
  authenticated: false,
  setAuthenticated: () => { },
  user: {
    username: "",
    id: ""
  },
  setUser: () => { }
})

const AuthContextProvider = ({ children }: { children: React.ReactNode }) => {
  const [authenticated, setAuthenticated] = useState(false)
  const [user, setUser] = useState<UserInfo>({ username: '', id: '' })

  const rounter = useRouter()

  useEffect(() => {
    const userInfo = localStorage.getItem("user_info")
    if (!userInfo) {
      if (window.location.pathname != "/signup") {
        rounter.push("/login")
        return
      }
    } else {
      const user: UserInfo = JSON.parse(userInfo)
      if (user) {
        setUser({
          username: user.username,
          id: user.id
        })
      }
      setAuthenticated(true)
    }
  }, [authenticated])
  return (
    <AuthContext.Provider
      value={{
        authenticated: authenticated, // change this to authnticated
        setAuthenticated: setAuthenticated,
        user: user,
        setUser: setUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export default AuthContextProvider
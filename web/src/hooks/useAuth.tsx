'use client'

import { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { useRouter } from 'next/navigation'
import { authApi } from '@/lib/api-client'
import type { User } from '@bookmarketing/shared-types'

interface AuthContextType {
  user: User | null
  token: string | null
  isLoading: boolean
  login: (email: string, password: string) => Promise<void>
  register: (email: string, password: string, fullName: string) => Promise<void>
  logout: () => Promise<void>
  refreshAuth: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const router = useRouter()

  const clearAuth = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('user')
    setToken(null)
    setUser(null)
  }

  const refreshAuthToken = async (refreshToken: string) => {
    try {
      const response = await authApi.refreshToken(refreshToken)
      localStorage.setItem('token', response.token)
      localStorage.setItem('refresh_token', response.refresh_token)
      setToken(response.token)
    } catch (error) {
      clearAuth()
      throw error
    }
  }

  useEffect(() => {
    // Check for existing token on mount
    const storedToken = localStorage.getItem('token')
    const storedRefreshToken = localStorage.getItem('refresh_token')
    const storedUser = localStorage.getItem('user')

    if (storedToken && storedUser) {
      // Browser storage is only readable after hydration, so the session is restored here.
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setToken(storedToken)
      setUser(JSON.parse(storedUser))

      // Try to refresh token if it's close to expiring
      if (storedRefreshToken) {
        refreshAuthToken(storedRefreshToken).catch(() => {
          // Refresh failed, clear auth
          clearAuth()
        })
      }
    }
    setIsLoading(false)
  }, [])

  const login = async (email: string, password: string) => {
    const data = await authApi.login(email, password)
    localStorage.setItem('token', data.token)
    localStorage.setItem('refresh_token', data.refresh_token)
    localStorage.setItem('user', JSON.stringify(data.user))
    setToken(data.token)
    setUser(data.user)
    router.push('/')
  }

  const register = async (email: string, password: string, fullName: string) => {
    const data = await authApi.register(email, password, fullName)
    localStorage.setItem('token', data.token)
    localStorage.setItem('refresh_token', data.refresh_token)
    localStorage.setItem('user', JSON.stringify(data.user))
    setToken(data.token)
    setUser(data.user)
    router.push('/')
  }

  const logout = async () => {
    const currentToken = localStorage.getItem('token')
    const refreshToken = localStorage.getItem('refresh_token')

    try {
      if (currentToken && refreshToken) {
        await authApi.logout(currentToken, refreshToken)
      }
    } catch (error) {
      // Logout may fail, but we still clear local state
      console.error('Logout error:', error)
    } finally {
      clearAuth()
      router.push('/login')
    }
  }

  const refreshAuth = async () => {
    const refreshToken = localStorage.getItem('refresh_token')
    if (!refreshToken) {
      clearAuth()
      return
    }

    try {
      await refreshAuthToken(refreshToken)
    } catch (error) {
      clearAuth()
      router.push('/login')
    }
  }

  return (
    <AuthContext.Provider value={{ user, token, isLoading, login, register, logout, refreshAuth }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

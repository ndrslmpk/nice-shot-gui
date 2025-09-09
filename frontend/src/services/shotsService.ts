import type { Shot } from '@/models/shot'

export type ShotsResponse = {
  data: Shot[]
}

const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL as string | undefined) ?? ''

export async function fetchShots(limit: number = 100): Promise<Shot[]> {
  const url = `${API_BASE_URL}/shots?limit=${encodeURIComponent(limit)}`
  const response = await fetch(url, {
    method: 'GET',
    headers: { Accept: 'application/json' },
  })
  if (!response.ok) {
    throw new Error(`Failed to fetch shots: ${response.status} ${response.statusText}`)
  }
  const contentType = response.headers.get('content-type') || ''
  if (!contentType.includes('application/json')) {
    throw new Error('Unexpected response format from /shots')
  }
  const json = await response.json()

  console.log('', json)
  // Accept either array or { data: [] }
  const records: Shot[] = Array.isArray(json) ? json : (json?.data ?? [])
  return records
}

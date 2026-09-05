import React from 'react'

interface EmptyStateProps {
  message: string
}

export default function EmptyState({ message }: EmptyStateProps) {
  return (
    <div className="empty-state">
      <svg width="48" height="48" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <rect x="3" y="6" width="18" height="15" rx="2" stroke="currentColor" strokeWidth="2" />
        <line x1="3" y1="10" x2="21" y2="10" stroke="currentColor" strokeWidth="2" />
      </svg>
      <p>{message}</p>
    </div>
  )
}

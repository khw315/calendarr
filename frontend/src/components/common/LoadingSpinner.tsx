import React from 'react'

interface LoadingSpinnerProps {
  label?: string
}

export default function LoadingSpinner({ label = 'Loading releases...' }: LoadingSpinnerProps) {
  return (
    <div className="loading">
      <div className="spinner" />
      <p>{label}</p>
    </div>
  )
}

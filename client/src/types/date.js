
export const formatDate = dateString => {
  if (dateString === '0001-01-01T00:00:00Z') {
    return `-`
  }
  
  const date = new Date(dateString)

  return `${date.getMonth()}.${date.getDate()}.${date.getFullYear()}`
}

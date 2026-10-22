// Split only list separators; the selected daemon validates each mailbox.
export function splitSMTPRecipients(value: string): string[] {
  const result: string[] = []
  let part = ''
  let quoted = false
  let escaped = false
  for (const character of value) {
    if (escaped) { part += character; escaped = false; continue }
    if (quoted && character === '\\') { part += character; escaped = true; continue }
    if (character === '"') quoted = !quoted
    if (!quoted && (character === ',' || character === '\n')) {
      if (part.trim()) result.push(part.trim())
      part = ''
    } else part += character
  }
  if (part.trim()) result.push(part.trim())
  return result
}

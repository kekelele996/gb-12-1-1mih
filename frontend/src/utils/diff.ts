// Simple line-based diff for document versions.
export interface DiffLine {
  type: 'same' | 'add' | 'remove'
  text: string
}

export function diffLines(oldText: string, newText: string): DiffLine[] {
  const oldLines = oldText.split('\n')
  const newLines = newText.split('\n')
  const oldSet = new Set(oldLines)
  const newSet = new Set(newLines)
  const result: DiffLine[] = []
  const max = Math.max(oldLines.length, newLines.length)
  for (let i = 0; i < max; i++) {
    const a = oldLines[i]
    const b = newLines[i]
    if (a === undefined && b !== undefined) {
      result.push({ type: 'add', text: b })
    } else if (b === undefined && a !== undefined) {
      result.push({ type: 'remove', text: a })
    } else if (a === b) {
      result.push({ type: 'same', text: a })
    } else {
      if (!newSet.has(a)) result.push({ type: 'remove', text: a })
      if (!oldSet.has(b)) result.push({ type: 'add', text: b })
      else result.push({ type: 'same', text: b })
    }
  }
  return result
}

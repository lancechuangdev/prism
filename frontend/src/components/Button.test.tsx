import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'

import { Button } from './Button'

describe('Button', () => {
  it('supports keyboard and pointer activation', async () => {
    const onClick = vi.fn()
    render(<Button onClick={onClick}>Continue</Button>)
    await userEvent.click(screen.getByRole('button', { name: 'Continue' }))
    expect(onClick).toHaveBeenCalledOnce()
  })
})

import '@testing-library/jest-dom/vitest'
import { cleanup, configure } from '@testing-library/react'
import { afterEach } from 'vitest'

configure({ asyncUtilTimeout: 10_000 })

document.documentElement.lang = 'en'
document.title = 'go-schedule'
afterEach(cleanup)

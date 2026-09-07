import '@testing-library/jest-dom/vitest'
import { cleanup } from '@testing-library/react'
import { afterEach } from 'vitest'

document.documentElement.lang = 'en'
document.title = 'go-schedule foundation proof'

afterEach(cleanup)

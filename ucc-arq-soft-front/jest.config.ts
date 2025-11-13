import type { Config } from 'jest';

const config: Config = {
  testEnvironment: 'jsdom',
  roots: ['<rootDir>/src'],
  setupFilesAfterEnv: ['<rootDir>/jest.setup.ts'],
  transform: {
    '^.+\\.(ts|tsx)$': ['ts-jest', { tsconfig: '<rootDir>/tsconfig.jest.json' }],
    '^.+\\.(js|jsx|mjs)$': 'babel-jest',
  },
  transformIgnorePatterns: [
    // Transform specific ESM deps for Jest (cross-platform): msw and until-async
    '/node_modules/(?!((msw|@mswjs|until-async)/))'
  ],
  moduleFileExtensions: ['ts', 'tsx', 'js', 'jsx', 'json'],
  moduleNameMapper: {
    // Map CSS modules and global styles to identity-obj-proxy or a custom mock
    '^.+\\.(css|scss|sass)$': '<rootDir>/__mocks__/styleMock.js',
    // Explicitly mock Swiper CSS entrypoints (they are imported as modules like 'swiper/css')
    '^swiper/css$': '<rootDir>/__mocks__/styleMock.js',
    '^swiper/css/.*$': '<rootDir>/__mocks__/styleMock.js',
    '^.+\\.(png|jpg|jpeg|gif|svg)$': '<rootDir>/__mocks__/fileMock.js',
    '^next/navigation$': '<rootDir>/__mocks__/next/navigation.js',
    '^next/image$': '<rootDir>/__mocks__/next/image.js',
    'swiper/react': '<rootDir>/__mocks__/swiper-react.tsx',
    'swiper/modules': '<rootDir>/__mocks__/swiper-modules.ts',
    '^@/(.*)$': '<rootDir>/src/$1'
  },
  // Whitelist only the production runtime code we care about for coverage.
  // This keeps any type-only folders from accidentally showing up in reports.
  collectCoverageFrom: [
    'src/app/**/*.{ts,tsx}',
    'src/components/**/*.{ts,tsx}',
    'src/context/**/*.{ts,tsx}',
    'src/utils/**/*.{ts,tsx}',
    'src/providers/**/*.{ts,tsx}',
    // Optional: include hooks or other top-level dirs if present
    'src/**/hooks/**/*.{ts,tsx}',
    // Excludes
    '!src/**/__tests__/**',
    '!src/**/*.d.ts',
    '!src/**/?(*.)+(spec|test).{ts,tsx}',
    '!src/**/__mocks__/**',
    // Exclude any folder named "types" anywhere
    '!**/types/**'
  ],
  // Ensure any imported type files under src/types are ignored even if required by code
  coveragePathIgnorePatterns: [
    '/node_modules/',
    // Regex-style patterns; handle Windows \\ and POSIX / separators
    '[/\\\\]src[/\\\\]types[/\\\\]',
    // And ignore any nested "types" segment regardless of root
    '[/\\\\]types[/\\\\]'
  ],
  coverageThreshold: {
    global: {
      // Tighten thresholds to target 80% global statements/lines
      statements: 80,
      lines: 80,
      functions: 75,
      branches: 60
    }
  },
  coverageReporters: [
    ['html', { skipEmpty: true }],
    'text',
    'lcov'
  ],
  // Do not ignore any tests; ensure all suites run
  testMatch: ['**/__tests__/**/*.test.ts?(x)'],
  clearMocks: true
};

export default config;
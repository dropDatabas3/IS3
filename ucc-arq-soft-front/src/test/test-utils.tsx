import React, { ReactElement } from 'react';
import { render as rtlRender } from '@testing-library/react';
import { AuthProvider } from '@/context/auth/AuthProvider';
import { CoursesProvider } from '@/context/courses/CoursesProvider';
import { UiProvider } from '@/context/ui/UiProvider';

// Future: allow dependency injection (apiClient) via context props once refactored
interface CustomRenderOptions {
  wrapper?: React.ComponentType<{ children: React.ReactNode }>;
}

function Providers({ children }: { children: React.ReactNode }) {
  return (
    <UiProvider>
      <AuthProvider>
        <CoursesProvider>{children}</CoursesProvider>
      </AuthProvider>
    </UiProvider>
  );
}

export * from '@testing-library/react';

export function render(ui: ReactElement, options: CustomRenderOptions = {}) {
  const { wrapper, ...rest } = options;
  if (wrapper) {
    return rtlRender(ui, { wrapper, ...rest });
  }
  return rtlRender(ui, { wrapper: Providers, ...rest });
}

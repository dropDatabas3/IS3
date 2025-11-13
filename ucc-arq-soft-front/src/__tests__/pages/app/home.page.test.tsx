import React from 'react';
import { render, screen } from '@/test/test-utils';
import Home from '@/app/page';

describe('Home page', () => {
  it('shows main hero content', async () => {
    render(<Home />);
    // Text from MainSection hero
    expect(await screen.findByText(/Welcome to DuoMingo/i)).toBeInTheDocument();
  });
});

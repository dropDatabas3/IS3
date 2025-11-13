/* Simple next/image mock for Jest */
// Renders a normal img to avoid Next.js image warnings during tests
const React = require('react');

const NextImage = (props) => {
  const { src, alt, width, height, loader, fetchPriority, ...rest } = props || {};
  // Strip Next.js-specific props to avoid DOM warnings
  return React.createElement('img', { src, alt: alt || 'image', width, height, ...rest });
};
module.exports = NextImage;
